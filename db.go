package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
	"time"
)

// Struct representing user of this website
type User struct {

	Username string
	Password string
	Admin bool
	JustSigned bool
}

// Struct representing news article of this website
type Article struct {

	Author string
	Header string
	Text string
	Importance int
	Date string
}

// Struct representing file's downloadable link of this website
type File struct {

	Uploader string
	Description string
	Url string
	Date string
}

type Message struct {

	Username string
	Message  string
	Date string
}

//   Private variable, used only internally from
//  this file
var db *sql.DB

// Open and check connection to database
func init() {

	var err error

	db, err = sql.Open("mysql", "Hayk:sqlisfun@/TechCompany")

	if err != nil {

		log.Fatal(err)
		return
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {

		log.Fatal(err)
		return
	}
}

// Get user by username and password
// If no user available by this username and password return nil
func GetUser(username, password string) *User {

	if isInjected(username) || isInjected(password) {

		return nil
	}

	query := "SELECT * FROM `users` WHERE BINARY `username` = ? AND BINARY `pass` = ?;"

	raws, err := db.Query(query, username, password)

	if err != nil {

		log.Println(err)
		return nil
	}
	defer raws.Close()

	var user *User = new(User)

	if raws.Next() {

		err := raws.Scan(
			&user.Username,
			&user.Password,
			&user.Admin,
			)

		if err != nil {

			log.Println(err)
			return nil
		}

		user.JustSigned = true

		return user
	}

	return nil
}

func ChangeUserInDB(username, password string, admin bool) error {

	if isInjected(username) || isInjected(password) {

		return securityError{ "Security error" }
	}

	adm := 0

	if admin == true {

		adm = 1
	}

	query := "UPDATE `users` set `pass` = ?, `admin` = ? WHERE `username` = ?;"

	_, err := db.Exec(query, password, strconv.Itoa(adm), username)
	if err != nil {

		return err
	}

	return nil
}

// Return all available articles
func GetArticles() []Article {

	result := make([]Article, 0)

	query := "SELECT * FROM `articles` WHERE `date` > ? ORDER BY `date` DESC;"

	raws, err := db.Query(query, time.Now().AddDate(0,0,-7).Format("2006-01-02"))
	if err != nil {

		log.Println(err)
		return nil
	}
	defer raws.Close()

	for raws.Next() {

		var article Article

		err := raws.Scan(
			&article.Author,
			&article.Header,
			&article.Text,
			&article.Importance,
			&article.Date,
			)

		if err != nil {

			log.Println(err)
			return nil
		}

		result = append(result, article)
	}

	return result
}

func WriteArticleToDB(article Article) error {

	query := "INSERT INTO `articles` VALUES (?, ?, ?, ?, ?);"

	_, err := db.Exec(
		query,
		article.Author,
		article.Header,
		article.Text,
		strconv.Itoa(article.Importance),
		article.Date)

	if err != nil {

		return err
	}

	return nil
}

// Return all available uploaded files
func GetFiles() []File {

	result := make([]File, 0)

	query := "SELECT * FROM `files` WHERE `date` > ? ORDER BY `date` DESC;"

	raws, err := db.Query(query, time.Now().AddDate(0,0,-7).Format("2006-01-02"))
	if err != nil {

		log.Println(err)
		return nil
	}
	defer raws.Close()

	for raws.Next() {

		var file File

		err := raws.Scan(
			&file.Uploader,
			&file.Description,
			&file.Url,
			&file.Date,
		)

		if err != nil {

			log.Println(err)
			return nil
		}

		result = append(result, file)
	}

	return result
}

func WriteFileToDB(file File) error {

	query := "INSERT INTO `files` VALUES (?, ?, ?, ?);"

	_, err := db.Exec(query, file.Uploader, file.Description, file.Url, file.Date)
	if err != nil {

		return err
	}

	return nil
}

func GetRecentMessages(count int) []Message {

	result := make([]Message, 0)

	query := "SELECT * FROM (SELECT * FROM `messages` ORDER BY `date` DESC LIMIT ?) AS M1 ORDER BY `date` ASC;"

	raws, err := db.Query(query, strconv.Itoa(count))
	if err != nil {

		log.Println(err)
		return nil
	}
	defer raws.Close()

	for i := 0; raws.Next(); i++ {

		var message Message

		err := raws.Scan(
			&message.Username,
			&message.Message,
			&message.Date,
		)

		if err != nil {

			log.Println(err)
			return nil
		}

		result = append(result, message)
	}

	return result
}

func WriteMessageToDB(message Message) error{

	query := "INSERT INTO `messages` VALUES (?, ?, ?);"

	_, err := db.Exec(query, message.Username, message.Message, message.Date)
	if err != nil {

		return err
	}

	return nil
}

func isInjected(str string) bool {

	chars := "=!|&()*^%+\\/?[]{}'\" "

	if strings.ContainsAny(str, chars) {

		return true
	}

	return false
}


