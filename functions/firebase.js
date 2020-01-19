const admin = require('firebase-admin');
admin.initializeApp();

const writeUserDataForUser = async (username, result) => {
  if (result) {
    result.data = JSON.stringify(result.data);
  }
  await admin.database().ref(`users/${username}`).set(result);
};

const readUserDataForUser = async (username) => {
  const data = await admin.database().ref(`users/${username}`).once('value');
  const result = data.val();
  if (result) {
    result.data = JSON.parse(result.data);
  }
  return result;
};

module.exports = {
  writeUserDataForUser,
  readUserDataForUser,
};