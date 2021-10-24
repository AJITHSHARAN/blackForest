package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/gorestapi/models"
	"github.com/go-pg/pg/v9"
	"github.com/magiconair/properties"
)

// Vehicle represents data about a record vehicle.

// INITIALIZE DB CONNECTION (TO AVOID TOO MANY CONNECTION)
var dbConnect *pg.DB
var prop *properties.Properties

func InitiateDB(db *pg.DB) {
	dbConnect = db
}

func InitiateProp(p *properties.Properties) {
	prop = p
	fmt.Print(prop.MustGetString("vehicle_select"))
}

func GetAllVehicles(w http.ResponseWriter, r *http.Request) {

	var result string
	var VehicleArray []models.Vehicle
	//fmt.Println(prop.MustGetString("host"))
	fmt.Print(dbConnect)
	_, err := dbConnect.Query(&result, prop.MustGetString("vehicle_select"))
	if err == nil {
		json.Unmarshal([]byte(result), &VehicleArray)
	}
	//fmt.Print(result)

	fmt.Println(VehicleArray)
	fmt.Println(len(VehicleArray))

	fmt.Print("success")
	var resp = map[string]interface{}{
		"status":  http.StatusOK,
		"message": "All Vehicles",
		"data":    VehicleArray,
	}
	json.NewEncoder(w).Encode(resp)
	return
}
