const Hapi = require('@hapi/hapi');
const { getPullRequestsForUser } = require('./github');
const { setupFirebase } = require('./firebase');

require('./setup');

const handler = async (req, h) => {
  try {
    const data = await getPullRequestsForUser(req.params.username);
    return data;
    res.json(data).code(200);
  } catch(e) {
    return h.response({
      error: e.toString(),
      stacktrace: e.stack
    }).code(500);
  }
};

const init = async () => {
  const server = Hapi.server({
    port: 3000,
    host: 'localhost'
  });

  server.route({
    method: 'GET',
    path: '/user/{username}',
    handler,
  });

  await setupFirebase();
  await server.start();
  console.log('Server running on %s', server.info.uri);
};

process.on('unhandledRejection', (err) => {
  console.log(err);
  process.exit(1);
});

init();
