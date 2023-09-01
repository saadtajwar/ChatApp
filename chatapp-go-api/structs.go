package main

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VisibleUser struct {
	Username string `json:"username"`
}
