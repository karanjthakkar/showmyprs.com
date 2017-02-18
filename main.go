package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strconv"
	"strings"
)

var router *gin.Engine
var client *github.Client

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

	// Define the route for the index page and display the index.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	router.GET("user/:username", func(c *gin.Context) {

		username := c.Param("username")
		responseType := c.Query("response_type")

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

			log.Println("waiting bro")

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

		if responseType == "json" {
			c.JSON(
				200,
				gin.H{
					"username":   username,
					"totalRepos": len(allRepos),
					"totalPrs":   len(allPrs),
					"allRepos":   allRepos,
				},
			)
		} else {
			// Call the HTML method of the Context to render a template
			c.HTML(
				// HTTP status
				200,
				// Use the index.html template
				"index.html",
				// Pass the data that the page uses (in this case, 'title')
				gin.H{
					"username":   username,
					"totalRepos": len(allRepos),
					"totalPrs":   len(allPrs),
					"allRepos":   allRepos,
				},
			)
		}

	})

	// Start serving the application
	router.Run()

}
