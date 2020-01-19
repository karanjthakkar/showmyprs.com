const functions = require('firebase-functions');
const { getPullRequestsForUser } = require('./github.js');
const _ = require('lodash');

exports.user = functions.https.onRequest(async (request, response) => {
  try {
    const username = _.get(request, 'params[0]', '').split('/')[1];
    if (username) {
      const data = await getPullRequestsForUser(username);
      response.send(data);
    } else {
      response.statusCode = 500;
      response.send({
        error: 'Please provide a username'
      });
    }
  } catch(error) {
    response.statusCode = 500;
    response.send({
      error: error.toString(),
      stacktrace: error.stack
    });
  }
});

