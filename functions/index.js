const functions = require('firebase-functions');
const { getPullRequestsForUser } = require('./github.js');

exports.user = functions.https.onRequest(async (request, response) => {
  try {
    const data = await getPullRequestsForUser(request.params[0]);
    response.send({
      ...data,
      config: process.env.FIREBASE_CONFIG
    });
  } catch(error) {
    response.send({
      error: error.toString(),
      stacktrace: error.stack,
      config: process.env.FIREBASE_CONFIG
    });
  }
});

