package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"reakgo/models"
	"reakgo/utility"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println("Form parsing failed !")
		}
		auth, err := Db.authentication.GetUserByEmail(r.FormValue("email"))
		if auth.Email == r.FormValue("email") {
			data2 := make(map[string]string)
			data2["header"] = "Sorry !"
			data2["type"] = "danger"
			data2["message"] = "This Email already exist!"
			utility.View.ExecuteTemplate(w, "flash", data2)
		} else {
			resp, _ := Db.authentication.InsertData(r.FormValue("full_name"), r.FormValue("email"), r.FormValue("pwd"), r.FormValue("mob"), r.FormValue("address"))
			log.Println("Bishwajeet Samal Insertion")
			log.Println(resp)
			if resp == true {
				data2 := make(map[string]string)
				data2["header"] = "Thanks !"
				data2["type"] = "success"
				data2["message"] = "User Added Successfully !"

				utility.View.ExecuteTemplate(w, "flash", data2)
				utility.RedirectTo(w, r, "/dashboard")
			}
		}

		//match := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(r.FormValue("password")))
		// 	if match != nil {
		// 		// Password match has failed
		// 		data := make(map[string]string)
		// 		data["type"] = "error"
		// 		data["message"] = "Test Message"
		// 		utility.View.ExecuteTemplate(w, "flash", data)
		// 	} else {
		// 		// Password match has been a success
		// 		sessionData := []utility.Session{
		// 			{Key: "username", Value: auth.Email},
		// 			{Key: "type", Value: "user"},
		// 		}
		// 		utility.SessionSet(w, r, sessionData)
		// 		utility.RedirectTo(w, r, utility.Config["appUrl"]+"/dashboard")
		// 	}
	}
	utility.View.ExecuteTemplate(w, "register", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println("Form parsing failed !")
		}
		auth, err := Db.authentication.GetUserByEmail(r.FormValue("email"))
		match := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(r.FormValue("password")))
		if match != nil {
			// Password match has failed
			data := make(map[string]string)
			data["type"] = "error"
			data["message"] = "Test Message"
			utility.View.ExecuteTemplate(w, "flash", data)
		} else {
			// Password match has been a success
			sessionData := []utility.Session{
				{Key: "username", Value: auth.Email},
				{Key: "type", Value: "user"},
			}
			utility.SessionSet(w, r, sessionData)
			utility.RedirectTo(w, r, utility.Config["appUrl"]+"/dashboard")
		}
	}
	utility.View.ExecuteTemplate(w, "index", nil)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println("Form parsing failed !")
		}
		auth, err := Db.authentication.GetUserByEmail(r.FormValue("email"))
		if err != nil {
			// Couldn't find the user in DB, Give fake info
			// Fake sending email delay
			time.Sleep(3)
		} else {
			// User returned successfully, Send email
			fp, err := Db.authentication.ForgotPassword(auth.Id)
			if err != nil {
				log.Println("Token Updation on DB Failed")
				log.Println(err)
			}
			data := make(map[string]string)
			data["token"] = fp
			buf := new(bytes.Buffer)
			err = utility.View.ExecuteTemplate(buf, "emailforgotpassword", data)
			if err != nil {
				log.Println("Template Parsing Failed!")
				log.Println(err)
			} else {
				utility.SendEmail("94b39c058d-2e1b95@inbox.mailtrap.io", "Forgot Password", buf.String())
			}
		}
	}
	utility.View.ExecuteTemplate(w, "forgotpassword", nil)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		token := r.URL.Query().Get("token")
		err := r.ParseForm()
		if err != nil {
			log.Println("Form parsing failed !")
		}
		resp, err := Db.authentication.TokenVerify(token, r.FormValue("password"))
		if resp {
			// Error
		} else {
			// Success
		}
	}
	utility.View.ExecuteTemplate(w, "changepassword", nil)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	resp, err := Db.data.All()
	if err != nil {
		log.Println("Form parsing failed !")
	}
	// Make sure the keys start with Capital letter to ensure export
	data := struct {
		TableData []models.Data
	}{
		resp,
	}
	utility.View.ExecuteTemplate(w, "dashboard", data)
	log.Println(data.TableData[0].Name)
}

func AjaxData(w http.ResponseWriter, r *http.Request) {
	resp, err := Db.data.All()
	if err != nil {
		log.Println("Form parsing failed !")
	}
	json, err := json.Marshal(resp)
	if err != nil {
		log.Println("Form parsing failed !")
	}
	w.Write([]byte(json))
}
