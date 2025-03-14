package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})

	http.HandleFunc("/login", login)

	http.HandleFunc("/jwt", jwt)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Login successful!")
}

func jwt(w http.ResponseWriter, r *http.Request) {

	header := "{\"alg\":\"HS256\",\"typ\":\"JWT\"}"
	payload := "{\"sub\":\"willams\",\"iat\":1516239022}"

	signedJWT := SIGJWT(header, payload, "123")
	

	
	jwt := fmt.Sprintf("%s.%s.%s", signedJWT.Header, signedJWT.Payload, signedJWT.Signature)
	response := Token{AccessToken: jwt}
		
	w.Header().Set("Content-Type", "application/json")
	
	json.NewEncoder(w).Encode(response)
}
