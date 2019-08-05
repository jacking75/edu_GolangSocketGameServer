package protocol

import (
"database/sql"
"fmt"
_ "github.com/go-sql-driver/mysql"
"log"
"sync"
)

var (
	once sync.Once
	dbConnect *sql.DB
)

func GetInstance() *sql.DB {
	once.Do(func() {
		if dbConnect == nil {
			var err error
			dbConnect, err = sql.Open("mysql", "suyeon:+saJ;65Q%y++@tcp(127.0.0.1:3306)/WOR")

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("inside singleton")
		}
	})
	return dbConnect
}

func Query(sql string) error {
	db := GetInstance()
	_, err := db.Exec(sql)

	return err
}

func Select(ID []byte) (int, int, int, int) {
	db := GetInstance()

	query := fmt.Sprint("SELECT Dollar, Plays, Win, Lose FROM BaccaratLogin WHERE ID = '",string(ID),"'")
	fmt.Println(query)

	var Dollar int
	var Plays int
	var Win int
	var Lose int
	_ = db.QueryRow(query).Scan(&Dollar, &Plays, &Win, &Lose)

	return Dollar, Plays, Win, Lose
}


func Insert(ID []byte, PW []byte) {
	db := GetInstance()

	query := fmt.Sprint("INSERT INTO BaccaratLogin VALUES ('",string(ID),"','",string(PW),"', 100, 0, 0, 0)")
	fmt.Println(query)

	_, _ = db.Exec(query)
}

func Update(setting string, ID []byte) {
	db := GetInstance()

	query := fmt.Sprint("UPDATE BaccaratLogin SET ",setting," WHERE ID = '",string(ID),"'")
	fmt.Println(query)

	_, _ = db.Exec(query)
}