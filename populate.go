package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/machinebox/graphql"
)

// AwaitClient test if the specified client can connect with its
// source. On this script, the client will need to contain a Product
// and a User table, so it queries for them. If the connection fails,
// it tries again after 4 seconds, within a limit of tries.
func AwaitClient(clt *graphql.Client) {
	// env controls the limit of connect tries.
	awaitLimit, _ := strconv.Atoi(getenv("GRAPHQL_CHECK_TABLE_LIMIT", "50"))

	for i := 0; i < awaitLimit; i++ {
		if err := clt.Run(context.Background(), graphql.NewRequest(EnvContent("GQL_CHECK_TABLES")), nil); err != nil {
			log.Print("Client not connecting, trying in 4 seconds")
			time.Sleep(4 * time.Second)
		} else {
			break
		}
	}

}

// MakeRequest receives a request content, uses client to make the
// request. It adds no-cache options to request. If request fails,
// it will crash the script.
func MakeRequest(clt *graphql.Client, content string) {
	// make a request and set headers
	req := graphql.NewRequest(content)
	req.Header.Set("Cache-Control", "no-cache")

	// run it and check response
	var res struct{}
	if err := clt.Run(context.Background(), req, &res); err != nil {
		log.Fatal(err)
	} else {
		log.Print(res)
	}
}

// EnvContent read data from a url stored in a env var.
func EnvContent(env string) string {
	// get url from env
	url := os.Getenv(env)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// read data from url
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}

func main() {
	// Create a client
	client := graphql.NewClient(os.Getenv("GRAPHQL_URL"))

	// Test if Client alive...
	AwaitClient(client)

	// Populate, use upsert mutation to add items only if missing,
	// to avoid multiple inserts, when restarting the docker setup.
	MakeRequest(client, EnvContent("GQL_POPULATE_FILE"))
}

// getenv returns an env variable value, or the fallback value
// specified.
func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}
