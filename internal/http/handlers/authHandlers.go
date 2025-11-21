package handlers

import (
	"net/http"

	"github.com/medidew/ApplicationTracker/internal/auth"
)

func (app *App) Register(response_writer http.ResponseWriter, request *http.Request) {
	// TODO: Validate email and username format

	password := request.FormValue("password")

	argon2auth := &auth.Argon2Auth{}
	argon2auth.SetDefaults()
	hashed_password := argon2auth.HashPassword([]byte(password))

	err := app.DB.CreateUser(
		request.FormValue("email"),
		request.FormValue("username"),
		argon2auth,
		hashed_password,
	)
	if err != nil {
		http.Error(response_writer, "DB insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response_writer.WriteHeader(http.StatusCreated)
	_, err = response_writer.Write([]byte("User created successfully"))
	if err != nil {
		http.Error(response_writer, "failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App) Login(response_writer http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	app.Logger.Info("user login request: " + username)
	password := request.FormValue("password")

	hashed_password, err := app.DB.GetUserHashedPassword(username)
	if err != nil {
		http.Error(response_writer, "failed to get user hashed password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	argon2auth, err := app.DB.GetUserArgon2Auth(username)
	if err != nil {
		http.Error(response_writer, "failed to get user argon2 auth: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !argon2auth.VerifyPassword([]byte(password), hashed_password) {
		http.Error(response_writer, "invalid username or password", http.StatusUnauthorized)
		return
	}

	app.SessionManager.RenewToken(request.Context())
	app.SessionManager.Put(request.Context(), "username", username)

	_, err = response_writer.Write([]byte("Login successful"))
	if err != nil {
		http.Error(response_writer, "failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *App) Logout(response_writer http.ResponseWriter, request *http.Request) {
	app.SessionManager.Destroy(request.Context())

	_, err := response_writer.Write([]byte("Logout successful"))
	if err != nil {
		http.Error(response_writer, "failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}