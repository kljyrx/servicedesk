package controllers

var (
	// UserContro 所有的controller类声明都在这儿
	UserContro = &UserController{}
)

type Response struct {
	Message string
}

type ResponseLogin struct {
	Response
	Token string
}
