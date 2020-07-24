package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
)

const URL = "https://test-api.sumsub.com"

func main() {
	postBody := `{
	"externalUserId": "SomeExternalUserId",
	"info": {
    "country": "GBR",
    "firstName": "SomeUserName",
    "lastName": "SomeLastName",
    "phone": "+449112081223",
    "dob": "2000-03-04",
    "placeOfBirth": "SomeCityName"
	},
	"requiredIdDocs": {
		"docSets": [{
				"idDocSetType": "IDENTITY",
				"types": ["PASSPORT","ID_CARD","DRIVERS","RESIDENCE_PERMIT"]
			},
			{
				"idDocSetType": "SELFIE",
				"types": ["SELFIE"]
			},
			{
				"idDocSetType": "PROOF_OF_RESIDENCE",
				"types": ["UTILITY_BILL"]
			}
		]
	}
}`

	b, err := makeSumsubRequest(
		"/resources/applicants",
		"POST",
		"application/json",
		[]byte(postBody))
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(string(b))
	ioutil.WriteFile("createApplicant.json", b, 0777)

	var ac ApplicantCreateResponse
	err = json.Unmarshal(b, &ac)
	if err != nil {
		log.Fatal(err)
	}
	///resources/applicants/{applicantId}/one
	p := fmt.Sprintf("/resources/applicants/%s/one", ac.ID)
	b, err = makeSumsubRequest(
		p,
		"GET",
		"application/json",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("getApplicant.json", b, 0777)

	var r ApplicantsResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(r)

}

//X-App-Token - an App Token that you generate in our dashboard
//X-App-Access-Sig - signature of the request in the hex format (see below)
//X-App-Access-Ts - number of seconds since Unix Epoch in UTC
func makeSumsubRequest(path, method, contentType string, body []byte) ([]byte, error) {

	request, err := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	ts := fmt.Sprintf("%d", time.Now().Add(1*time.Minute).Unix())

	token := "tst:3FW3pkDK8zF86y7dPJwLhwB9.GyY8I5XL2Pkg2SKXjiMwWOWiVEtS6f0T"

	request.Header.Add("X-App-Token", token)

	secret := "jAVYqzpfuU9iewQ2PWdI7Q3knFkr083a"
	request.Header.Add("X-App-Access-Sig", sign(ts, secret, method, path, &body))
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

func sign(ts string, secret string, method string, path string, body *[]byte) string {
	hash := hmac.New(sha256.New, []byte(secret))
	data := []byte(ts + method + path)

	if body != nil {
		data = append(data, *body...)
	}

	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

type ApplicantCreateResponse struct {
	ID             string `json:"id"`
	CreatedAt      string `json:"createdAt"`
	Key            string `json:"key"`
	ClientID       string `json:"clientId"`
	InspectionID   string `json:"inspectionId"`
	ExternalUserID string `json:"externalUserId"`
	Info           struct {
		FirstName    string `json:"firstName"`
		FirstNameEn  string `json:"firstNameEn"`
		LastName     string `json:"lastName"`
		LastNameEn   string `json:"lastNameEn"`
		Dob          string `json:"dob"`
		PlaceOfBirth string `json:"placeOfBirth"`
		Country      string `json:"country"`
		Phone        string `json:"phone"`
	} `json:"info"`
	Env            string `json:"env"`
	RequiredIDDocs struct {
		DocSets []struct {
			IDDocSetType string   `json:"idDocSetType"`
			Types        []string `json:"types"`
		} `json:"docSets"`
	} `json:"requiredIdDocs"`
	Review struct {
		Reprocessing           bool   `json:"reprocessing"`
		CreateDate             string `json:"createDate"`
		ReviewStatus           string `json:"reviewStatus"`
		NotificationFailureCnt int    `json:"notificationFailureCnt"`
		Priority               int    `json:"priority"`
		AutoChecked            bool   `json:"autoChecked"`
	} `json:"review"`
	Type string `json:"type"`
}

type ApplicantsResponse struct {
	List struct {
		Items []struct {
			ID             string `json:"id"`
			CreatedAt      string `json:"createdAt"`
			Key            string `json:"key"`
			ClientID       string `json:"clientId"`
			InspectionID   string `json:"inspectionId"`
			ExternalUserID string `json:"externalUserId"`
			Info           struct {
				FirstName    string `json:"firstName"`
				FirstNameEn  string `json:"firstNameEn"`
				MiddleName   string `json:"middleName"`
				MiddleNameEn string `json:"middleNameEn"`
				LastName     string `json:"lastName"`
				LastNameEn   string `json:"lastNameEn"`
				Dob          string `json:"dob"`
				Gender       string `json:"gender"`
				Country      string `json:"country"`
				Phone        string `json:"phone"`
				IDDocs       []struct {
					IDDocType    string `json:"idDocType"`
					Country      string `json:"country"`
					FirstName    string `json:"firstName"`
					FirstNameEn  string `json:"firstNameEn"`
					MiddleName   string `json:"middleName"`
					MiddleNameEn string `json:"middleNameEn"`
					LastName     string `json:"lastName"`
					LastNameEn   string `json:"lastNameEn"`
				} `json:"idDocs"`
			} `json:"info"`
			Env               string `json:"env"`
			ApplicantPlatform string `json:"applicantPlatform"`
			RequiredIDDocs    struct {
				DocSets []struct {
					IDDocSetType  string   `json:"idDocSetType"`
					Types         []string `json:"types"`
					VideoRequired string   `json:"videoRequired,omitempty"`
				} `json:"docSets"`
			} `json:"requiredIdDocs"`
			Review struct {
				ElapsedSincePendingMs int    `json:"elapsedSincePendingMs"`
				ElapsedSinceQueuedMs  int    `json:"elapsedSinceQueuedMs"`
				Reprocessing          bool   `json:"reprocessing"`
				CreateDate            string `json:"createDate"`
				ReviewDate            string `json:"reviewDate"`
				StartDate             string `json:"startDate"`
				ReviewResult          struct {
					ReviewAnswer string `json:"reviewAnswer"`
				} `json:"reviewResult"`
				ReviewStatus           string `json:"reviewStatus"`
				NotificationFailureCnt int    `json:"notificationFailureCnt"`
				Priority               int    `json:"priority"`
				AutoChecked            bool   `json:"autoChecked"`
			} `json:"review"`
			Lang string `json:"lang"`
			Type string `json:"type"`
		} `json:"items"`
		TotalItems int `json:"totalItems"`
	} `json:"list"`
}
