package main

import (
	"os"
	"testing"

	"github.com/SolaTyolo/graphqltopostman"
)

func TestXxx(t *testing.T) {
	url := "https://swapi-graphql.netlify.app/.netlify/functions/index"
	if err := graphqltopostman.Convert("scm-service", url, os.Stdout); err != nil {
		panic(err)
	}
}
