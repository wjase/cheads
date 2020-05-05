package auth

import (
	//...
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	// import the jwt-go library
	"github.com/dgrijalva/jwt-go"
	//...
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// Credentials Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Claims Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//AuthorisedFn wrapped fn taking the unpacked claims for use
type AuthorisedFn func(claims *Claims, w http.ResponseWriter, r *http.Request)

// Signin Create the Signin handler
func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	fmt.Printf("Create new user token for %s\n", creds.Username)

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   tokenString,
	// 	Expires: expirationTime,
	// })

	w.Write([]byte(tokenString))

}

// AuthnError error for authentication
type AuthnError struct {
	Msg    string
	Status int
}

func (e *AuthnError) Error() string {
	return fmt.Sprintf("Authentication error: %d -%s", e.Status, e.Msg)
}

// GetTokenClaimsFromCookie returns the claims from the request cookie
func GetTokenClaimsFromCookie(r *http.Request) (*Claims, error) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, nil
		}
		// For any other type of error, return a bad request status
		return nil, &AuthnError{err.Error(), http.StatusBadRequest}
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	return ClaimsFromToken(tknStr)

}

// GetTokenClaimsFromParam returns the claims from the request cookie
func GetTokenClaimsFromParam(r *http.Request) (*Claims, error) {
	// We can obtain the token from the requests token param
	if token, ok := r.URL.Query()["token"]; ok {
		tknStr := token[0]

		return ClaimsFromToken(tknStr)
	}
	fmt.Printf("did not get claim")

	// Get the JWT string from the param
	return nil, nil

}

// ClaimsFromToken - parses token string
func ClaimsFromToken(tknStr string) (*Claims, error) {
	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, &AuthnError{"Signature Invalid", http.StatusUnauthorized}
		}
		return nil, &AuthnError{err.Error(), http.StatusBadRequest}
	}
	if !tkn.Valid {
		return nil, &AuthnError{"Invalid token", http.StatusUnauthorized}
	}

	return claims, nil
}

// GetTokenClaimsFromRequest gets token from either cookie or param
func GetTokenClaimsFromRequest(r *http.Request) (*Claims, error) {
	claims, err := GetTokenClaimsFromCookie(r)
	if err != nil {
		return nil, err
	}
	if claims == nil {
		claims, err = GetTokenClaimsFromParam(r)
	}
	if err != nil {
		if err != nil {
			return nil, err
		}
	}

	if claims == nil {
		return nil, &AuthnError{"Invalid token", http.StatusUnauthorized}
	}

	return claims, nil
}

// ContextKey a string key
type ContextKey string

// ContextClaimsKey string constant for accessing claims from the request context
const ContextClaimsKey ContextKey = "claims"

//Authorised wraps a http handler in a jwt check
func Authorised(fn func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := GetTokenClaimsFromRequest(r)
		if err != nil {
			if err, ok := err.(*AuthnError); ok {
				fmt.Println(err)
				w.WriteHeader(err.Status)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), ContextClaimsKey, claims)
		fn(w, r.WithContext(ctx))
	}
}

// Refresh renews a token
func Refresh(w http.ResponseWriter, r *http.Request) {
	claims, err := GetTokenClaimsFromCookie(r)
	if err != nil {
		if err, ok := err.(*AuthnError); ok {
			fmt.Println(err)
			w.WriteHeader(err.Status)
			return
		}
		fmt.Println(err)
		return
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
