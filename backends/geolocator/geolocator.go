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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type geoInfo struct {
	lat string
	lon string
}

// IPInfo is the {} containing info for single IP address
type IPInfo struct {
	IP  string `json:"ip"`
	Lat string `json:"lat"`
	Lon string `json:"lon"`
	Num string `json:"num"`
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

func fetchAllIP(db *sql.DB, tableName string) []string {
	// Returns an array of all IP addresses in traffic table
	row, err := db.Query(`SELECT ip from ` + tableName)
	checkerr(err)
	defer row.Close()

	ipArr := make([]string, 0)
	var ip string

	for row.Next() {
		err := row.Scan(&ip)
		checkerr(err)

		// Strip '"' from start and end of the field data
		ip = strings.Trim(ip, `"`)

		ipArr = append(ipArr, ip)
	}

	return ipArr
}

func toJSON(respStr []byte) map[string]interface{} {
	var jsonStr map[string]interface{}
	err := json.Unmarshal(respStr, &jsonStr)
	checkerr(err)

	return jsonStr
}

func getGeoInfo(ip string) geoInfo {
	// Get Lat and Lon for given IP address
	// If cannot fetch geoinfo, return both lat and lon as "0"

	lat := "000"
	lon := "000"

	resp, err := http.Get("http://ip-api.com/json/" + ip + "?fields=status,lat,lon")
	if err != nil {
		fmt.Println(err)
		return geoInfo{lat, lon}
	}
	defer resp.Body.Close()
	respStr, err := ioutil.ReadAll(resp.Body)
	checkerr(err)

	/* Sample response: {
		"status": "success",
		"lat": 45.511,
		"lon": -73.5561
	}
	*/
	respJSON := toJSON(respStr)

	// Break down into individual fields
	if respJSON["status"].(string) == "success" {
		lat = float64ToStr(respJSON["lat"].(float64))
		lon = float64ToStr(respJSON["lon"].(float64))
	}

	return geoInfo{lat, lon}
}

func float64ToStr(inFl64 float64) string {
	// The API's accuracy is 4 decimals (~11m)
	return strconv.FormatFloat(inFl64, 'f', 4, 64)
}

func checkIPExist(db *sql.DB, ip string) bool {
	// check if IP exists
	row, err := db.Query("SELECT ip FROM geo WHERE ip = ?", ip)
	checkerr(err)
	defer row.Close()

	for row.Next() {
		return true
	}

	return false
}

func addGeoRecord(db *sql.DB, ip, lat, lon string) {
	// Insert IP with Lat and Lon into the geo table
	insert, err := db.Query("INSERT INTO geo (ip, lat, lon) VALUES ( ?, ?, ? )", ip, lat, lon)
	checkerr(err)
	defer insert.Close()
}

func createDB(db *sql.DB) {
	// Check if the table exists, if not then create a new table
	create, err := db.Query("CREATE TABLE IF NOT EXISTS geo (ip VARCHAR(25) NOT NULL, lat VARCHAR(25) NOT NULL, lon VARCHAR(25) NOT NULL)")
	checkerr(err)
	defer create.Close()
}

func createGeoJSON(db *sql.DB) {
	// Fetch all IP and their geo info
	row, err := db.Query(
		`select geo.ip, geo.lat, geo.lon, traffic.num_attempts from geo inner join traffic on geo.ip=traffic.ip`)
	checkerr(err)
	defer row.Close()

	ipInfoArr := make([]IPInfo, 0)

	var ip, num, lat, lon string

	for row.Next() {
		err := row.Scan(&ip, &lat, &lon, &num)
		checkerr(err)

		ipInfoArr = append(
			ipInfoArr,
			IPInfo{IP: ip, Lat: lat, Lon: lon, Num: num})
	}

	// Convert into JSON
	res, err := json.Marshal(ipInfoArr)
	checkerr(err)

	// Write to file
	f, err := os.Create("geo.json")
	checkerr(err)
	defer f.Close()

	f.WriteString(string(res))
	f.Sync()
}

func main() {
	// Reuse db instance
	db := openDb()
	defer db.Close()
	createDB(db)

	for true {
		for _, ip := range fetchAllIP(db, "traffic") {
			// Check if already exists in geo table
			fmt.Println("Check " + ip)

			if checkIPExist(db, ip) {
				fmt.Println(ip + " exists, skip")
				continue
			}

			fmt.Println(ip + " not exist")
			geo := getGeoInfo(ip)

			fmt.Println("Fetched ip: " + ip + " lat: " + geo.lat + " and lon: " + geo.lon)

			// Check for fetch request lat lon variable
			if geo.lat == "000" && geo.lon == "000" {
				fmt.Println("Detected error in previous fetch response, abort current ip")
				continue
			}

			// Store into geo table
			addGeoRecord(db, ip, geo.lat, geo.lon)
			fmt.Println("Added!")

			// IP-API has a rate limit of 150 requests per minute, sleep for 0.5 second
			time.Sleep(500 * time.Millisecond)
		}

		// Generate JSON file
		fmt.Println("Outputting JSON file")
		createGeoJSON(db)

		fmt.Print("Sleep for 12 hours")
		time.Sleep(12 * time.Hour)
	}
}
