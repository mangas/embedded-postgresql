package embeded_postgresql

import (
	"testing"
	"os"
	"path"

	_ "github.com/lib/pq"
	"fmt"
	"database/sql"
)

func TestCreateDatabase(t *testing.T) {

	sConfig := StartupConfig{true, "9.6.5-1"}
	dConfig := DBConfig{46782, "postgres"}

	os.Remove(path.Join(os.Getenv("HOME"), ".postgres-embedded", "9.6.5-1"))

	rc := StartPostgres(sConfig, dConfig)

	connStr := fmt.Sprintf("user=postgres dbname=postgres sslmode=disable port=%v host=127.0.0.1", dConfig.Port)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		t.Errorf("Connection should not have errors but had: %v", err)
	}

	rows := db.QueryRow("select 100")

	var i int

	err = rows.Scan(&i)

	if err != nil {
		t.Errorf("Scan should not have errors but had: %v", err)
	}

	fmt.Println(fmt.Sprintf("i = %v", i))


	defer func() {
		recover()
		StopPostGres(rc)
	}()
}