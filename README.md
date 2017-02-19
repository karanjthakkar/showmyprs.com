# Show My PR's

## Motivation 😰

I wanted to create an `Open Source` section on [my website](https://karanjthakkar.com) to showcase some of the contributions that I have done over the years. Finding the PR's, sorting them, grouping them was a big pain. I couldnt't find any service that did this already. 


## Features 💅🏻💥

1. Show the actual status of each pull request: `open` / `closed` / `merged`
2. All your Pull Requests grouped by repositories: `/user/<username>`
3. Add `?response_type=json` to consume the json response directly


## Implementation notes 🙇🏻

1. Fetch all the pull requests from the github `search/issues` endpoint

  `https://api.github.com/search/issues?q=type:pr+author:gaearon+is:public`

2. For each pull request, fetch the repository data (stars/forks)

  `https://api.github.com/repos/:owner/:repo`

  **Note**: *The repository data is cached in-memory and is reused between subsequent calls.*

3. The search endpoint returns status as `closed` even for pull requests that have been merged. So, for each closed pull request, the merge status is fetched using the below endpoint.

  `https://api.github.com/repos/:owner/:repo/pulls/:number/merge`

4. These pull requests are then grouped based on the repository and then sent to the client.


## Contributing 👯

- Make sure you have installed Go and setup your workspace
- Get the latest code: `go get github.com/karanjthakkar/showmyprs.com`
- Start the server: `go build && ./showmyprs.com` (This requires the `GITHUB_TOKEN` environment variable. [Heres how](https://github.com/blog/1509-personal-api-tokens) you can get yours)

If you need help figuring out how to contribute (since its written in Go), hit me up on [Twitter](https://twitter.com/geekykaran) or [Email](mailto:karanjthakkar@gmail.com). I would love to help you get set up ☺️


## Deploying 🚀

- Pre-requisite: `npm install`
- `grunt build deploy`: This builds the binary for my Ubuntu linux instance on Amazon EC2 and uploads it along with the web assets via `scp`.



## Coming Next 🔥 (**Want to help? 👇🏻**)

- [ ] Sorting based on the stars/forks that a repository has ([Issue #4](https://github.com/karanjthakkar/showmyprs.com/issues/))
- [ ] Sorting based on the number of PR's for a repository ([Issue #4](https://github.com/karanjthakkar/showmyprs.com/issues/))
- [ ] Response caching for faster profile loads ([Issue #2](https://github.com/karanjthakkar/showmyprs.com/issues/))
- [ ] Filters based on PR state: Open/Merged/Closed ([Issue #3](https://github.com/karanjthakkar/showmyprs.com/issues/))


## License

MIT © [Karan Thakkar](https://karanjthakkar.com)
