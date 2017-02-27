package main

import (
	"bytes"
	"encoding/json"
	timeago "github.com/ararog/timeago"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Util to concat two strings
func concat(a string, b string) string {
	var buffer bytes.Buffer
	buffer.WriteString(a)
	buffer.WriteString(b)
	return buffer.String()
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

// Group pull requests by repo based on the repo url
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

// Save the response as JSON in the local cache
func SaveDataAsJson(data Response, username string) {
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

// Get time ago string based on the last synced time
func getTimeAgo(raw int64) string {
	end := time.Unix(raw, 0)
	start := time.Now()
	timeAgoInString, _ := timeago.TimeAgoWithTime(start, end)
	return timeAgoInString
}
