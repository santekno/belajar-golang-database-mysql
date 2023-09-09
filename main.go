package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// initiate context
	ctx := context.Background()

	// initiate get connection to db
	db := GetConnection()
	defer db.Close()

	// get function insert into database
	// InsertIntoDB(ctx, db)

	// get data customer
	// GetDataCustomer(ctx, db)

	// username := "santekno'; DROP TABLE user; #"
	// password := "santekno"
	// Login(ctx, db, username, password)
	// Register(ctx, db, username, password)

	// email := "santekno@gmail.com"
	// comment := "hello world"
	// lastID := InsertComment(ctx, db, email, comment)
	// fmt.Println("Last Insert ID: ", lastID)

	PrepareStatement(ctx, db)

	TransactionDatabase(ctx, db)

}

func GetConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:belajargolang@tcp(localhost:3306)/belajar-golang?parseTime=true")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)
	return db
}

func InsertIntoDB(ctx context.Context, db *sql.DB) {
	_, err := db.ExecContext(ctx, "INSERT INTO customer(id,name) VALUES('santekno','Santekno');")
	if err != nil {
		panic(err)
	}
	fmt.Println("success insert data to database")
}

func GetDataCustomer(ctx context.Context, db *sql.DB) {
	rows, err := db.QueryContext(ctx, "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float32
		var birthDate sql.NullTime
		var createdAt time.Time
		var married bool
		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}
		fmt.Println("Id : ", id)
		fmt.Println("Name: ", name)
		if email.Valid {
			fmt.Println("Email: ", email)
		}
		fmt.Println("Balance: ", balance)
		fmt.Println("Rating: ", rating)
		if birthDate.Valid {
			fmt.Println("Birth Date: ", birthDate)
		}
		fmt.Println("Married: ", married)
		fmt.Println("Created At: ", createdAt)
	}
}

func Login(ctx context.Context, db *sql.DB, username, password string) bool {
	sqlQuery := "SELECT username FROM user WHERE username=? AND password=? LIMIT 1"
	rows, err := db.QueryContext(ctx, sqlQuery, username, password)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var username string
		rows.Scan(&username)
		fmt.Println("Success Login ", username)
		return true
	} else {
		fmt.Println("Failed Login")
	}
	return false
}

func Register(ctx context.Context, db *sql.DB, username, password string) bool {
	_, err := db.ExecContext(ctx, "INSERT INTO user(username, password) VALUE(?,?)", username, password)
	if err != nil {
		panic(err)
	}

	fmt.Println("success insert new user")
	return true
}

func InsertComment(ctx context.Context, db *sql.DB, email, comment string) int64 {
	sqlQuery := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	result, err := db.ExecContext(ctx, sqlQuery, email, comment)
	if err != nil {
		panic(err)
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	return insertID
}

func PrepareStatement(ctx context.Context, db *sql.DB) {
	query := "INSERT INTO comments(email,comment) VALUE(?, ?)"
	statement, err := db.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	for i := 0; i < 10; i++ {
		email := "Santekno " + strconv.Itoa(i) + "@gmail.com"
		comment := "Komentar ke " + strconv.Itoa(i)
		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("comment id ", id)
	}
}

func TransactionDatabase(ctx context.Context, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	query := "INSERT INTO comments(email,comment) VALUE(?, ?)"

	// do transaction
	for i := 0; i < 10; i++ {
		email := "Santekno " + strconv.Itoa(i) + "@gmail.com"
		comment := "Komentar ke " + strconv.Itoa(i)
		result, err := tx.ExecContext(ctx, query, email, comment)
		if err != nil {
			panic(err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("comment id ", id)
	}

	err = tx.Rollback()
	if err != nil {
		panic(err)
	}
}
