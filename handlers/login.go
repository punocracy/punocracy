package handlers

// GetLogin is responsible for logins
// func GetLogin(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")

// 	t, _ := template.ParseFiles("templates/dashboard.html.tmpl", "templates/login.html.tmpl")
// 	t.Execute(w, nil)
// }

// // PostLogin is responsible for logins
// func PostLogin(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")

// 	r.ParseForm()

// 	username := r.Form["username"]
// 	password := r.Form["password"]

// 	// TODO: Authenticate user
// 	logrus.Infoln(username, password)

// 	// TODO: Redirect if sucessful authentication
// 	http.Redirect(w, r, "/", http.StatusFound)
// }
