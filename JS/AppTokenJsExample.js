const axios = require('axios');
const crypto = require('crypto');
const fs = require('fs');
const FormData = require('form-data');

// These parameters should be used for all requests
const SUMSUB_APP_TOKEN = 'YOUR_SUMSUB_APP_TOKEN'; // Example: tst:uY0CgwELmgUAEyl4hNWxLngb.0WSeQeiYny4WEqmAALEAiK2qTC96fBad
const SUMSUB_SECRET_KEY = 'YOUR_SUMSUB_SECRET_KEY'; // Example: Hej2ch71kG2kTd1iIUDZFNsO5C1lh5Gq
const SUMSUB_BASE_URL = 'https://test-api.sumsub.com'; // Please don't forget to change when switching to production

var config = {};
config.baseURL= SUMSUB_BASE_URL;

axios.interceptors.request.use(createSignature, function (error) {
  return Promise.reject(error);
})

// This function creates signature for the request as described here: https://developers.sumsub.com/api-reference/#app-tokens

function createSignature(config) {
  console.log('Creating a signature for the request...');

  var ts = Math.floor(Date.now() / 1000);
  const signature = crypto.createHmac('sha256',  SUMSUB_SECRET_KEY);
  signature.update(ts + config.method.toUpperCase() + config.url);

  if (config.data instanceof FormData) {
    signature.update (config.data.getBuffer());
  } else if (config.data) {
    signature.update (config.data);
  }

  config.headers['X-App-Access-Ts'] = ts;
  config.headers['X-App-Access-Sig'] = signature.digest('hex');

  return config;
}

// These functions configure requests for specified method

// https://developers.sumsub.com/api-reference/#creating-an-applicant
function createApplicant(externalUserId) {
  console.log("Creating an applicant...");

  var method = 'post';
  var url = '/resources/applicants';
  var ts = Math.floor(Date.now() / 1000);
  
  // Setup requierd documents set
  var requiredIdDocs = {
      docSets: 
      [
          {
              idDocSetType: 'IDENTITY',
              types: ['PASSPORT','ID_CARD','DRIVERS','RESIDENCE_PERMIT']
          },
          {
              idDocSetType: 'SELFIE',
              types: ['SELFIE']
          }
      ]
  };

  var body = {
      externalUserId: externalUserId,
      requiredIdDocs: requiredIdDocs
  };

  var headers = {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'X-App-Token': SUMSUB_APP_TOKEN
  };

  config.method = method;
  config.url = url;
  config.headers = headers;
  config.data = JSON.stringify(body);

  return config;
}

// https://developers.sumsub.com/api-reference/#adding-an-id-document
function addDocument(applicantId) {
  console.log("Adding document to the applicant...");

  var method = 'post';
  var url = `/resources/applicants/${applicantId}/info/idDoc`;
  var filePath = 'resources/sumsub-logo.png';

  var metadata = {
      idDocType: 'PASSPORT',
      country: 'GBR'
  };

  var form = new FormData();
  form.append('metadata', JSON.stringify(metadata));
  var content = fs.readFileSync(filePath); 
  form.append('content', content, filePath);

  var headers = {
    'Accept': 'application/json',
    'X-App-Token': SUMSUB_APP_TOKEN
  };

  config.method = method;
  config.url = url;
  config.headers = Object.assign(headers, form.getHeaders());
  config.data = form;

  return config;
}

// https://developers.sumsub.com/api-reference/#getting-applicant-status-sdk
function getApplicantStatus(applicantId) {
  console.log("Getting the applicant status...");

  var method = 'get';
  var url = `/resources/applicants/${applicantId}/status`;

  var headers = {
    'Accept': 'application/json',
    'X-App-Token': SUMSUB_APP_TOKEN
  };

  config.method = method;
  config.url = url;
  config.headers = headers;
  config.data = null;

  return config;
}

// https://developers.sumsub.com/api-reference/#access-tokens-for-sdks
function createAccessToken (externalUserId, ttlInSecs = 600) {
  console.log("Creating an access token for initializng SDK...");

  var method = 'post';
  var url = `/resources/accessTokens?userId=${externalUserId}&ttlInSecs=${ttlInSecs}`;

  var headers = {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'X-App-Token': SUMSUB_APP_TOKEN
  };

  config.method = method;
  config.url = url;
  config.headers = headers;
  config.data = null;

  return config;
}

// This section contains requests to server using configuration functions
// The description of the flow can be found here: https://developers.sumsub.com/api-flow/#api-integration-phases

        // Such actions are presented below:
        // 1) Creating an applicant
        // 2) Adding a document to the applicant
        // 3) Getting applicant status
        // 4) Getting access tokens for SDKs

async function main() {
  externalUserId = "random-JSToken-" + Math.random().toString(36).substr(2, 9);
  console.log("External UserID: ", externalUserId); 

  response = await axios(createApplicant(externalUserId))
    .then(function (response) {
      console.log("Response:\n", response.data);
      return response;
    })
    .catch(function (error) {
      console.log("Error:\n", error.response.data);
    });
  
  const applicantId = response.data.id;
  console.log("ApplicantID: ", applicantId);

  response = await axios(addDocument(applicantId))
  .then(function (response) {
    console.log("Response:\n", response.data);
    return response;
  })
  .catch(function (error) {
    console.log("Error:\n", error.response.data);
  });

  response = await axios(getApplicantStatus(applicantId))
  .then(function (response) {
    console.log("Response:\n", response.data);
    return response;
  })
  .catch(function (error) {
    console.log("Error:\n", error.response.data);
  });

  response = await axios(createAccessToken(externalUserId, 1200))
  .then(function (response) {
    console.log("Response:\n", response.data);
    return response;
  })
  .catch(function (error) {
    console.log("Error:\n", error.response.data);
  });

}

main();
