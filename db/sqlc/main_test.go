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


func TestMain(m *testing.M){
conn,error:= sql.Open(dbDriver,dbSource)
if error!=nil{
	log.Fatal("Failed to connect to db",error)
}
testQueries=New(conn)
os.Exit(m.Run())
}