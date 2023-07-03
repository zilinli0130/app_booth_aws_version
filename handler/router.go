package handler

import (
	"net/http"
	jwt "github.com/form3tech-oss/jwt-go"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func InitRouter() http.Handler {

	// Create jwt middleware before the handlers
	// mySignKey -> global variable inside handler.go
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
        ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
            return []byte(mySigningKey), nil
        },
        SigningMethod: jwt.SigningMethodHS256,
    })

    router := mux.NewRouter()

	// Only authenticated requests are allowed
    router.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST")
  	router.Handle("/checkout", jwtMiddleware.Handler(http.HandlerFunc(checkoutHandler))).Methods("POST")    
	router.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET")
	router.Handle("/app/{id}", jwtMiddleware.Handler(http.HandlerFunc(deleteHandler))).Methods("DELETE")

	// General requests
    router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
    router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

	// Allow CORS
	originsOk := handlers.AllowedOrigins([]string{"*"})
    headersOk := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})

    return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}
