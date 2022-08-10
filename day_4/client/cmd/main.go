package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"io/ioutil"
	"log"
	"net/http"
	apiclient "swagger/code-gen/client"
	"swagger/code-gen/client/operations"
)

var (
	CAcertPath     string
	clientCertPath string
	keyPath        string
	candyType      string
	countCandy     int64
	money          int64
)

func init() {
	flag.StringVar(&CAcertPath, "ca", "../ca/minica.pem", "path for certificate of CA")
	flag.StringVar(&clientCertPath, "cert", "127.0.0.1/cert.pem", "path for certificate of client")
	flag.StringVar(&keyPath, "key", "127.0.0.1/key.pem", "path for key of client")
	flag.StringVar(&candyType, "k", "", "type of candy")
	flag.Int64Var(&countCandy, "c", 0, "number of candies in order")
	flag.Int64Var(&money, "m", 0, "amount of money were put in vending machine")
	flag.Parse()
}

func main() {
	certCA, err := ioutil.ReadFile(CAcertPath)
	certClient, err2 := ioutil.ReadFile(clientCertPath)
	keyClient, err3 := ioutil.ReadFile(keyPath)
	if err != nil || err2 != nil || err3 != nil {
		log.Fatalf("Can't open cert file of CA, cert or key")
	}
	cp, _ := x509.SystemCertPool()
	cp.AppendCertsFromPEM(certCA)

	// create the transport
	transport := httptransport.New("127.0.0.1:3333", "", []string{"https"})
	ptrTlsConfig := &transport.Transport.(*http.Transport).TLSClientConfig
	if *ptrTlsConfig == nil {
		cer, _ := tls.X509KeyPair(certClient, keyClient)
		*ptrTlsConfig = &tls.Config{
			RootCAs:      cp,
			Certificates: []tls.Certificate{cer},
		}
	}

	// create the API client, with the transport
	client := apiclient.New(transport, strfmt.Default)

	// make the request to get all items
	resp, err := client.Operations.BuyCandy(operations.NewBuyCandyParams().WithOrder(operations.BuyCandyBody{
		CandyCount: &countCandy,
		CandyType:  &candyType,
		Money:      &money,
	}))
	if err != nil {
		log.Fatal(err)
	}
	res, _ := json.MarshalIndent(*resp.Payload, "", "  ")
	fmt.Println(string(res))
}
