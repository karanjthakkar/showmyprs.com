// const functions = require('firebase-functions');
// const firebase = require('firebase');
const admin = require('firebase-admin');
admin.initializeApp();

// const firebaseConfig = {
//   apiKey: functions.config().showmyprs.firebase_api_key,
//   authDomain: functions.config().showmyprs.firebase_auth_domain,
//   databaseURL: functions.config().showmyprs.firebase_database_url,
//   projectId: functions.config().showmyprs.firebase_project_id,
//   storageBucket: functions.config().showmyprs.firebase_storage_bucket,
//   messagingSenderId: functions.config().showmyprs.firebase_messaging_sender_id,
//   appId: functions.config().showmyprs.firebase_app_id
// };

const writeUserDataForUser = async (username, result) => {
  // const app = firebase.initializeApp(firebaseConfig);
  // const database = await firebase.database(app);
  if (result) {
    result.data = JSON.stringify(result.data);
  }
  await admin.database().ref(`users/${username}`).set(result);
  // await database.goOffline();
  // await app.delete();
};

const readUserDataForUser = async (username) => {
  // const app = firebase.initializeApp(firebaseConfig);
  // const database = await firebase.database(app);
  const data = await admin.database().ref(`users/${username}`).once('value');
  const result = data.val();
  if (result) {
    result.data = JSON.parse(result.data);
  }
  // await database.goOffline();
  // await app.delete();
  return result;
};

module.exports = {
  writeUserDataForUser,
  readUserDataForUser,
};