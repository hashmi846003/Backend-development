package db
import (
"testing"
 "database/sql"
 "log"
_ "github.com/lib/pq"
 "os"
)

const (
	dbDriver= "postgres"
	dbSource= "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)
var testQueries *Queries
var testDB *sql.DB


func TestMain(m *testing.M){
	var err error
	
testDB,err= sql.Open(dbDriver,dbSource)
if err!=nil{
	log.Fatal("Failed to connect to db",err)
}
testQueries=New(testDB)
os.Exit(m.Run())
}