/*
	PLEASE RUN THIS SCRIPT UNDER SERVER DIRECTORY !
	AS IT NEEDS TO USE RELATIVE PATH TO READ GEO JSON FILE !
*/

package main

import (
	"database/sql"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func openDb() *sql.DB {
	db, err := sql.Open("mysql", "")
	checkerr(err)
	return db
}

// Tunnel struct
type Tunnel struct {
	HTTPID string `json:"http_id"`
	Epoch  string `json:"epoch"`
	SrcIP  string `json:"src_ip"`
	DstIP  string `json:"dst_ip"`
	Data   string `json:"data"`
}

// Command struct
type Command struct {
	InputID string `json:"input_id"`
	Epoch   string `json:"epoch"`
	SrcIP   string `json:"src_ip"`
	Cmd     string `json:"cmd"`
}

// Cred struct
type Cred struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NumAttempts string `json:"num_attempts"`
}

func getTunnelRecords() []Tunnel {
	// Return latest X tunnel records based on decreasing ID as a slice of structs
	db := openDb()
	defer db.Close()

	row, err := db.Query("SELECT http_id, epoch, src_ip, dst_ip, data FROM http order by http_id desc limit 24")
	checkerr(err)
	defer row.Close()

	var httpID, epoch, srcIP, dstIP, data string
	var t Tunnel

	ts := make([]Tunnel, 0)

	for row.Next() {
		err := row.Scan(&httpID, &epoch, &srcIP, &dstIP, &data)
		checkerr(err)

		// Strip '"' from start and end of the field data
		srcIP = strings.Trim(srcIP, `"`)
		dstIP = strings.Trim(dstIP, `"`)
		data = strings.Replace(strings.Trim(data, `"'`), `\\`, `\`, -1)

		// Craft slice of structs
		t = Tunnel{httpID, epoch, srcIP, dstIP, data}
		ts = append(ts, t)
	}

	return ts
}

func getCommandRecords() []Command {
	// Return latest X command records based on decreasing ID as a slice of structs
	db := openDb()
	defer db.Close()

	row, err := db.Query("SELECT input_id, epoch, src_ip, cmd FROM input order by input_id desc limit 24")
	checkerr(err)
	defer row.Close()

	var inputID, epoch, srcIP, cmd string
	var c Command

	cs := make([]Command, 0)

	for row.Next() {
		err := row.Scan(&inputID, &epoch, &srcIP, &cmd)
		checkerr(err)

		cmd = strings.Replace(cmd, `\\`, `\`, -1)

		// Craft slice of structs
		c = Command{inputID, epoch, srcIP, cmd}
		cs = append(cs, c)
	}

	return cs
}

func listLoginCreds() []Cred {
	// Return top 20 login creds based on decreasing number of attempts as a slice of structs
	db := openDb()
	defer db.Close()

	row, err := db.Query(`SELECT username, password, num_attempts FROM dictionary WHERE username!="\"root\"" AND password!="\"root\"" order by num_attempts desc limit 30`)
	checkerr(err)
	defer row.Close()

	var username, password, numAttempts string
	var c Cred

	cs := make([]Cred, 0)

	for row.Next() {
		err := row.Scan(&username, &password, &numAttempts)
		checkerr(err)

		// Strip '"' from start and end of the field data
		username = strings.Trim(username, `"`)
		password = strings.Trim(password, `"' `)

		// Craft slice of structs
		c = Cred{username, password, numAttempts}
		cs = append(cs, c)
	}

	return cs
}

func serverHeader(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
}

func fetchGeoJSON() string {
	// Read file containing Geo json and return the content
	path := "../geolocator/geo.json"

	b, err := ioutil.ReadFile(path)
	checkerr(err)

	return string(b)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Change Header to allow CORS
	router.Use(serverHeader)

	router.GET("/tunnel", func(c *gin.Context) {
		c.JSON(200, getTunnelRecords())
	})

	router.GET("/cmd", func(c *gin.Context) {
		c.JSON(200, getCommandRecords())
	})

	router.GET("/login", func(c *gin.Context) {
		c.JSON(200, listLoginCreds())
	})

	router.GET("/geo", func(c *gin.Context) {
		c.String(200, fetchGeoJSON())
	})

	router.Run(":8001")
}
