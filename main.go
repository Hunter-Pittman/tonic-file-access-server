package main

import (
	"flag"
	"tonic-file-access-server/server"
)

func main() {
	apiToken := flag.String("api-token", "secret", "Provide a secret API token for use with the tonic API")

	flag.Parse()

	if *apiToken == "secret" {
		panic("Please provide a secret API token (use --help for more info)")
	}

	server.Init(*apiToken)

	//server.Init("bruh")
}
