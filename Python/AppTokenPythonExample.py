import time
import requests
import hmac
import hashlib
import json

class Const:
    SUMSUB_SECRET_KEY = "YOUR_SECRET_KEY"  # Example: Hej2ch71kG2kTd1iIUDZFNsO5C1lh5Gq
    SUMSUB_APP_TOKEN = "YOUR_APP_TOKEN"  # Example: tst:uY0CgwELmgUAEyl4hNWxLngb.0WSeQeiYny4WEqmAALEAiK2qTC96fBad
    SUMSUB_TEST_BASE_URL = "https://test-api.sumsub.com" # Please don't forget to change when switching to production


CONST = Const()
applicantId = 0
imageId = 0

def createApplicant():
# https://developers.sumsub.com/api-reference/#creating-an-applicant
    global applicantId
    body = {"externalUserId": 'SomeUserID'}
    headers = {'Content-Type': 'application/json',
               'Content-Encoding': 'utf-8'
               }
    resp = sign_request(requests.Request('POST', CONST.SUMSUB_TEST_BASE_URL+'/resources/applicants?levelName=basic-kyc-level',
                                         data = json.dumps(body),
                                         headers=headers
                                         ))
    s = requests.Session()
    ourresponse = s.send(resp)
    applicantId = (ourresponse.json()['id'])
    print('The applicant was successfully created:', applicantId)


def addDocument():
# https://developers.sumsub.com/api-reference/#adding-an-id-document
    global imageId, applicantId
    with open('img.jpg', 'wb') as handle:
        response = requests.get('https://fv2-1.failiem.lv/thumb_show.php?i=gdmn9sqy&view', stream=True)
        if not response.ok:
            print(response)

        for block in response.iter_content(1024):
            if not block:
                break
            handle.write(block)
    payload = {"metadata": '{"idDocType":"PASSPORT", "country":"USA"}'}
    resp = sign_request(
        requests.Request('POST',  CONST.SUMSUB_TEST_BASE_URL+'/resources/applicants/'+applicantId+'/info/idDoc',
                         data=payload,
                         files=[('content', open('img.jpg', 'rb'))]
                         ))
    sw = requests.Session()
    ourresponse = sw.send(resp)
    imageId = (ourresponse.headers['X-Image-Id'])
    print('Identifier of the added document:', imageId)


def getApplicantStatus():
# https://developers.sumsub.com/api-reference/#getting-applicant-status-api
    global applicantId
    url =  CONST.SUMSUB_TEST_BASE_URL+'/resources/applicants/'+applicantId+'/requiredIdDocsStatus'
    resp = sign_request(requests.Request('GET', url))
    s = requests.Session()
    ourresponse = s.send(resp)
    print(ourresponse.text)

def getAccessToken():
# https://developers.sumsub.com/api-reference/#access-tokens-for-sdks
    global applicantId
    params = {"userId": applicantId, "ttlInSecs": '600'}
    headers = {'Content-Type': 'application/json',
               'Content-Encoding': 'utf-8'
               }
    resp = sign_request(requests.Request('POST', CONST.SUMSUB_TEST_BASE_URL+'/resources/accessTokens',
                                         params=params,
                                         headers=headers
                                         ))
    s = requests.Session()
    ourresponse = s.send(resp)
    print(ourresponse.text)
    token = (ourresponse.json()['token'])
    print('Token:', token)

def sign_request(request: requests.Request) -> requests.PreparedRequest:
    prepared_request = request.prepare()
    now = int(time.time())
    method = request.method.upper()
    path_url = prepared_request.path_url  # includes encoded query params
    # could be None so we use an empty **byte** string here
    body = b'' if prepared_request.body is None else prepared_request.body
    if type(body) == str:
        body = body.encode('utf-8')
    data_to_sign = str(now).encode('utf-8') + method.encode('utf-8') + path_url.encode('utf-8') + body
    # hmac needs bytes
    signature = hmac.new(
        CONST.SUMSUB_SECRET_KEY.encode('utf-8'),
        data_to_sign,
        digestmod=hashlib.sha256
    )
    prepared_request.headers['X-App-Token'] = CONST.SUMSUB_APP_TOKEN
    prepared_request.headers['X-App-Access-Ts'] = str(now)
    prepared_request.headers['X-App-Access-Sig'] = signature.hexdigest()
    return prepared_request

 # Such actions are presented below:
 # 1) Creating an applicant
 # 2) Adding a document to the applicant
 # 3) Getting applicant status
 # 4) Getting access token
createApplicant()
addDocument()
getApplicantStatus()
getAccessToken()
