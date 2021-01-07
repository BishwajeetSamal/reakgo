package utility

import (
    "net/http"
    "html/template"
    "database/sql"
    "log"
    "github.com/gorilla/sessions"
    mail "gopkg.in/mail.v2"
    "crypto/tls"
    "strconv"
)

// Template Pool
var View *template.Template
// Session Store
var Store *sessions.CookieStore
// DB Connections
var Db *sql.DB

type Session struct {
    Key string
    Value string
}

var Config = map[string]string{
    "appUrl": "http://localhost:4000",
    "smtpHost": "smtp.mailtrap.io",
    "smtpEmail": "test@test.com",
    "smtpUser": "1bd57713f81f6d",
    "smtpPassword": "7ec4deb203e5e1",
    "smtpPort": "2525",
}

func RedirectTo(w http.ResponseWriter, r *http.Request, url string) {
    http.Redirect(w, r, url, 302)
}

func SessionSet(w http.ResponseWriter, r *http.Request, data []Session) {
    session, _ := Store.Get(r, "session-name")
    // Set some session values.
    for _, dataSingle := range data {
        session.Values[dataSingle.Key] = dataSingle.Value
    }
    // Save it before we write to the response/return from the handler.
    err := session.Save(r, w)
    log.Println(err)
}

func SessionGet(r *http.Request, key string) interface{} {
    session, _ := Store.Get(r, "session-name")
    // Set some session values.
    return session.Values[key]
}


func CheckACL(w http.ResponseWriter, r *http.Request, minLevel int) bool {
    userType := SessionGet(r, "type")
    var level int = 0
    switch(userType){
    case "user":
        level = 1
    case "admin":
        level = 2
    default:
        level = 0
    }
    if(level >= minLevel){
        return true
    } else {
        RedirectTo(w, r, Config["appUrl"]+"/login")
        return false
    }
}

func SendEmail(to string, subject string, body string){
    m := mail.NewMessage()
    // Set E-Mail sender
    m.SetHeader("From", Config["smtpEmail"])
    // Set E-Mail receivers
    m.SetHeader("To", to)
    // Set E-Mail subject
    m.SetHeader("Subject", subject)
    // Set E-Mail body. You can set plain text or html with text/html
    m.SetBody("text/html", body)
    // Settings for SMTP server
    port, err := strconv.Atoi(Config["smtpPort"])
    if (err != nil){
        // Set default port if we get something else
        port = 587
    }
    d := mail.NewDialer(Config["smtpHost"], port, Config["smtpUser"], Config["smtpPassword"])
    // This is only needed when SSL/TLS certificate is not valid on server.
    // In production this should be set to false.
    d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
    // Now send E-Mail
    if err := d.DialAndSend(m); err != nil {
        log.Println(err)
    }
    return
}
