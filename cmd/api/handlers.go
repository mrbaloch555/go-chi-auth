package main

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mrbaloch555/go-chi-auth/common"
	"github.com/mrbaloch555/go-chi-auth/models"
)

func (app *Config) Login(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.UserModel.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := app.Models.UserModel.PasswordMatches(user, requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	token, err := app.createToken(user, w)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "logged in successfully",
		Data:    map[string]any{"user": user, "token": token},
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) Register(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Active    int    `json:"active"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user := models.User{
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Email:     requestPayload.Email,
		Password:  requestPayload.Password,
		Active:    requestPayload.Active,
	}

	_, err = app.Models.UserModel.Insert(&user)

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "registered successfully",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) createToken(user *models.User, w http.ResponseWriter) (*common.JWTOutput, error) {
	expirationTimeConv, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	expirationTime := time.Now().Add(time.Duration(expirationTimeConv) * time.Hour)
	claims := &common.Claims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return nil, err
	}

	jwtOutput := common.JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}

	return &jwtOutput, nil
}
