package main

import (
	"database/sql"
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
	db, err := sql.Open("mysql", "threatdef:194122602@tcp(108.61.169.45:3306)/threatdef")
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

func getTunnelRecords() []Tunnel {
	// Return latest X tunnel records based on decreasing ID in JSON format
	db := openDb()
	defer db.Close()

	row, err := db.Query("SELECT http_id, epoch, src_ip, dst_ip, data FROM http order by http_id desc limit 10")
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
	// Return latest X command records based on decreasing ID in JSON format
	db := openDb()
	defer db.Close()

	row, err := db.Query("SELECT input_id, epoch, src_ip, cmd FROM input order by input_id desc limit 10")
	checkerr(err)
	defer row.Close()

	var inputID, epoch, srcIP, cmd string
	var c Command

	cs := make([]Command, 0)

	for row.Next() {
		err := row.Scan(&inputID, &epoch, &srcIP, &cmd)
		checkerr(err)

		// Strip '"' from start and end of the field data
		srcIP = strings.Trim(srcIP, `"`)
		cmd = strings.Replace(strings.Trim(cmd, `"' `), `\\`, `\`, -1)

		// Craft slice of structs
		c = Command{inputID, epoch, srcIP, cmd}
		cs = append(cs, c)
	}

	return cs
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/tunnel", func(c *gin.Context) {
		c.JSON(200, getTunnelRecords())
	})

	router.GET("/cmd", func(c *gin.Context) {
		c.JSON(200, getCommandRecords())
	})

	router.Run(":8001")
}