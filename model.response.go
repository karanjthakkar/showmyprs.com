package main

type Response struct {
	Username   string `json:"username"`
	TotalRepos int    `json:"totalRepos"`
	TotalPrs   int    `json:"totalPrs"`
	AllRepos   []Repo `json:"allRepos"`
}
