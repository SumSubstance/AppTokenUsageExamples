package model

type Info struct {
	FirstName    string  `json:"firstName,omitempty"`
	FirstNameEn  string  `json:"firstNameEn,omitempty"`
	MiddleName   string  `json:"middleName,omitempty"`
	MiddleNameEn string  `json:"middleNameEn,omitempty"`
	LastName     string  `json:"lastName,omitempty"`
	LastNameEn   string  `json:"lastNameEn,omitempty"`
	Dob          string  `json:"dob,omitempty"` //yyyy-mm-dd format
	Gender       string  `json:"gender,omitempty"`
	Country      string  `json:"country,omitempty"`
	Phone        string  `json:"phone,omitempty"`
	IdDocs       []IdDoc `json:"idDocs,omitempty"`
}
