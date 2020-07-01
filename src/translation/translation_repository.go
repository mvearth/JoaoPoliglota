package translation

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	USER    = "postgres"
	PASS    = "banco"
	DBNAME  = "joaopoliglota"
	SSLMODE = "disable"
)

func Connect() *sql.DB {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", USER, PASS, DBNAME, SSLMODE)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}

func InsertTranslation(translation Translation) (bool, error) {
	con := Connect()
	defer con.Close()
	sql := "INSERT INTO translations (idiom, standard_key, translation) VALUES($1, $2, $3) RETURNING translation_id"
	stmt, err := con.Prepare(sql)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(translation.Idiom, translation.StandardKey, translation.Translation)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetTranslation(standardKey, idiom string) (Translation, error) {
	con := Connect()
	defer con.Close()
	sql := "SELECT * FROM translations where idiom = $1 AND standard_key = $2"
	rs, err := con.Query(sql, idiom, standardKey)
	if err != nil {
		return Translation{}, err
	}
	defer rs.Close()
	var translation Translation
	for rs.Next() {
		err := rs.Scan(&translation.ID, &translation.Idiom, &translation.StandardKey, &translation.Translation)
		if err != nil {
			return translation, err
		}
	}
	return translation, nil
}

func TestConnection() {
	con := Connect()
	defer con.Close()
	err := con.Ping()
	if err != nil {
		fmt.Errorf("%s", err.Error())
		return
	}

	fmt.Println("Database connected!")
}
