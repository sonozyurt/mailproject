package main

import (
	"fmt"
	"mailprojesi/datamanagement"
	"net/http"
	"text/template"
	"time"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))

}

func (config *Config) homepage(w http.ResponseWriter, req *http.Request) {
	data := config.addData(req)
	err := tpl.ExecuteTemplate(w, "homepage.gohtml", data)
	config.err(err)
}

func (config *Config) signup(w http.ResponseWriter, req *http.Request) {
	FlashMessage := config.session.PopString(req.Context(), "flash")
	err := tpl.ExecuteTemplate(w, "signup.gohtml", FlashMessage)
	config.err(err)
}

func (config *Config) postSignup(w http.ResponseWriter, req *http.Request) {
	email := req.FormValue("email")
	password := req.FormValue("password")
	checkPassword := req.FormValue("check")
	name := req.FormValue("firstname")
	lastname := req.FormValue("lastname")
	_, err := config.Users.GetByEmail(email)
	if checkPassword != password {
		config.session.Put(req.Context(), "flash", "Passwords dont match")
		http.Redirect(w, req, "/signup", http.StatusSeeOther)
		return
	}
	if err == nil {
		config.session.Put(req.Context(), "flash", "email already taken,you can login")
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	} else {
		user := datamanagement.User{
			Email:     email,
			FirstName: name,
			LastName:  lastname,
			Password:  password,
			Active:    0,
			IsAdmin:   0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err := config.Users.Insert(user)
		if err != nil {
			config.err(err)
			config.session.Put(req.Context(), "flash", "couldn't create user")
			http.Redirect(w, req, "/signup", http.StatusSeeOther)
			return
		}
		msg := message{
			from:    "sendingmail@example.om",
			to:      user.Email,
			message: fmt.Sprintf("acitvation link: localhost%s/activate/?mail=%s", Port, user.Email),
		}
		config.sendMail(msg)
		config.session.Put(req.Context(), "flash", "activation email sent,activate your account then you can login")
		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}

}

func (config *Config) postLogin(w http.ResponseWriter, req *http.Request) {

	email := req.FormValue("email")
	password := req.FormValue("password")
	user, err := config.Users.GetByEmail(email)
	config.err(err)
	if user.Active == 0 {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		msg := message{
			from:    "sendingmail@example.om",
			to:      user.Email,
			message: fmt.Sprintf("acitvation link: localhost%s/activate/?mail=%s", Port, user.Email),
		}
		config.sendMail(msg)
		config.session.Put(req.Context(), "flash", "user in not activated,new activation link has been sent check your email.")
		config.err(err)
		return
	}
	if user.Password != password {
		msg := message{
			from:    "sendingmail@example.om",
			to:      user.Email,
			message: "wrong login attempt",
		}
		config.sendMail(msg)
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		config.session.Put(req.Context(), "flash", "wrong password")

		return
	}
	config.session.Put(req.Context(), "userID", user.ID)
	config.session.Put(req.Context(), "user", user)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (config *Config) login(w http.ResponseWriter, req *http.Request) {
	FlashMessage := config.session.PopString(req.Context(), "flash")
	err := tpl.ExecuteTemplate(w, "login.gohtml", FlashMessage)
	config.err(err)
}

func (config *Config) activate(w http.ResponseWriter, req *http.Request) {
	FlashMessage := config.session.PopString(req.Context(), "flash")
	err := tpl.ExecuteTemplate(w, "homepage.gohtml", FlashMessage)
	config.err(err)
	q := req.URL.Query()
	mail := q.Get("mail")
	user, err := config.Users.GetByEmail(mail)
	config.err(err)
	user.Active = 1
	err = user.Update()
	config.err(err)
	config.session.Put(req.Context(), "flash", "user activated, now you can login")
	msg := message{
		from:    "sendingmail@example.om",
		to:      user.Email,
		message: "account activated",
	}
	config.sendMail(msg)
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func (config *Config) logout(w http.ResponseWriter, req *http.Request) {
	err := config.session.Destroy(req.Context())
	config.err(err)
	config.session.Put(req.Context(), "flash", "logged out")
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
func (config *Config) userOnly(w http.ResponseWriter, req *http.Request) {
	data := config.addData(req)
	if !config.session.Exists(req.Context(), "user") {
		config.session.Put(req.Context(), "flash", "you must be logged in to see this page.")
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(w, "userOnly.gohtml", data)
	config.err(err)
}
