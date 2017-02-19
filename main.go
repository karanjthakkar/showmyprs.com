package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var router *gin.Engine
var client *github.Client
var CACHE_PATH = "./.cache/"

func concat(a string, b string) string {
	var buffer bytes.Buffer
	buffer.WriteString(a)
	buffer.WriteString(b)
	return buffer.String()
}

func setupGithub() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: os.Getenv("GITHUB_TOKEN"),
		},
	)

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client = github.NewClient(tc)
}

// Parse the github PR url of the format:
// https://github.com/ubilabs/react-geosuggest/pull/253
// to retrive the owner: ubilabs, repo: react-geosuggest, number: 253
func parseGithubPrUrl(url string) ParsedPrUrl {
	prParts := strings.Split(url, "/")
	owner := prParts[3]
	repo := prParts[4]
	pullNumber, _ := strconv.Atoi(prParts[6])

	return ParsedPrUrl{
		Owner:  owner,
		Repo:   repo,
		Number: pullNumber,
	}
}

func groupByRepo(allPrs []PullRequest, CachedRepo map[string]Repo) []Repo {
	var allRepos []Repo
	currentRepoIndex := len(allRepos)
	urlMap := make(map[string]int)
	for _, pr := range allPrs {
		// If not found in map, its a new repo
		// Save its index
		// If found, push the PR object to the pr list
		// inside the repo object
		if _, isCached := urlMap[pr.RepoUrl]; isCached == false {
			urlMap[pr.RepoUrl] = currentRepoIndex
			allRepos = append(allRepos, CachedRepo[pr.RepoUrl])
			allRepos[currentRepoIndex].PullRequests = append(allRepos[currentRepoIndex].PullRequests, pr)
			currentRepoIndex = currentRepoIndex + 1
		} else {
			allRepos[urlMap[pr.RepoUrl]].PullRequests = append(allRepos[urlMap[pr.RepoUrl]].PullRequests, pr)
		}
	}
	return allRepos
}

func SaveDataAsJson(data gin.H, username string) {
	dir := CACHE_PATH
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(dir, 0755)
		} else {
			log.Println(err)
		}
	}

	path := concat(dir, username)
	os.Remove(path)

	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	ioutil.WriteFile(path, b, 0644)
}

func main() {

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Setup template folder
	router.LoadHTMLGlob("templates/*")

	router.Static("/assets", "./assets")

	// Setup Github client
	setupGithub()

	// Cache repo data
	CachedRepo := make(map[string]Repo)

	// Default landing page
	router.GET("/", func(c *gin.Context) {
		// Call the HTML method of the Context to render a template
		c.HTML(
			// HTTP status
			200,
			// Use the index.html template
			"index.html",
			gin.H{},
		)
	})

	// User handler
	router.GET("user/:username", func(c *gin.Context) {

		username := c.Param("username")
		responseType := c.Query("response_type")
		cachePath := concat(CACHE_PATH, username)

		var resultData gin.H

		if _, err := os.Stat(cachePath); err == nil {
			cacheData, err := ioutil.ReadFile(cachePath)
			// Error reading the file
			if err != nil {
				log.Println(err)
			}
			data := gin.H{}
			err = json.Unmarshal(cacheData, &data)
			// Error unmarshaling the file
			if err != nil {
				log.Println(err)
			}
			resultData = data
		} else {
			// Search options to override the default 30 and fetch max 100 per page
			opt := &github.SearchOptions{
				ListOptions: github.ListOptions{
					PerPage: 100,
				},
			}

			var allRepos []Repo
			var allPrs []PullRequest

			// Continuously fetch all PR's
			for {
				prs, resp, _ := client.Search.Issues(
					// Search query to find PR's
					fmt.Sprintf("type:pr author:%s is:public", username),
					opt,
				)

				// Iterate over all closed pull requests to
				// see which of them is merged and which one isn't
				// Also for each PR we are going to fetch the actual
				// repo's stats such as stars, pr's etc.
				for _, githubPrObject := range prs.Issues {
					parsedPrUrl := parseGithubPrUrl(*githubPrObject.HTMLURL)

					// Get stats
					repoUrl := strings.Join(
						[]string{
							"https://api.github.com",
							parsedPrUrl.Owner,
							parsedPrUrl.Repo,
						},
						"/",
					)

					// Cache repo stats and only make calls for new ones
					var repo Repo
					if _, isCached := CachedRepo[repoUrl]; isCached == false {
						repoData, _, _ := client.Repositories.Get(parsedPrUrl.Owner, parsedPrUrl.Repo)
						repo = Repo{
							Stars:        *repoData.StargazersCount,
							Forks:        *repoData.ForksCount,
							Name:         *repoData.FullName,
							Url:          *repoData.HTMLURL,
							PullRequests: []PullRequest{},
						}
						CachedRepo[repoUrl] = repo
					} else {
						repo = CachedRepo[repoUrl]
					}

					// Get merged status
					if *githubPrObject.State == "closed" {
						isPrMerged, _, _ := client.PullRequests.IsMerged(
							parsedPrUrl.Owner,
							parsedPrUrl.Repo,
							parsedPrUrl.Number,
						)

						if isPrMerged {
							*githubPrObject.State = "merged"
						}
					}

					pr := PullRequest{
						Url:     *githubPrObject.HTMLURL,
						Title:   *githubPrObject.Title,
						State:   *githubPrObject.State,
						RepoUrl: repoUrl,
					}

					allPrs = append(allPrs, pr)

				}

				if resp.NextPage == 0 {
					break
				}

				opt.ListOptions.Page = resp.NextPage
			}

			allRepos = groupByRepo(allPrs, CachedRepo)

			data := gin.H{
				"username":   username,
				"totalRepos": len(allRepos),
				"totalPrs":   len(allPrs),
				"allRepos":   allRepos,
			}

			resultData = data

			go SaveDataAsJson(data, username)
		}

		if responseType == "json" {
			c.JSON(
				200,
				resultData,
			)
		} else {
			// Call the HTML method of the Context to render a template
			c.HTML(
				// HTTP status
				200,
				// Use the user.html template
				"user.html",
				// Pass the data that the page uses (in this case, 'title')
				resultData,
			)
		}

	})

	// Start serving the application
	router.Run()

}
