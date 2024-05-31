package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SolaTyolo/graphqltopostman"
	"github.com/SolaTyolo/httpclient"
	"github.com/atombender/go-jsonschema/pkg/generator"
)

const (
	DEFAULT_OUT_FILE = "internal/postman/model.go"
	DEFAULT_PACKAGE  = "postman"
	TEMP_FILE        = "./collection.json"
)

// TODO: use quick type to conver to go struct

// Path: cmd/update_collection_schema/main.go
func main() {
	// outfile 注入
	var out string
	var p string

	flag.StringVar(&out, "out", DEFAULT_OUT_FILE, "output file name")
	flag.StringVar(&p, "package", DEFAULT_PACKAGE, "output go package")

	flag.Parse()

	rawBytes, err := postmanCollectionSchema()
	if err != nil {
		panic(err)
	}

	defer os.Remove(TEMP_FILE)
	if err := generate(rawBytes, out, p); err != nil {
		panic(err)
	}
}

func generate(data []byte, output, p string) error {
	if err := os.WriteFile(TEMP_FILE, data, 0644); err != nil {
		return err
	}

	cfg := generator.Config{
		Warner: func(message string) {
			fmt.Fprintf(os.Stderr, "Warning: %s", message)
		},
		ExtraImports:        false,
		Capitalizations:     nil,
		DefaultOutputName:   output,
		DefaultPackageName:  p,
		SchemaMappings:      []generator.SchemaMapping{},
		ResolveExtensions:   nil,
		YAMLExtensions:      nil,
		StructNameFromTitle: true,
		Tags:                []string{"json"},
		OnlyModels:          true,
	}

	generator, err := generator.New(cfg)
	if err != nil {
		return err
	}

	if err := generator.DoFile(TEMP_FILE); err != nil {
		return err
	}

	for fileName, source := range generator.Sources() {
		if err := os.MkdirAll(filepath.Dir(fileName), 0o755); err != nil {
			return err
		}
		w, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
		if _, err = w.Write(source); err != nil {
			return err
		}
		_ = w.Close()
	}
	return nil
}

// get postman collection jsonschema spec
func postmanCollectionSchema() ([]byte, error) {
	url := graphqltopostman.POSTMAN_SCHEMA_URL
	response, err := httpclient.Default().Get(context.Background(), url, nil)
	if err != nil {
		return nil, err
	}
	return response.RawBody, nil
}
