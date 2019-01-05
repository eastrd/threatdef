/*
	This script fetches all IP addresses from MySQL database periodically,
	1. Checks each address:
		- If already exists in geo table, then skip to next address
		- Otherwise, fetch geolocation info and insert the info with IP to geo table.
	2. Output a JSON file containing all existing IP with their lat & lon, and number of records.
*/

package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type geoInfo struct {
	lat string
	lon string
}

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func openDb() *sql.DB {
	// Returns an authenticated db instance
	db, err := sql.Open("mysql", "")
	checkerr(err)
	return db
}

func fetchAllIP(db *sql.DB) []string {
	// Returns an array of all IP addresses in traffic table
	row, err := db.Query(`SELECT ip from traffic`)
	checkerr(err)
	defer row.Close()

	ipArr := make([]string, 0)
	var ip string

	for row.Next() {
		err := row.Scan(&ip)
		checkerr(err)

		// Strip '"' from start and end of the field data
		ip = strings.Trim(ip, `"`)
		fmt.Println(ip)

		ipArr = append(ipArr, ip)
	}

	return ipArr
}

func getGeoInfo(ip string) geoInfo {
	// Get Lat and Lon for given IP address
	return geoInfo{"0", "0"}
}

func main() {
	for true {
		// Reuse db instance
		db := openDb()
		defer db.Close()

		for _, ip := range fetchAllIP(db) {
			// Check if already exists in geo table
			fmt.Println("Check" + ip)

			// check if IP exists
			row, err := db.Query("SELECT ip, num_attempts FROM traffic WHERE ip = ?", ip)
			checkerr(err)
			defer row.Close()

			ipExist := false

			for row.Next() {
				ipExist = true
			}

			if ipExist {
				fmt.Println(ip + "exists, skip")
			} else {
				fmt.Println(ip + " not exist")
				getGeoInfo(ip)
			}
		}
		fmt.Print("Sleep for 12 hours")
		time.Sleep(12 * time.Hour)
	}
}
