package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid"
)

var (
	dbx *sqlx.DB
)

type ApplicationUser struct {
	ID           string    `json:"id" db:"id"`
	ScreenID     string    `json:"screen_id" db:"screen_id"`
	ScreenName   string    `json:"screen_name" db:"screen_name"`
	LastModified time.Time `json:"last_modified" db:"last_modified"`
	Height       string    `json:"height" db:"height"`
}

type DataStatus struct {
	RecordDate   string    `json:"record_date" db:"record_date"`
	UserID       int       `json:"user_id" db:"user_id"`
	Weight       int       `json:"weight" db:"weight"`
	Lastmodified time.Time `json:"lastmodified" db:"lastmodified"`
	Bpf          int       `json:"bpf" db:"bpf"`
}

type UserSystemInfo struct {
	ID             string    `json:"id" db:"id"`
	MainAddress    string    `json:"main_address" db:"main_address"`
	Passhash       string    `json:"passhash" db:"passhash"`
	LoginFailed    int       `json:"login_failed" db:"login_failed"`
	IsLocked       bool      `json:"is_locked" db:"is_locked"`
	LastmodifiedAt time.Time `json:"lastmodified_at" db:"lastmodified_at"`
}

func main() {
	//テストコミット1

	// テストコメント



   // テストコメント4


	
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("failed to reed DB port number from an environment variable MYSWL_PORT.\nError: %s", err.Error())
	}

	user := os.Getenv("MYSQL_USER")
	if user == "" {
		user = "punipi"
	}
	dbname := os.Getenv("MYSQL_DBNAME")
	if dbname == "" {
		dbname = "mydb"
	}
	password := os.Getenv("MYSQL_PASSWORD")
	if password == "" {
		password = "punipi"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user,
		password,
		host,
		port,
		dbname,
	)

	dbx, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %s.", err.Error())
	}
	defer dbx.Close()

	http.HandleFunc("/get_user", getUser)
	http.HandleFunc("/set_user", setUser)
	http.ListenAndServe(":3333", nil)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")

	var user ApplicationUser
	err := dbx.Get(&user,
		"SELECT * FROM `application_user` WHERE `id` = ?",
		id,
	)
	if err != nil {
		log.Print(err)
		return
	}

	fmt.Printf(
		"id: %s\nscreen_id: %s\nscreen_name: %s\nlast_modifed: %s\nheight: %s\n",
		user.ID,
		user.ScreenID,
		user.ScreenName,
		user.LastModified,
		user.Height,
	)
	return
}

func setUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	screenID := query.Get("screen_id")
	screenName := query.Get("screen_name")
	lastModified := query.Get("last_modified")
	height := query.Get("height")

	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	tx := dbx.MustBegin()
	_, err := tx.Exec("INSERT INTO `application_user` (`id`, `screen_id`, `screen_name`, `last_modified`, `height`) VALUES (?, ?, ?, ?, ?)",
		id.String(),
		screenID,
		screenName,
		lastModified,
		height,
	)

	if err != nil {
		log.Print(err)
		tx.Rollback()
		return
	}
	tx.Commit()
}
