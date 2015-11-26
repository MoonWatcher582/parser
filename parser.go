package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var FILE_NAME = "src/Parser/losangeles.txt"
var DB_NAME = "src/Parser/restaurants.db"

var restaurantRegex = regexp.MustCompile(`^\s*(.+?) - rank: (\d+), stars:(\d+\.?\d*)\s*$`)

// Restaurant spelled Restuarant because data contains a typo
var reviewRegex = regexp.MustCompile(`^\s*(\d) by (.*?)\(Reviews: (\d+|None), Restuarant Reviews: (\d+|None), Helpful Votes: (\d+|None)\) on ([\d-]+)\s*$`)

type parser struct {
	*sql.DB
	users map[string]string
}

func (p parser) databaseInsert(sql, class string, args ...interface{}) (int64, error) {
	res, err := p.Exec(sql, args...)
	if err != nil {
		fmt.Println("Failed to update database with", class, "!", err, "sql was", sql)
		return -1, err
	}
	return res.LastInsertId()
}

func (p parser) userExists(userName string) bool {
	_, ok := p.users[userName]
	return ok
}

func (p parser) insertUser(name, reviews, restaurantReviews, helpfulVotes string) error {
	if b := p.userExists(name); b {
		return nil
	}

	sql := "INSERT INTO reviewer VALUES (NULL, ?, ?, ?, ?)"
	id, err := p.databaseInsert(sql, "reviewer", name, reviews, restaurantReviews, helpfulVotes)
	if err != nil {
		fmt.Println("Failed to update database with reviewer!", err)
		return err
	}

	p.users[name] = string(id)

	return nil
}

func (p parser) insertReview(restaurantId, user, stars, date string) error {
	sql := "INSERT INTO review VALUES (NULL, ?, (SELECT u.id FROM reviewer u WHERE u.user_name = ?), ?, ?)"
	_, err := p.databaseInsert(sql, "review", restaurantId, user, stars, date)
	if err != nil {
		fmt.Println("Failed to update database with review!", err)
		return err
	}
	return nil
}

func (p parser) parseReview(line, restaurantId string) error {
	lineData := reviewRegex.FindStringSubmatch(line)
	stars := lineData[1]
	reviewer := lineData[2]
	numReviews := lineData[3]
	numRestaurantReviews := lineData[4]
	helpfulVotes := lineData[5]
	date := lineData[6]

	if numReviews == "None" {
		numReviews = "0"
	}

	if numRestaurantReviews == "None" {
		numRestaurantReviews = "0"
	}

	if helpfulVotes == "None" {
		helpfulVotes = "0"
	}

	err := p.insertUser(reviewer, numReviews, numRestaurantReviews, helpfulVotes)
	if err != nil {
		fmt.Println("Error writing reviewer to database!", err)
		return err
	}

	err = p.insertReview(restaurantId, reviewer, stars, date)
	if err != nil {
		fmt.Println("Error writing review to database!", err)
		return err
	}
	return nil
}

func (p parser) insertRestaurant(name, rank, stars string) (string, error) {
	sql := "INSERT INTO restaurant VALUES (NULL, ?, ?, ?);"
	id, err := p.databaseInsert(sql, "restaurant", name, rank, stars)
	if err != nil {
		fmt.Println("Failed to update database with restaurant!", err)
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func (p parser) parseRestaurant(line string) (string, error) {
	lineData := restaurantRegex.FindStringSubmatch(line)
	name := lineData[1]
	rank := lineData[2]
	stars := lineData[3]
	id, err := p.insertRestaurant(name, rank, stars)
	if err != nil {
		fmt.Println("Error writing restaurant to database!", err)
		return "", err
	}
	return id, nil
}

func main() {
	connection, err := sql.Open("sqlite3", DB_NAME)
	if err != nil {
		fmt.Println("Error opening database connection", err)
		return
	}
	defer connection.Close()
	parser := parser{connection, make(map[string]string)}
	file, err := os.Open(FILE_NAME)
	if err != nil {
		fmt.Println("Error opening file '", FILE_NAME, "':", err)
	}
	scanner := bufio.NewScanner(file)
	var currentRestaurantId string
	counter := 0
	for scanner.Scan() {
		counter += 1
		if counter%15 == 0 {
			fmt.Print(" ", counter, "/37287 lines parsed\r")
		}
		line := scanner.Text()
		if len(line) == 0 || line[0] == '-' {
			continue
		}

		if _, err = strconv.Atoi(string(line[0])); err == nil {
			//review
			err = parser.parseReview(line, currentRestaurantId)
			if err != nil {
				fmt.Println("Failed to parse a review!", err)
				return
			}
		} else {
			// restaurant
			currentRestaurantId, err = parser.parseRestaurant(line)
			if err != nil {
				fmt.Println("Failed to parse a restaurant!", err)
				return
			}
		}
	}
}
