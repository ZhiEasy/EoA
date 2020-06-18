package models

type UpdateUserInfoReq struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Pwd string `json:"pwd"`
	CPwd string `json:"cpwd"`
}