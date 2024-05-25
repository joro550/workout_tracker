package user_pages

import "net/http"

type RegisterModel struct {
	UserExists bool
}

type LoginModel struct {
	Error bool
}

type UserLogin struct {
	Username string
	Password string
}

func UserLoginFromRequest(r *http.Request) UserLogin {
	return UserLogin{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
}
