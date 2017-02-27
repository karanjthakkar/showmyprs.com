package main

type Response struct {
	LastSyncedAt              int64  `json:"last_synced_at"`
	LastSyncedAtString        string `json:"last_synced_at_string"`
	LastSyncedAtStringVerbose string `json:"last_synced_at_string_verbose"`
	Username                  string `json:"username"`
	TotalRepos                int    `json:"totalRepos"`
	TotalPrs                  int    `json:"totalPrs"`
	AllRepos                  []Repo `json:"allRepos"`
}
