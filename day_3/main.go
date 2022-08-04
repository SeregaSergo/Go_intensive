package main

import (
	"elastic_study/internals/api"
	"elastic_study/internals/db"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	dataFile   string
	schemaFile string
)

func init() {
	flag.StringVar(&dataFile, "data", "", "File with data to upload in ES")
	flag.StringVar(&schemaFile, "schema", "schema.json", "File with schema of document type in ES")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
}

func main() {
	es := db.NewESClient("places", "place")
	s := api.NewService(10, 3, time.Minute, "signingKey", es)

	// Upload data to DB if there is a flag with datafile
	if len(dataFile) != 0 {
		if err := es.UploadData(dataFile, schemaFile); err != nil {
			log.Fatalf("Error: %s", err)
		}
	}
	log.Fatal(http.ListenAndServe(":8888", s))
}
