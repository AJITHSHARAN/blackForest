package controllers

import (
	"encoding/json"
	"fmt"

	models "example.com/gorestapi/models"

	"net/http"

	"log"
	//	"github.com/go-pg/pg/v9"
)

// Vehicle represents data about a record vehicle.

func GetAllFeedback(w http.ResponseWriter, r *http.Request) {
	var result string
	var FeedbackArray []models.Feedback
	fmt.Print("getAllFeedback")
	fmt.Print(dbConnect)
	_, err := dbConnect.Query(&result, prop.MustGetString("feedback_select"))
	if err == nil {
		json.Unmarshal([]byte(result), &FeedbackArray)
	}
	var resp = map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Feedback Details",
		"data":    FeedbackArray,
	}
	json.NewEncoder(w).Encode(resp)

}

func PostFeedback(w http.ResponseWriter, r *http.Request) {

	postFeed := &models.Feedback{}
	json.NewDecoder(r.Body).Decode(postFeed)
	fmt.Print(postFeed)
	memo := postFeed.Memo
	rating := postFeed.Rating
	custId := postFeed.CustomerID

	var result string

	var insertQuery string

	var rat string
	var cu string

	rat = fmt.Sprint(rating)
	cu = fmt.Sprint(custId)

	insertQuery = prop.MustGetString("feedback_insert") + " VALUES(" + rat + ", '" + memo + "'," + cu + ")"
	log.Printf(insertQuery)

	fmt.Print(cu)

	fmt.Print(insertQuery)
	_, insertError := dbConnect.Query(&result, insertQuery)

	fmt.Print(insertError)

	if insertError != nil {
		var resp = map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		}
		json.NewEncoder(w).Encode(resp)
	}

}
