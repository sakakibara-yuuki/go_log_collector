//
// main.go
// Copyright (C) 2024 sakakibara <sakakibara@organon>
//
// Distributed under terms of the MIT license.
//


package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// User represents the user structure in the log
type User struct {
	Age  int    `json:"age"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// LogEntry represents the log structure in the log file
type LogEntry struct {
	User User   `json:"user"`
	Dist string `json:"dist"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
	Src   string `json:"src"`
	Time  string `json:"time"`
}

// PostgreSQL databaseに接続
func connectDB() (*sql.DB, error) {
	connStr := "user=root password=password dbname=logdb sslmode=disable host=localhost port=5432"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// users tableの作成
func createUsersTable(db *sql.DB) error {
  // UUIDでもいいかも
  _, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, age integer, name varchar(500), role char(15))")
  return err
}

// users tableにデータを挿入
func insertUser(tx *sql.Tx, user User) error {
	_, err := tx.Exec("INSERT INTO users (age, name, role) VALUES ($1, $2, $3)", user.Age, user.Name, user.Role)
	return err
}

func main() {
	// 引数の数をチェック
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run main.go <logfile>")
	}

	// ログファイルを開く
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// データベースに接続
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

    // ユーザーテーブルを作成
    err = createUsersTable(db)
    if err != nil {
        log.Fatal(err)
    }

	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
    // ファイルを1行ずつ読み込む
	for scanner.Scan() {
		// JSONをパース, 構造体に格納
		var entry LogEntry
		err := json.Unmarshal([]byte(scanner.Text()), &entry)
		if err != nil {
			log.Println("Error parsing log entry:", err)
			tx.Rollback() // ロールバック
			return
		}

		// ユーザーデータをusers tableに挿入
		err = insertUser(tx, entry.User)
		if err != nil {
			log.Println("Error inserting user:", err)
			tx.Rollback() // ロールバック
			return
		}
	}

	// 上記のコードが通ったならコミット
	err = tx.Commit()
	if err != nil {
		log.Fatal("Transaction commit failed:", err)
	}

	fmt.Println("All entries inserted successfully")
}

