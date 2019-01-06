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
	db := openDb()
	defer db.Close()

	// Create table if not exist
	create, err := db.Query("CREATE TABLE IF NOT EXISTS dictionary (username VARCHAR(255) NOT NULL, password VARCHAR(255) NOT NULL, num_attempts INT NOT NULL)")
	checkerr(err)
	defer create.Close()

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

func openDb() *sql.DB {
	db, err := sql.Open("mysql", "")
	checkerr(err)
	return db
}

func addIPstats(srcIP string) {
	// Log IP statistics information
	db := openDb()
	defer db.Close()

	// Create table if not exist
	create, err := db.Query("CREATE TABLE IF NOT EXISTS traffic (ip VARCHAR(25) NOT NULL, num_attempts INT NOT NULL)")
	checkerr(err)
	defer create.Close()

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
	db := openDb()
	defer db.Close()

	// Create table if not exist
	create, err := db.Query("CREATE TABLE IF NOT EXISTS http (http_id INT AUTO_INCREMENT, epoch VARCHAR(25) NOT NULL, src_ip VARCHAR(25) NOT NULL, dst_ip VARCHAR(25) NOT NULL, data TEXT, PRIMARY KEY (http_id))")
	checkerr(err)
	defer create.Close()

	// Insert into http table
	insert, err := db.Query("INSERT INTO http (epoch, src_ip, dst_ip, data) VALUES ( ?, ?, ?, ? )", epoch, srcIP, dstIP, data)
	checkerr(err)
	defer insert.Close()
}

func addInput(epoch, srcIP, cmd string) {
	// Log input interaction commands
	db := openDb()
	defer db.Close()

	// Create table if not exist
	create, err := db.Query("CREATE TABLE IF NOT EXISTS input (input_id INT AUTO_INCREMENT, epoch VARCHAR(25) NOT NULL, src_ip VARCHAR(25) NOT NULL, cmd TEXT, PRIMARY KEY (input_id))")
	checkerr(err)
	defer create.Close()

	// Insert into input table
	insert, err := db.Query("INSERT INTO input (epoch, src_ip, cmd) VALUES (?, ?, ?)", epoch, srcIP, cmd)
	checkerr(err)
	defer insert.Close()
}

func processJSON(jsonStr map[string]interface{}) {
	// Detect event type based on EventID and collect relevant information and send them to mysql
	epoch := strconv.FormatFloat(jsonStr["epoch"].(float64), 'f', 0, 64)
	srcIP := strings.Trim(jsonStr["src_ip"].(string), `"`)

	addIPstats(srcIP)

	switch jsonStr["eventid"] {

	case "cowrie.login.success", "cowrie.login.failed":
		username := strings.Trim(jsonStr["username"].(string), `"`)
		password := strings.Trim(jsonStr["password"].(string), `"' `)
		addLoginAttempt(username, password)

	case "cowrie.direct-tcpip.data":
		// Tunnel Request: epoch int, src_ip varchar, dst_ip varchar, data varchar
		dstIP := strings.Trim(jsonStr["dst_ip"].(string), `"`)
		data := strings.Trim(jsonStr["data"].(string), `"'`)
		// Check that the Data is in plain text instead of HTTPS connection (i.e. No "\\x{number}")
		if !strings.Contains(data, "\\\\x") && (strings.Contains(data, "GET") || strings.Contains(data, "POST")) {
			addTunnelData(epoch, srcIP, dstIP, data)
		}

	case "cowrie.command.input":
		// Capture Command inputs
		cmd := strings.Replace(strings.Trim(jsonStr["input"].(string), `"' `), `\\`, `\`, -1)
		if cmd != `""` {
			addInput(epoch, srcIP, cmd)
		}
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"plus": "midoriya",
	}))

	authorized.POST("/signal", func(c *gin.Context) {
		processJSON(toJSON(c))
		c.JSON(200, gin.H{
			"signal": "true",
		})
	})

	router.Run(":8080")
}
