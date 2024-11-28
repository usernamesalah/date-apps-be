package model

type User struct {
	UID         string `json:"uid"`
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Password    string `json:"-"`
}
