package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func toJSON(c *gin.Context) map[string]interface{} {
	// Bytes to String
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)

	// fmt.Println(buf.String())
	// Convert to JSON object
	var jsonStr map[string]interface{}
	err := json.Unmarshal([]byte(buf.String()), &jsonStr)
	checkerr(err)
	// fmt.Println(jsonStr["eventid"])
	return jsonStr
}

func addLoginAttempt(username, password string) {
	// Log statistics for brute force combinations

	// Connect to MySQL
	// fmt.Println(username + ":" + password)
	db, err := sql.Open("mysql", "threatdef:194122602@tcp(108.61.169.45:3306)/threatdef")
	checkerr(err)
	defer db.Close()

	// check if login combination exists
	row, err := db.Query("SELECT username, password, num_attempts FROM dictionary WHERE username = ? AND password = ?", username, password)
	checkerr(err)
	defer row.Close()

	count := 0
	for row.Next() {
		count++
		fmt.Println("Found a record...")
		var numAttempts string
		err := row.Scan(&username, &password, &numAttempts)
		checkerr(err)

		// Increment number of attempts
		stmt, err := db.Prepare("UPDATE dictionary set num_attempts=? where username=? AND password=?")
		checkerr(err)
		defer stmt.Close()

		fmt.Println("Num attempts = " + numAttempts)
		numAttemptsINT, err := strconv.Atoi(numAttempts)
		checkerr(err)
		newNumAttemptsINT := numAttemptsINT + 1
		newNumAttempts := strconv.Itoa(newNumAttemptsINT)

		stmt.Exec(newNumAttempts, username, password)
	}

	if count == 0 {
		fmt.Println("Didn't find any records for combination " + username + ":" + password)
		// Insert into dictionary table as a new combination
		insert, err := db.Query("INSERT INTO dictionary VALUES ( ?, ?, 1 )", username, password)
		checkerr(err)
		defer insert.Close()
	}
}

func addIPstats(srcIP string) {
	// Log IP statistics information
	db, err := sql.Open("mysql", "threatdef:194122602@tcp(108.61.169.45:3306)/threatdef")
	checkerr(err)
	defer db.Close()

	// check if IP exists
	row, err := db.Query("SELECT ip, num_attempts FROM traffic WHERE ip = ?", srcIP)
	checkerr(err)
	defer row.Close()

	count := 0
	for row.Next() {
		count++
		fmt.Println("Found a record...")
		var numAttempts string
		err := row.Scan(&srcIP, &numAttempts)
		checkerr(err)

		// Increment number of attempts
		stmt, err := db.Prepare("UPDATE traffic set num_attempts=? where ip=?")
		checkerr(err)
		defer stmt.Close()

		fmt.Println("Num attempts for " + srcIP + " : " + numAttempts)
		numAttemptsINT, err := strconv.Atoi(numAttempts)
		checkerr(err)
		newNumAttemptsINT := numAttemptsINT + 1
		newNumAttempts := strconv.Itoa(newNumAttemptsINT)

		stmt.Exec(newNumAttempts, srcIP)
	}

	if count == 0 {
		fmt.Println("Didn't find any records for IP " + srcIP)
		// Insert into dictionary as a new IP
		insert, err := db.Query("INSERT INTO traffic VALUES ( ?, 1 )", srcIP)
		checkerr(err)
		defer insert.Close()
	}
}

func addTunnelData(epoch, srcIP, dstIP, data string) {
	// Log tunneling HTTP data
	db, err := sql.Open("mysql", "threatdef:194122602@tcp(108.61.169.45:3306)/threatdef")
	checkerr(err)
	defer db.Close()

	// Create table if not exist
	create, err := db.Query("CREATE TABLE IF NOT EXISTS http (http_id INT AUTO_INCREMENT, epoch INT NOT NULL, src_ip VARCHAR(25) NOT NULL, dst_ip VARCHAR(25) NOT NULL, data TEXT, PRIMARY KEY (http_id))")
	checkerr(err)
	defer create.Close()

	// Insert into http table
	insert, err := db.Query("INSERT INTO http (epoch, src_ip, dst_ip, data) VALUES ( ?, ?, ?, ? )", epoch, srcIP, dstIP, data)
	checkerr(err)
	defer insert.Close()
}

func processJSON(jsonStr map[string]interface{}) {
	// Detect event type based on EventID and collect relevant information and send them to mysql
	srcIP, _ := json.Marshal(jsonStr["src_ip"])
	epoch, _ := json.Marshal(jsonStr["epoch"])
	addIPstats(string(srcIP))
	switch jsonStr["eventid"] {
	case "cowrie.login.success", "cowrie.login.failed":
		username, _ := json.Marshal(jsonStr["username"])
		password, _ := json.Marshal(jsonStr["password"])
		addLoginAttempt(string(username), string(password))

	case "cowrie.direct-tcpip.data":
		// Tunnel Request: epoch int, src_ip varchar, dst_ip varchar, data varchar
		dstIP, _ := json.Marshal(jsonStr["dst_ip"])
		data, _ := json.Marshal(jsonStr["data"])
		// Check that the Data is in plain text instead of HTTPS connection (i.e. No "\\x{number}")
		if !strings.Contains(string(data), "\\\\x") {
			addTunnelData(string(epoch), string(srcIP), string(dstIP), string(data))
		}
	}
}

func main() {
	r := gin.Default()

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"plus": "midoriya",
	}))

	authorized.POST("/signal", func(c *gin.Context) {
		processJSON(toJSON(c))
		c.JSON(200, gin.H{
			"signal": "true",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
