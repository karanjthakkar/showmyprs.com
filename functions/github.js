const Octokit = require('@octokit/rest');
const functions = require('firebase-functions');
const FirebaseUtils = require('./firebase');

const client = new Octokit({
  auth: functions.config().showmyprs.github_token,
});

const getRepoNameFromPullRequestUrl = (url) => url.replace('https://api.github.com/repos/', '');
const getRepoUrlFromPullRequestUrl = (url) => url.replace('https://api.github.com/repos/', 'https://github.com/');

const getPullRequestDetailsFromHtmlUrl = (url) => {
	const parts = url.split('/');
	const owner = parts[3];
	const repo = parts[4];
	const pull_number = parts[6];

	return {
		owner,
		repo,
		pull_number,
	};
};

const getMergedStatusForPullRequest = async (item) => {
  if (item.state === 'closed') {
    const params = getPullRequestDetailsFromHtmlUrl(item.html_url);
    try {
      const isMerged = await client.pulls.checkIfMerged(params);
      return 'merged';
    } catch(e) {
      return 'closed';
    }
  }

  return item.state;
};

const addMergedStateToItem = async (item) => {
  return {
    ...item,
    state: await getMergedStatusForPullRequest(item),
  };
};

const addItemsToResult = async (items, result) => {
  const promises = items.map(item => addMergedStateToItem(item));
  const modifiedItems = await Promise.all(promises);
  modifiedItems.forEach(modifiedItem => {
    const name = getRepoNameFromPullRequestUrl(modifiedItem.repository_url);
    if (result.hasOwnProperty(name)) {
      result[name].pullRequests.push(modifiedItem);
    } else {
      result[name] = {
        name,
        url: getRepoUrlFromPullRequestUrl(modifiedItem.repository_url),
        pullRequests: [modifiedItem],
      };
    }
  });
  return result;
};

const getPullRequestsForUser = async (username) => {
  const getPullRequests = async (page = 1, fetched_count = 0, result = {}) => {
    const data = await client.search.issuesAndPullRequests({
      q: `type:pr author:${username} is:public`,
      per_page: 100,
      page,
    });

    // Update fetched count to decide whether to continue fetch
    fetched_count += data.data.items.length;
    result = await addItemsToResult(data.data.items, result);

    // Github has an upper limit to allow searching only the first 1000 items
    // in any search result :(
    if (data.data.total_count > fetched_count && fetched_count < 1000) {
      page += 1;
      return getPullRequests(page, fetched_count, result);
    } else {
      return {
        username,
        data: result,
        total_prs: data.data.total_count,
        total_repos: Object.keys(result).length,
        last_fetched_at: Date.now()
      };
    }
  };

  const data = await FirebaseUtils.readUserDataForUser(username);

  if (data) {
    return data;
  } else {
    const data = await getPullRequests();
    await FirebaseUtils.writeUserDataForUser(username, data);
    return data;
  }
};

module.exports = {
  getPullRequestsForUser
};