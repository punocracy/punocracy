package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

// GetSignup generates the user signup page
func GetSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/users/users-external.html.tmpl", "templates/users/signup.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, nil)
}

// PostSignup reads the new user's credentials and stores them in the user DB if they are valid
// After signing up, the user is autmatically logged in.
func PostSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	db := r.Context().Value("db").(*sqlx.DB)

	username := r.FormValue("Username")
	email := r.FormValue("Email")
	password := r.FormValue("Password")
	passwordAgain := r.FormValue("PasswordAgain")

	_, err := models.NewUser(db).Signup(nil, username, email, password, passwordAgain)
	if err != nil {
		// TODO: Redirect to Login maybe with an error message
		logrus.Infoln(err)
		libhttp.HandleErrorJson(w, err)
		return
	}

	PostLogin(w, r)
}

// GetLoginWithoutSession generates the login page without checking if an existing user has already logged in
func GetLoginWithoutSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/users/users-external.html.tmpl", "templates/users/login.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, nil)
}

// GetLogin first checks if a user is already logged in.
// If the user is already logged in, the user is redirected to the home page.
// If not GetLoginWithoutSession is called
func GetLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")

	currentUserInterface := session.Values["user"]
	if currentUserInterface != nil {
		http.Redirect(w, r, "/now", 302)
		return
	}

	GetLoginWithoutSession(w, r)
}

// PostLogin handles user authentication.
// If the user has an account in the system, he/she is redirected to the home page
// If the user used the wrong credentials, redirect them to the login page with an error message.
func PostLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	db := r.Context().Value("db").(*sqlx.DB)
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	username := r.FormValue("Username")
	password := r.FormValue("Password")

	u := models.NewUser(db)

	user, err := u.GetUserByUsernameAndPassword(nil, username, password)
	if err != nil {
		logrus.Errorln(err.Error())
		http.Redirect(w, r, "login/", http.StatusFound)
	}

	session, _ := sessionStore.Get(r, "punocracy-session")
	session.Values["user"] = user

	err = session.Save(r, w)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	http.Redirect(w, r, "/now", 302)
}

// GetLogout deletes the current user from the session and redirects to the main page
func GetLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")

	delete(session.Values, "user")
	session.Save(r, w)

	http.Redirect(w, r, "/now", 302)
}

// PostPutDeleteUsersID will redirect to either the PutUsersID or the DeleteUsersID handlers depending on the typpe of request
func PostPutDeleteUsersID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	method := r.FormValue("_method")
	if method == "" || strings.ToLower(method) == "post" || strings.ToLower(method) == "put" {
		PutUsersID(w, r)
	} else if strings.ToLower(method) == "delete" {
		DeleteUsersID(w, r)
	}
}

// PutUsersID updates the user data
func PutUsersID(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromPath(w, r)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	db := r.Context().Value("db").(*sqlx.DB)

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")

	currentUser := session.Values["user"].(*models.UserRow)

	if currentUser.ID != userID {
		err := errors.New("modifying other user is not allowed")
		libhttp.HandleErrorJson(w, err)
		return
	}

	password := r.FormValue("Password")
	passwordAgain := r.FormValue("PasswordAgain")

	u := models.NewUser(db)

	currentUser, err = u.UpdateUsernameAndPasswordByID(nil, currentUser.ID, currentUser.Username, password, passwordAgain)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	// Update currentUser stored in session.
	session.Values["user"] = currentUser
	err = session.Save(r, w)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	http.Redirect(w, r, "/now", 302)
}

// DeleteUsersID is not implemented
func DeleteUsersID(w http.ResponseWriter, r *http.Request) {
	err := errors.New("delete method is not implemented")
	libhttp.HandleErrorJson(w, err)
	return
}
