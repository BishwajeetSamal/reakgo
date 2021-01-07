package controllers

import (
    "reakgo/models"
    "reakgo/utility"
)

type Env struct {
    authentication interface {
        GetUserByEmail(email string) (models.Authentication, error)
        ForgotPassword(id int32) (string, error)
        TokenVerify(token string, newPassword string) (bool, error)
    }
    data interface {
        All() ([]models.Data, error)
    }
}

var Db *Env

func init(){
    // Initialize DB
    Db = &Env{
        authentication: models.AuthenticationModel{DB: utility.Db},
        data: models.DataModel{DB: utility.Db},
    }
}
