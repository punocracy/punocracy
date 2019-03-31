package handlers

import (
    "github.com/alvarosness/goodsample/models"
    "github.com/alvarosness/goodsample/libhttp"
    "github.com/gorilla/sessions"
    "html/template"
    "net/http"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")

    sessionStore := r.Context().Value( "sessionStore").(sessions.Store)

    session, _ := sessionStore.Get(r, "goodsample-session")
    currentUser, ok := session.Values["user"].(*models.UserRow)
    if !ok {
        http.Redirect(w, r, "/logout", 302)
        return
    }

    data := struct {
        CurrentUser *models.UserRow
    }{
        currentUser,
    }

    tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")
    if err != nil {
        libhttp.HandleErrorJson(w, err)
        return
    }

    tmpl.Execute(w, data)
}
