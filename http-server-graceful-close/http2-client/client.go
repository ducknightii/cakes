package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/net/http2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const url = "https://ducky.night.mac:8000/healthz"

var client *http.Client

func main() {
	client = &http.Client{}
	// Create a pool with the server certificate since it is not signed
	// by a known CA
	caCert, err := ioutil.ReadFile("/Users/hanlei/Documents/work/IJ/aibee_code/src/github.com/ducknightii/cakes/http-server-graceful-close/http2-server/server.crt")
	if err != nil {
		log.Fatalf("Reading server certificate: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	// Create TLS configuration with the certificate of the server
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	// Use the proper transport in the client
	/*client.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}*/
	client.Transport = &http2.Transport{
		TLSClientConfig: tlsConfig,
	}

	for i := 0; i < 10; i++ {
		go healthz()
		time.Sleep(time.Second * 3)
	}
	time.Sleep(time.Second * 60)

}

func healthz() {
	fmt.Println(time.Now(), "curl...")
	// Perform the request
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Failed get: %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed reading response body: %s", err)
		return
	}
	fmt.Printf(
		"Got response %d: %s %s\n",
		resp.StatusCode, resp.Proto, string(body),
	)
}
