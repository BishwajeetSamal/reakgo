package controllers

import (
    "reakgo/models"
    "reakgo/utility"
)

type Env struct {
    authentication interface {
        GetUserByEmail(email string) (models.Authentication, error)
        ForgotPassword(id int32) (string, error)
        TokenVerify(token string) (bool, error)
    }
}

var Db *Env

func init(){
    // Initialize DB
    Db = &Env{
        authentication: models.AuthenticationModel{DB: utility.Db},
    }
}
