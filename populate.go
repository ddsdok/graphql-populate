package main

import (
	"os"
	"strconv"
	"time"

	"github.com/ddsgok/gql"
	"github.com/ddsgok/log"
)

// AwaitClient test if the specified client can connect with its
// source. Make the check via query file at GQL_CHECK_TABLES url
// If the connection fails, it tries again after 4 seconds, within a
// limit of tries.
func AwaitClient() (err error) {
	// env controls the limit of connect tries.
	envLimit := os.Getenv("GRAPHQL_CHECK_TABLE_LIMIT")

	var awaitLimit int
	if len(envLimit) > 0 {
		if awaitLimit, err = strconv.Atoi(envLimit); err != nil {
			return
		}
	} else {
		awaitLimit = 50
	}

	for i := 0; i < awaitLimit; i++ {
		if _, err = gql.ReadRequest(os.Getenv("GQL_CHECK_TABLES")).SetHeader("Cache-Control", "no-cache").Run(); err != nil {
			log.Print("message: client not connecting, trying in 4 seconds")
			time.Sleep(4 * time.Second)
		} else {
			err = nil
			break
		}
	}

	return
}

// Populate uses query file from GQL_POPULATE_FILE URL and runs it.
// It adds no-cache options to request.
func Populate() (err error) {
	_, err = gql.ReadRequest(os.Getenv("GQL_POPULATE_FILE")).SetHeader("Cache-Control", "no-cache").Run()
	return
}

func main() {
	gql.Connect()
	defer gql.Disconnect()

	// Test if Client alive...
	if err := AwaitClient(); err != nil {
		log.Fatal("fatal: awaiting client failed err=%s\n", err.Error())
	}

	// Populate, use upsert mutation to add items only if missing,
	// to avoid multiple inserts, when restarting the docker setup.
	if err := Populate(); err != nil {
		log.Fatal("fatal: problem when populating err=%s\n", err.Error())
	}
}
