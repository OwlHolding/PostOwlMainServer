package main

import (
	"database/sql"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var DatabaseClient *sql.DB

func InitDatabase(config ServerConfig) {
	MySqlClient, err := sql.Open("mysql", config.SqlUser+":"+config.SqlPass+"@/postowl")
	if err != nil {
		log.Fatal(err)
	}

	MySqlClient.SetMaxOpenConns(config.MaxSqlConns)
	MySqlClient.SetMaxIdleConns(config.MaxSqlIdleConns)

	err = MySqlClient.Ping()
	if err != nil {
		log.Fatal(err)
	}

	DatabaseClient = MySqlClient
	log.Println("DataBase: inited")
}

func DataBaseAddUser(userID int64) bool {
	_, err := DatabaseClient.Exec("INSERT INTO users VALUES (?, '')", userID)
	sqlerr, _ := err.(*mysql.MySQLError)
	if err != nil {
		if sqlerr.Number == 1062 { // primary key already exists
			return false
		} else {
			panic(err)
		}
	}

	log.Printf("DataBase: create user %d \n", userID)
	return true
}

func DataBaseDelUser(userID int64) bool {
	_, err := DatabaseClient.Exec("DELETE FROM users WHERE id=?", userID)

	if err != nil {
		return false
	}

	log.Printf("DataBase: delete user %d \n", userID)
	return true
}

func DataBaseAddChannel(userID int64, channel string) bool {
	result, err := DatabaseClient.Query("SELECT * FROM users WHERE id=?", userID)
	if err != nil {
		panic(err)
	}

	var id int64
	var row_channels string
	result.Next()
	result.Scan(&id, &row_channels)

	if id == 0 {
		return false
	}

	if strings.Contains(row_channels, channel) {
		return false
	}

	channels := strings.Split(row_channels, "&")
	row_channels = strings.Join(append(channels, channel), "&")

	_, err = DatabaseClient.Exec("UPDATE users SET channels=? WHERE id=?", row_channels, userID)
	if err != nil {
		panic(err)
	}

	log.Printf("DataBase: Add channel %s to user %d \n", channel, userID)
	return true
}

func DataBaseDelChannel(userID int64, channel string) bool {
	result, err := DatabaseClient.Query("SELECT * FROM users WHERE id=?", userID)
	if err != nil {
		panic(err)
	}

	var id int64
	var row_channels string
	result.Next()
	result.Scan(&id, &row_channels)

	if id == 0 {
		return false
	}

	if !strings.Contains(row_channels, channel) {
		return false
	}

	channels := strings.Split(row_channels, "&")
	var new_channels []string
	for _, one_channnel := range channels {
		if one_channnel != channel {
			new_channels = append(new_channels, one_channnel)
		}
	}

	row_channels = strings.Join(new_channels, "&")

	_, err = DatabaseClient.Exec("UPDATE users SET channels=? WHERE id=?", row_channels, userID)
	if err != nil {
		panic(err)
	}

	log.Printf("DataBase: Del channel %s to user %d \n", channel, userID)
	return true
}

func DataBaseInfo(userID int64) []string {
	result, err := DatabaseClient.Query("SELECT * FROM users WHERE id=?", userID)
	if err != nil {
		panic(err)
	}

	var id int64
	var channels string
	result.Next()
	result.Scan(&id, &channels)

	return strings.Split(channels, "&")
}

func DataBaseGetUsers(channel string) []int64 {
	query := "SELECT id FROM users WHERE channels LIKE '%" + channel + "%'"
	result, err := DatabaseClient.Query(query)
	if err != nil {
		panic(err)
	}

	var users []int64
	var user int64
	for result.Next() {
		result.Scan(&user)
		users = append(users, user)
	}

	return users
}

func DataBaseIsUserExist(userID int64) bool {
	result, err := DatabaseClient.Query("SELECT id FROM users WHERE id=?", userID)
	if err != nil {
		panic(err)
	}

	var id int64
	result.Next()
	result.Scan(&id)

	if id != 0 {
		return true
	}

	return false
}
