package db

import (
	"bytes"
	"elastic_study/internals/api"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dariubs/percent"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/gocarina/gocsv"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func NewESClient(indexName string, docType string) *ElasticAPI {
	// Create the Elasticsearch client
	var err error
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	return &ElasticAPI{
		Client:    client,
		indexName: indexName,
		docType:   docType,
	}
}

type bulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			ID     string `json:"_id"`
			Result string `json:"result"`
			Status int    `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
				Cause  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"caused_by"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}

type ElasticAPI struct {
	Client    *elasticsearch.Client
	indexName string
	docType   string
}

func (es *ElasticAPI) UploadData(dataFile string, schemaFilepath string) error {
	var (
		buf bytes.Buffer
		res *esapi.Response
		raw map[string]interface{}
		blk *bulkResponse

		batch  = 255
		places []*api.Place

		numItems   int
		numErrors  int
		numIndexed int
		numBatches int
		currBatch  int
	)

	// making CSV Reader
	CSVFile, err := os.OpenFile(dataFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer CSVFile.Close()

	// setting Tab as delimiter
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '\t'
		return r
	})

	// Load places from file
	if err := gocsv.UnmarshalFile(CSVFile, &places); err != nil {
		log.Fatalf("Error of reading CSV file: %s", err)
	}
	count := len(places)
	if count%batch == 0 {
		numBatches = count / batch
	} else {
		numBatches = (count / batch) + 1
	}
	fmt.Printf("Num of places %d\nNum of batches %d\n", count, numBatches)

	// recreating Index
	if res, err = es.Client.Indices.Delete([]string{es.indexName}); err != nil {
		log.Fatalf("Cannot delete index: %s", err)
	}
	res, err = es.Client.Indices.Create(es.indexName)
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}

	// adding mapping schema
	schemaFile, err := os.Open(schemaFilepath)
	if err != nil {
		log.Fatalf("Error of oppening the schema file: %s", err)
	}
	res, err = es.Client.Indices.PutMapping(schemaFile,
		es.Client.Indices.PutMapping.WithIndex(es.indexName),
		es.Client.Indices.PutMapping.WithDocumentType(es.docType),
		es.Client.Indices.PutMapping.WithIncludeTypeName(true))
	if err != nil {
		log.Fatalf("Error of applying mapping schema: %s", err)
	}

	// Loop over the collection and making index requests
	for i, p := range places {
		numItems++
		currBatch = i / batch
		if i == count-1 {
			currBatch++
		}

		// Prepare the metadata payload
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, p.ID, "\n"))

		// Prepare the data payload: encode place to JSON
		data, err := json.Marshal(p)
		if err != nil {
			log.Fatalf("Cannot encode place %d: %s", p.ID, err)
		}
		data = append(data, "\n"...)

		// Append payloads to the buffer (ignoring write errors)
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

		//When a threshold is reached, execute the Bulk() request with body from buffer
		if i > 0 && i%batch == 0 || i == count-1 {
			fmt.Printf("[%f]\n", percent.PercentOf(currBatch, numBatches))

			res, err = es.Client.Bulk(bytes.NewReader(buf.Bytes()), es.Client.Bulk.WithIndex(es.indexName))
			if err != nil {
				log.Fatalf("Failure indexing batch %d: %s", currBatch, err)
			}

			// If the whole request failed, print error and mark all documents as failed
			if res.IsError() {
				numErrors += numItems
				if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
					log.Fatalf("Failure to to parse response body: %s", err)
				} else {
					log.Printf("  Error: [%d] %s: %s",
						res.StatusCode,
						raw["error"].(map[string]interface{})["type"],
						raw["error"].(map[string]interface{})["reason"],
					)
				}
				// A successful response might still contain errors for particular documents...
			} else {
				if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
					log.Fatalf("Failure to to parse response body: %s", err)
				} else {
					for _, d := range blk.Items {
						// any HTTP status above 201 is the error
						if d.Index.Status > 201 {
							numErrors++
							// print the response status and error information ...
							log.Printf("  Error: [%d]: %s: %s: %s: %s",
								d.Index.Status,
								d.Index.Error.Type,
								d.Index.Error.Reason,
								d.Index.Error.Cause.Type,
								d.Index.Error.Cause.Reason,
							)
						} else {
							// ... otherwise increase the success counter.
							numIndexed++
						}
					}
				}
			}

			// Close the response body, to prevent reaching the limit for goroutines or file handles
			res.Body.Close()

			// Reset the buffer and items counter
			buf.Reset()
			numItems = 0
		}
	}
	if numIndexed == count {
		return nil
	} else {
		return errors.New("some places weren't uploaded")
	}
}

func (es *ElasticAPI) GetPlaces(limit int, offset int) ([]api.Place, int, error) {
	var (
		batchNum  int
		scrollID  string
		totalHits int
	)
	res, err := es.Client.Search(
		es.Client.Search.WithIndex("places"),
		es.Client.Search.WithSort("_doc"),
		es.Client.Search.WithSize(limit),
		es.Client.Search.WithScroll(time.Minute),
	)
	if err != nil {
		return []api.Place{}, 0, err
	}

	// Handle the first batch of data and extract the scrollID
	jsonRes := read(res.Body)
	res.Body.Close()
	totalHits = int(gjson.GetBytes(jsonRes, "hits.total.value").Int())
	if offset < 0 || totalHits < (offset*limit) {
		return []api.Place{}, 0, errors.New("Invalid 'page' value: " + strconv.Itoa(offset))
	}

	for {
		// if we scrolled for needed page
		if batchNum == offset {
			hits := gjson.GetBytes(jsonRes, "hits.hits.#._source").Array()
			array := make([]api.Place, len(hits))
			for i, hit := range hits {
				err = json.Unmarshal([]byte(hit.String()), &array[i])
			}
			return array, totalHits, err
		} else {
			scrollID = gjson.GetBytes(jsonRes, "_scroll_id").String()
		}

		// Perform the scroll request and pass the scrollID and scroll duration
		res, err := es.Client.Scroll(es.Client.Scroll.WithScrollID(scrollID),
			es.Client.Scroll.WithScroll(time.Minute))
		if err != nil {
			break
		}
		if res.IsError() {
			err = errors.New(res.String())
			break
		}
		batchNum++
		jsonRes = read(res.Body)
		res.Body.Close()
	}
	return []api.Place{}, totalHits, err
}

func (es *ElasticAPI) GetRecommendations(lat float64, lon float64, size int) ([]api.Place, error) {
	/*
			Query sort in Elasticsearch.
			It will produce query like this:
			"sort": [
		    {
		      "_geo_distance": {
		        "location": {
		          "lat": 55.674,
		          "lon": 37.666
		        },
		        "order": "asc",
		        "unit": "km",
		        "mode": "min",
		        "distance_type": "arc",
		        "ignore_unmapped": true
		      }
		    }
			]
	*/
	var buf bytes.Buffer
	sort := map[string]interface{}{
		"sort": map[string]interface{}{
			"_geo_distance": map[string]interface{}{
				"location": map[string]interface{}{
					"lat": lat,
					"lon": lon,
				},
				"order":           "asc",
				"unit":            "km",
				"mode":            "min",
				"distance_type":   "arc",
				"ignore_unmapped": true,
			},
		},
	}

	// We encode from map string-interface into json format.
	json.NewEncoder(&buf).Encode(sort)

	// Process the query
	search, _ := es.Client.Search(
		es.Client.Search.WithSize(size),
		es.Client.Search.WithIndex(es.indexName),
		es.Client.Search.WithBody(&buf),
	)
	jsonSearch := read(search.Body)
	search.Body.Close()
	hits := gjson.GetBytes(jsonSearch, "hits.hits.#._source").Array()
	array := make([]api.Place, len(hits))
	for i, hit := range hits {
		json.Unmarshal([]byte(hit.String()), &array[i])
	}
	return array, nil
}

func read(r io.Reader) []byte {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.Bytes()
}
