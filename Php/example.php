<?php
declare(strict_types = 1);

use Sumsub\AppTokenUsageExample\SumsubClient;

require_once __DIR__ . '/vendor/autoload.php';

// The description of the authorization method is available here: https://developers.sumsub.com/api-reference/#app-tokens
// Don't forget to change token and secret key values to production ones when switching to production
const SUMSUB_SECRET_KEY = 'SUMSUB_SECRET_KEY'; // Example: Hej2ch71kG2kTd1iIUDZFNsO5C1lh5Gq
const SUMSUB_APP_TOKEN = 'SUMSUB_APP_TOKEN'; // Example: sbx:uY0CgwELmgUAEyl4hNWxLngb.0WSeQeiYny4WEqmAALEAiK2qTC96fBad

// The description of the flow can be found here: https://developers.sumsub.com/api-flow/#api-integration-phases

// Such actions are presented below:
// 1) Creating an applicant
// 2) Adding a document to the applicant
// 3) Getting applicant status
// 4) Getting access token

$externalUserId = uniqid(); // Use your internal UserID instead in production code
$levelName = 'basic-kyc-level';

$testObject = new SumsubClient(SUMSUB_APP_TOKEN, SUMSUB_SECRET_KEY);

$applicantId = $testObject->createApplicant($externalUserId, $levelName);
echo 'The applicant was successfully created: ' . $applicantId . PHP_EOL;

$imageId = $testObject->addDocument(
    $applicantId,
    __DIR__ . '/resources/sumsub-logo.png',
    ['idDocType' => 'PASSPORT', 'country' => 'GBR'],
);
echo 'Identifier of the added document: ' . $imageId . PHP_EOL;

$applicantStatusInfo = $testObject->getApplicantStatus($applicantId);
echo 'Applicant status (json): ' . json_encode($applicantStatusInfo, JSON_PRETTY_PRINT) . PHP_EOL;

$accessTokenInfo = $testObject->getAccessToken($externalUserId, $levelName);
echo 'Access token (json): ' . json_encode($accessTokenInfo, JSON_PRETTY_PRINT) . PHP_EOL;