package controllers

import (
    "net/http"
    "reakgo/utility"
    "log"
    "time"
    "bytes"
    "golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
    if (r.Method == "POST"){
        err := r.ParseForm()
        if (err != nil){
            log.Println("Form parsing failed !")
        }
        auth, err := Db.authentication.GetUserByEmail(r.FormValue("email"))
        match := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(r.FormValue("password")))
        if(match != nil){
            // Password match has failed
            data := make(map[string]string)
            data["type"] = "error"
            data["message"] = "Test Message"
            utility.View.ExecuteTemplate(w, "flash", data)
        } else {
            // Password match has been a success
            sessionData := []utility.Session{
                {Key:"username", Value:auth.Email},
                {Key:"type", Value:"user"},
            }
            utility.SessionSet(w, r, sessionData)
            utility.RedirectTo(w, r, utility.Config["appUrl"]+"/dashboard")
        }
    }
    utility.View.ExecuteTemplate(w, "index", nil)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
    if(r.Method == "POST"){
        err := r.ParseForm()
        if (err != nil){
            log.Println("Form parsing failed !")
        }
        auth, err := Db.authentication.GetUserByEmail(r.FormValue("email"))
        if (err != nil){
            // Couldn't find the user in DB, Give fake info
            // Fake sending email delay
            time.Sleep(3)
        } else {
            // User returned successfully, Send email
            fp, err := Db.authentication.ForgotPassword(auth.Id)
            if(err != nil){
                log.Println("Token Updation on DB Failed")
                log.Println(err)
            }
            data := make(map[string]string)
            data["token"] = fp
            buf := new(bytes.Buffer)
            err = utility.View.ExecuteTemplate(buf, "emailforgotpassword", data)
            if(err != nil){
                log.Println("Template Parsing Failed!")
                log.Println(err)
            } else {
                utility.SendEmail("94b39c058d-2e1b95@inbox.mailtrap.io", "Forgot Password", buf.String())
            }
        }
        log.Println(auth)
    } else {
        if(r.URL.Query().Get("token") != ""){
            log.Println("token encountered")
        }
    }
    utility.View.ExecuteTemplate(w, "forgotpassword", nil)
}


func Dashboard(w http.ResponseWriter, r *http.Request) {
    Db.authentication.TokenVerify("zNF-vce-NSqhkNcUdrbUcnYxmF8Um8NPk-spPSox2UiMydjlzxFGO3T9iS3M")
    utility.View.ExecuteTemplate(w, "dashboard", nil)
}
