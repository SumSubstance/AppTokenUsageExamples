package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"example/model"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
)

const URL = "https://api.sumsub.com"
const SumsubAppToken = "sbx:6L6rqHEtRVvBKKt7P1A03k2x.h6OsEOXWpyaXAjvBVNnx3ccXNGTBLHkw" // Example: sbx:uY0CgwELmgUAEyl4hNWxLngb.0WSeQeiYny4WEqmAALEAiK2qTC96fBad
const SumsubSecretKey = "EraepapR4Grr2vI1eZWtTkFDhbhsC5EI"                             // Example: Hej2ch71kG2kTd1iIUDZFNsO5C1lh5Gq
//Please don't forget to change token and secret key values to production ones when switching to production

func main() {
	var levelName = "basic-kyc-level"
	var externalUserId = uuid.NewString()

	var applicant = model.Applicant{}
	var fixedInfo = model.Info{}
	fixedInfo.Country = "GBR"
	fixedInfo.FirstName = "someName"
	applicant.FixedInfo = fixedInfo
	applicant.ExternalUserID = externalUserId

	// https://docs.sumsub.com/reference/create-applicant
	applicant = CreateApplicant(applicant, levelName)

	// https://docs.sumsub.com/reference/add-id-documents
	idDoc := AddDocument(applicant.ID)
	fmt.Println(idDoc)

	// https://docs.sumsub.com/reference/get-applicant-data
	applicant = GetApplicantInfo(applicant)

	// https://docs.sumsub.com/reference/generate-access-token-query
	accessToken := GenerateAccessToken(applicant)

	fmt.Println(accessToken.Token)
}

func GenerateAccessToken(applicant model.Applicant) model.AccessToken {
	postBody, _ := json.Marshal(model.AccessToken{
		UserId: applicant.ExternalUserID,
	})

	b, err := _makeSumsubRequest("/resources/accessTokens/sdk",
		"POST",
		"application/json",
		postBody)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(string(b))
	ioutil.WriteFile("generateAccessToken.json", b, 0777)

	var token model.AccessToken
	err = json.Unmarshal(b, &token)

	return token
}

func CreateApplicant(applicant model.Applicant, levelName string) model.Applicant {
	postBody, _ := json.Marshal(applicant)

	b, err := _makeSumsubRequest(
		"/resources/applicants?levelName="+levelName,
		"POST",
		"application/json",
		postBody)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(string(b))
	ioutil.WriteFile("createApplicant.json", b, 0777)

	var ac model.Applicant
	err = json.Unmarshal(b, &ac)
	if err != nil {
		log.Fatal(err)
	}

	return ac
}

func GetApplicantInfo(applicant model.Applicant) model.Applicant {
	p := fmt.Sprintf("/resources/applicants/%s/one", applicant.ID)
	b, err := _makeSumsubRequest(
		p,
		"GET",
		"application/json",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("getApplicant.json", b, 0777)

	var r model.Applicant
	err = json.Unmarshal(b, &r)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(r)

	return r
}

func AddDocument(applicantId string) model.IdDoc {
	file, err := os.Open("resources/images/sumsub-logo.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	meta, err := json.Marshal(model.IdDoc{
		IdDocType: "PASSPORT",
		Country:   "GBR",
	})

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	var fw io.Writer
	if fw, err = w.CreateFormFile("content", file.Name()); err != nil {
		log.Fatal(err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		log.Fatal(err)
	}

	if fw, err = w.CreateFormField("metadata"); err != nil {
		log.Fatal(err)
	}
	if _, err = io.Copy(fw, strings.NewReader(string(meta))); err != nil {
		log.Fatal(err)
	}
	w.Close()

	resp, err := _makeSumsubRequest(
		"/resources/applicants/"+applicantId+"/info/idDoc",
		"POST",
		w.FormDataContentType(),
		b.Bytes(),
	)

	var doc model.IdDoc
	err = json.Unmarshal(resp, &doc)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

// X-App-Token - an App Token that you generate in our dashboard
// X-App-Access-Sig - signature of the request in the hex format (see below)
// X-App-Access-Ts - number of seconds since Unix Epoch in UTC
func _makeSumsubRequest(path, method, contentType string, body []byte) ([]byte, error) {

	request, err := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	ts := fmt.Sprintf("%d", time.Now().Unix())

	request.Header.Add("X-App-Token", SumsubAppToken)

	request.Header.Add("X-App-Access-Sig", _sign(ts, SumsubSecretKey, method, path, &body))
	request.Header.Add("X-App-Access-Ts", ts)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", contentType)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return b, nil
}

func _sign(ts string, secret string, method string, path string, body *[]byte) string {
	hash := hmac.New(sha256.New, []byte(secret))
	data := []byte(ts + method + path)

	if body != nil {
		data = append(data, *body...)
	}

	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
