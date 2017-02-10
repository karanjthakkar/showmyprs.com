package main

type ParsedPrUrl struct {
  Owner string `json:"owner"`
  Repo string `json:"repo"`
  Number int `json:"pullNumber"`
}