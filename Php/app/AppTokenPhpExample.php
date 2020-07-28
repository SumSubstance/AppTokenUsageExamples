<?php

namespace App;

require_once __DIR__ . '/../vendor/autoload.php';

use GuzzleHttp;
use GuzzleHttp\Psr7\MultipartStream;

// The description of the authorization method is available here: https://developers.sumsub.com/api-reference/#app-tokens
define("SUMSUB_SECRET_KEY", "SUMSUB_SECRET_KEY"); // Example: Hej2ch71kG2kTd1iIUDZFNsO5C1lh5Gq
define("SUMSUB_APP_TOKEN", "SUMSUB_APP_TOKEN"); // Example: tst:uY0CgwELmgUAEyl4hNWxLngb.0WSeQeiYny4WEqmAALEAiK2qTC96fBad
define("SUMSUB_TEST_BASE_URL", "https://api.sumsub.com"); // Please don't forget to change when switching to production

class AppTokenGuzzlePhpExample
{
    public function createApplicant()
        // https://developers.sumsub.com/api-reference/#creating-an-applicant
    {
        $requestBody = [
            'externalUserId' => uniqid(),
            'requiredIdDocs' => [
                'docSets' => array(
                    [
                        'idDocSetType' => 'IDENTITY',
                        'types' => array('PASSPORT', 'ID_CARD', 'DRIVERS')
                    ],
                    [
                        'idDocSetType' => 'SELFIE',
                        'types' => array('SELFIE')
                    ]
                )
            ]];

        $url = '/resources/applicants';
        $request = new GuzzleHttp\Psr7\Request('POST', SUMSUB_TEST_BASE_URL . $url);
        $request = $request->withHeader('Content-Type', 'application/json');
        $request = $request->withBody(GuzzleHttp\Psr7\stream_for(json_encode($requestBody)));

        $responseBody = $this->sendHttpRequest($request, $url)->getBody();
        return json_decode($responseBody)->{'id'};
    }

    public function sendHttpRequest($request, $url)
    {
        $client = new GuzzleHttp\Client();
        $ts = round(strtotime("now"));

        $request = $request->withHeader('X-App-Token', SUMSUB_APP_TOKEN);
        $request = $request->withHeader('X-App-Access-Sig', $this->createSignature($ts, $request->getMethod(), $url, $request->getBody()));
        $request = $request->withHeader('X-App-Access-Ts', $ts);

        try {
            $response = $client->send($request);
            if ($response->getStatusCode() != 200 && $response->getStatusCode() != 201) {
                // https://developers.sumsub.com/api-reference/#errors
                // If an unsuccessful answer is received, please log the value of the "correlationId" parameter.
                // Then perhaps you should throw the exception. (depends on the logic of your code)
            }
        } catch (GuzzleHttp\Exception\GuzzleException $e) {
            error_log($e);
        }

        return $response;
    }

    private function createSignature($ts, $httpMethod, $url, $httpBody)
    {
        return hash_hmac('sha256', $ts . strtoupper($httpMethod) . $url . $httpBody, SUMSUB_SECRET_KEY);
    }

    public function addDocument($applicantId)
        // https://developers.sumsub.com/api-reference/#adding-an-id-document
    {
        $metadata = ['idDocType' => 'PASSPORT', 'country' => 'GBR'];
        $file = __DIR__ . '/resources/images/sumsub-logo.png';

        $multipart = new MultipartStream([
            [
                "name" => "metadata",
                "contents" => json_encode($metadata)
            ],
            [
                'name' => 'content',
                'contents' => fopen($file, 'r')
            ],
        ]);

        $url = "/resources/applicants/" . $applicantId . "/info/idDoc";
        $request = new GuzzleHttp\Psr7\Request('POST', SUMSUB_TEST_BASE_URL . $url);
        $request = $request->withBody($multipart);

        return $this->sendHttpRequest($request, $url)->getHeader("X-Image-Id")[0];
    }

    public function getApplicantStatus($applicantId)
        // https://developers.sumsub.com/api-reference/#getting-applicant-status-api
    {
        $url = "/resources/applicants/" . $applicantId . "/requiredIdDocsStatus";
        $request = new GuzzleHttp\Psr7\Request('GET', SUMSUB_TEST_BASE_URL . $url);

        return $this->sendHttpRequest($request, $url)->getBody();
    }

    public function getAccessToken($applicantId)
        // https://developers.sumsub.com/api-reference/#access-tokens-for-sdks
    {
        $url = "/resources/accessTokens?userId=" . $applicantId;
        $request = new GuzzleHttp\Psr7\Request('POST', SUMSUB_TEST_BASE_URL . $url);

        return $this->sendHttpRequest($request, $url)->getBody();
    }

}

// The description of the flow can be found here: https://developers.sumsub.com/api-flow/#api-integration-phases

// Such actions are presented below:
// 1) Creating an applicant
// 2) Adding a document to the applicant
// 3) Getting applicant status
// 4) Getting access token

$testObject = new AppTokenGuzzlePhpExample();

$applicantId = $testObject->createApplicant();
echo "The applicant was successfully created: " . $applicantId . PHP_EOL;

$imageId = $testObject->addDocument($applicantId);
echo "Identifier of the added document: " . $imageId . PHP_EOL;

$applicantStatusStr = $testObject->getApplicantStatus($applicantId);
echo "Applicant status (json string): " . $applicantStatusStr;

$accessTokenStr = $testObject->getAccessToken($applicantId);
echo "Access token (json string): " . $accessTokenStr;
