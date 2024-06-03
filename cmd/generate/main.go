package main

import (
	"flag"
	"os"

	"github.com/SolaTyolo/graphqltopostman"
)

func main() {
	// outfile 注入
	var url string

	flag.StringVar(&url, "url", "", "graphql endpoint url")

	flag.Parse()

	if err := graphqltopostman.Convert("scm-service", url, os.Stdout); err != nil {
		panic(err)
	}
}
