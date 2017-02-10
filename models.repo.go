package main

type Repo struct {
  Url string `json:"url"`
  Name string `json:"name"`
  Stars int `json:"stars"`
  Forks int `json:"forks"`
  PullRequests []PullRequest `json:"prs"`
}
