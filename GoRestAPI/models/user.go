package models

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	tableName struct{} `pg:"user_auth"`
	Email     string   `json:"email" pg:"email"`
	Password  string   `json:"password" pg:"password"`
	Hash      string   `json:"hash" pg:"hash"`
	Active    int      `json:"active" pg:"active"`
}
