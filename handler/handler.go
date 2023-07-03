package handler

import (
	"appstore/model"
	"appstore/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/pborman/uuid"
	"regexp"
	jwt "github.com/form3tech-oss/jwt-go"
	"time"
	"github.com/gorilla/mux"
)

var (
	mySigningKey = []byte("secret")
)

/*
* App item handler
*/
func uploadHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Received one upload request")

	// Read username from Jwt token payload 
	user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]

	// Parse the request body
	app := model.App{
		Id:          uuid.New(),
		User:        username.(string),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	fmt.Printf("%v,%T", price, price)
	if err != nil {
		fmt.Println(err)
	}
   	app.Price = int(price * 100.0)

	// Get the media file
	file, _, err := r.FormFile("media_file")
	if err != nil {
		http.Error(w, "Media file is not available", http.StatusBadRequest)
		fmt.Printf("Media file is not available %v\n", err)
		return
	}

	// Save the App item
	err = service.SaveApp(&app, file)
	if err != nil {
		http.Error(w, "Failed to save app to backend", http.StatusInternalServerError)
		fmt.Printf("Failed to save app to backend %v\n", err)
		return
	}
   	fmt.Println("App is saved successfully.")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one search request")
	w.Header().Set("Content-Type", "application/json")
	title := r.URL.Query().Get("title")
	description := r.URL.Query().Get("description")
 
	var apps []model.App
	var err error
	apps, err = service.SearchApps(title, description)
	if err != nil {
		http.Error(w, "Failed to read Apps from backend", http.StatusInternalServerError)
		return
	}
 
	js, err := json.Marshal(apps)
	if err != nil {
		http.Error(w, "Failed to parse Apps into JSON format", http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one checkout request")
	w.Header().Set("Content-Type", "text/plain")
 
	appID := r.FormValue("appID")
	s, err := service.CheckoutApp(r.Header.Get("Origin"), appID)
	if err != nil {
		fmt.Println("Checkout failed.")
		w.Write([]byte(err.Error()))
		return
	}
 
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s.URL))
 
	fmt.Println("Checkout process started!")
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for delete")

    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"].(string)
    id := mux.Vars(r)["id"]

	fmt.Printf("Deleting app: %v\n", id)
	fmt.Printf("User name: %v\n", username)

    if err := service.DeleteApp(id, username); err != nil {
        http.Error(w, "Failed to delete app from backend", http.StatusInternalServerError)
        fmt.Printf("Failed to delete app from backend %v\n", err)
        return
    }
    fmt.Println("App is deleted successfully")
}


/*
* User handler
*/
func signupHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signup request")
    w.Header().Set("Content-Type", "text/plain")

    decoder := json.NewDecoder(r.Body)
    var user model.User
    if err := decoder.Decode(&user); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }

	// Regular expression to validate username
    if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        fmt.Printf("Invalid username or password\n")
        return
    }

    success, err := service.AddUser(&user)
    if err != nil {
        http.Error(w, "Failed to save user to Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to save user to Elasticsearch %v\n", err)
        return
    }

    if !success {
        http.Error(w, "User already exists", http.StatusBadRequest)
        fmt.Println("User already exists")
        return
    }
    fmt.Printf("User added successfully: %s.\n", user.Username)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signin request")
    w.Header().Set("Content-Type", "text/plain")

    // Get User information from client
    decoder := json.NewDecoder(r.Body)
    var user model.User
    if err := decoder.Decode(&user); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }

	// Check if user already exists
    exists, err := service.CheckUser(user.Username, user.Password)
    if err != nil {
        http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read user from Elasticsearch %v\n", err)
        return
    }

    if !exists {
        http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized)
        fmt.Printf("User doesn't exists or wrong password\n")
        return
    }

	// Create JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

	// Generate encoded token and send it as response
    tokenString, err := token.SignedString(mySigningKey)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        fmt.Printf("Failed to generate token %v\n", err)
        return
    }
    w.Write([]byte(tokenString))
}


 
 