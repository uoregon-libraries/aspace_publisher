package db

import "os"

type DBPool struct {
  Pguser string
  Pgpass string
  Pgdb string
  Pgaddress string
}

func (db *DBPool)Init(){
  db.Pguser = os.Getenv("POSTGRES_USER")
  db.Pgpass = os.Getenv("POSTGRES_PASSWORD")
  db.Pgdb = os.Getenv("POSTGRES_DB")
  db.Pgaddress = os.Getenv("DATABASE_URL")
}
