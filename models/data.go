package models

import (
    "database/sql"
    "reakgo/utility"
    "log"
)

type Data struct {
    Id int32
    Name string
}

type DataModel struct {
    DB *sql.DB
}

func (data DataModel) All () ([]Data, error) {
    var allRows []Data

    rows,err := utility.Db.Query("SELECT * FROM data")
    if err != nil {
        log.Println("Query failed")
        log.Println(err)
    }
    defer rows.Close()
    for rows.Next(){
        var singleRow Data
        err := rows.Scan(&singleRow.Id, &singleRow.Name)
        if err != nil {
            log.Println("Query Scan failed")
            log.Println(err)
        }
        allRows = append(allRows, singleRow)
    }
    if err = rows.Err(); err != nil{
        log.Println("Couldn't fetch from data")
        log.Println(err)
    }

    return allRows, err
}
