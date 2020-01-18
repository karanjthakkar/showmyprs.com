const firebase = require('firebase');

let FirebaseDatabase;

const setupFirebase = async () => {
  const app = firebase.initializeApp({
    apiKey: process.env.FIREBASE_API_KEY,
    authDomain: process.env.FIREBASE_AUTH_DOMAIN,
    databaseURL: process.env.FIREBASE_DATABASE_URL,
    projectId: process.env.FIREBASE_PROJECT_ID,
    storageBucket: process.env.FIREBASE_STORAGE_BUCKET,
    messagingSenderId: process.env.FIREBASE_MESSAGING_SENDER_ID,
    appId: process.env.FIREBASE_APP_ID
  });
  FirebaseDatabase = await firebase.database(app);
};

const writeUserDataForUser = async (username, result) => {  
  if (result) {
    result.data = JSON.stringify(result.data);
  }
  await FirebaseDatabase.ref(`users/${username}`).set(result);
};

const readUserDataForUser = async (username) => {
  const data = await FirebaseDatabase.ref(`users/${username}`).once('value');
  const result = data.val();
  if (result) {
    result.data = JSON.parse(result.data);
  }
  return result;
};

module.exports = {
  setupFirebase,
  writeUserDataForUser,
  readUserDataForUser,
};