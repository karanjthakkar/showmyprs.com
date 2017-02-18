package main

type PullRequest struct {
	Url     string `json:"url"`
	Title   string `json:"title"`
	State   string `json:"state"`
	RepoUrl string `json:"-"`
}
