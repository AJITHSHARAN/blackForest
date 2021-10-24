package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"example.com/gorestapi/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var bigstr string
	user := &models.UserRequest{}
	json.NewDecoder(r.Body).Decode(user)
	userObj := []models.User{}
	p, _ := rand.Prime(rand.Reader, 64)
	err := dbConnect.Model(&userObj).Where("id = " + user.Email + " and active = 1").Select()
	if err == nil {
		if len(userObj) == 1 {
			var resp = map[string]interface{}{
				"message": "already logged in user",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
	bigstr = p.String()
	SendEmail(user.Email, bigstr)
	var resp = map[string]interface{}{
		"message": "OTP was sent to that Email successfully",
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	insertiIntoUserAuth(user.Email, string(pass), bigstr)
	json.NewEncoder(w).Encode(resp)
}

func insertiIntoUserAuth(email string, passwd string, token string) {
	user := []models.User{}
	NewUser := &models.User{}
	NewUser.Active = 0
	NewUser.Email = email
	NewUser.Password = passwd
	NewUser.Hash = token
	fmt.Print(dbConnect)
	err := dbConnect.Model(&user).Where("email = ?", email).Select()
	if err == nil {
		if len(user) >= 1 {
			NewUser1 := &models.User{}
			NewUser1 = &user[0]
			_, err := dbConnect.Model(NewUser1).Column("hash").Where("email =?", email).Update()
			fmt.Println(err)
		} else {
			_, err := dbConnect.Model(NewUser).Insert()
			fmt.Println(err)
		}
	}
}

func OtpAuthenticate(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	userArray := []models.User{}
	//var result string
	q := r.URL.Query()
	hash := q["token"][0]
	email := q["email"][0]
	user.Active = 1
	custProfile := &models.CustomerProfile{}
	err := dbConnect.Model(&userArray).Where("email = ?", email).Select()
	if err == nil {
		if userArray[0].Active == 1 {
			var resp = map[string]interface{}{
				"message": "This Email is already verified",
			}
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			if userArray[0].Hash == hash {
				_, err := dbConnect.Model(user).Column("active").Where("email =?", email).Update()
				fmt.Println(err)
				custProfile.Email = userArray[0].Email
				custProfile.Password = userArray[0].Password
				custProfile.AuthType = "General"
				_, err = dbConnect.Model(custProfile).Insert()
				fmt.Println(err)
				if err == nil {
					var resp = map[string]interface{}{
						"message": "This Email is verified",
					}
					json.NewEncoder(w).Encode(resp)
					return
				} else {
					var resp = map[string]interface{}{
						"message": "something went wrong",
					}
					json.NewEncoder(w).Encode(resp)
					return
				}

			} else {
				var resp = map[string]interface{}{
					"message": "The link is expired",
				}
				json.NewEncoder(w).Encode(resp)
			}
		}

	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	user := &models.UserRequest{}
	json.NewDecoder(r.Body).Decode(user)
	var custProfile []models.CustomerProfile
	err := dbConnect.Model(&custProfile).Where("i_user=?", user.Email).Select()
	if err == nil {
		errf := bcrypt.CompareHashAndPassword([]byte(custProfile[0].Password), []byte(user.Password))
		if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
			fmt.Println("password not match")
			var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
			json.NewEncoder(w).Encode(resp)
			return
		}
		expiresAt := time.Now().Add(time.Minute * 100000).Unix()
		tk := &models.Token{
			Email: user.Email,
			StandardClaims: &jwt.StandardClaims{
				ExpiresAt: expiresAt,
			},
		}

		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

		tokenString, error := token.SignedString([]byte("secret"))
		if error != nil {
			fmt.Println(error)
		}

		var resp = map[string]interface{}{"status": false, "message": "logged in"}
		resp["token"] = tokenString //Store the token in the response
		resp["user"] = custProfile
		json.NewEncoder(w).Encode(resp)
		return

	} else {
		fmt.Println(err)
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		json.NewEncoder(w).Encode(resp)
		return
	}

}

func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func SendEmail(email string, token string) {

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	var message gmail.Message
	temp := []byte("From: 'me'\r\n" +
		"reply-to: sender@gmail.com\r\n" +
		"To:  " + email + "\r\n" +
		"Subject: " + "Regarding OTP" + "\r\n" +
		"\r\n" + "Hi,\r\n Welcome to BlackForest, The verification link is " + prop.MustGetString("host_url") + "/otp/authenticate?email=" + email + "&token=" + token)

	message.Raw = base64.StdEncoding.EncodeToString(temp)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Fatalf("Unable to send. %v", err)
	}
}

func FetchUsers(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user")
	fmt.Print(user)
}
