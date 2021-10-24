package models

type CustomerProfile struct {
	tableName struct{} `pg:"CustomerProfile"`
	Id        int      `json:"id" pg:"i_cust"`
	Password  string   `json:"password" pg:"c_pwd"`
	AuthType  string   `json:"auth_type" pg:"c_auth"`
	Email     string   `json:"email" pg:"i_user"`
}
