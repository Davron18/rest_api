package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gopherschool/http-rest-api/internal/app/apiserver"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()

	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		log.Fatal(err)
	}
	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		fmt.Println(err)
		log.Fatal()
	}

}

func listUsers(db *sql.DB, pageNumber int, pageSize int) {
	offset := (pageNumber - 1) * pageSize
	rows, err := db.Query(`
		SELECT  id, email,encrypted_password
		FROM users
		ORDER BY id 
		LIMIT $1 OFFSET $2`, pageSize, offset)

	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	fmt.Println("Users:")
	for rows.Next() {
		var id int
		var email, encrypted_password string
		err := rows.Scan(&id, &email, &encrypted_password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(email, encrypted_password)
	}
}

func getUserByID(db *sql.DB, userID int) {
	var id int
	var email, encryptedPassword string
	err := db.QueryRow(`
		SELECT id, email, encrypted_password
		FROM users
		WHERE id = $1`, userID).Scan(&id, &email, &encryptedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User not found")
		} else {
			log.Fatal(err)
		}
		return

	}
	fmt.Println(email, encryptedPassword)

}

func editUser(db *sql.DB, userID int, newEmail string, newEncryptedPassword string) {
	_, err := db.Exec(`
		UPDATE users
		SET email = $1, encrypted_password = $2
		WHERE id = $3`, newEmail, newEncryptedPassword, userID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User updated")
}

func deleteUser(db *sql.DB, userID int) {
	_, err := db.Exec(`
		DELETE FROM users WHERE id = $1`, userID)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User deleted")
}
