package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	//	"github.com/go-pg/pg/v9"
)

// Vehicle represents data about a record vehicle.
type company struct {
	ID         int    `json:"ID"`
	Name       string `json:"name"`
	AddressLn1 string `json:"addressLn1"`
	AddressLn2 string `json:"addressLn2"`
	City       string `json:"city"`
	State      string `json:"state"`
	Zip        string `json:"zip"`
	Email      string `json:"email"`
	PhoneNum   string `json:"phoneNum"`
	TollFree   string `json:"tollFree"`
	FbSocial   string `json:"fbSocial"`
	TwSocial   string `json:"twSocial"`
	InsSocial  string `json:"insSocial"`
}

func GetCompanyDetails(w http.ResponseWriter, r *http.Request) {
	var result string
	var CompanyArray company
	fmt.Print("getCompanyDetails")
	fmt.Print(dbConnect)
	_, err := dbConnect.Query(&result, `select json_build_object('ID', "I_COMP",'name', "C_COMP_NAME", 'addressLn1', "C_REG_ADD_LN1", 'addressLn2', "C_REG_ADD_LN2", 'city', "C_CITY", 'state', "C_STATE", 'zip', "C_ZIP", 'email', "C_EMAIL", 'phoneNum', "N_PH_NUM", 'tollFree', "N_TOLLFREE", 'fbSocial', "C_FB_SOCIAL", 'twSocial', "C_TW_SOCIAL",'insSocial', "C_INS_SOCIAL" ) result from public."Company"as comp`)
	if err == nil {
		json.Unmarshal([]byte(result), &CompanyArray)
	}
	var resp = map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Company Details",
		"data":    CompanyArray,
	}
	json.NewEncoder(w).Encode(resp)

}
