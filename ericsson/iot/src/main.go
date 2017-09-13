package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
)

type Config struct {
	Port      string `json:"port"`
	AppiotURL string `json:"appiot_url"`
}

func GetConfig() *Config{
	file, e := ioutil.ReadFile("ericsson/iot/resources/config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(file))

	var con Config
	json.Unmarshal(file, &con)
	return &con
}

func dialTLS(network, addr string) (net.Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{ServerName: host}

	tlsConn := tls.Client(conn, cfg)
	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		return nil, err
	}

	cs := tlsConn.ConnectionState()
	cert := cs.PeerCertificates[0]

	// Verify here
	cert.VerifyHostname(host)
	log.Println(cert.Subject)

	return tlsConn, nil
}


func main() {
	// Getting Configuration from json file
	//config := GetConfig()

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		//Host:   config.AppiotURL,
		Host:   "kddiappiot.sensbysigma.com",

	})

	// Set a custom DialTLS to access the TLS connection state
	proxy.Transport = &http.Transport{DialTLS: dialTLS}

	// Change req.Host so badssl.com host check is passed
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = req.URL.Host
	}

	port := ":"
	//port += config.Port
	port += "8888"

	fmt.Printf("starting server\n")

	//log.Fatal(http.ListenAndServeTLS(port, "ericsson/iot/resources/certificate.pem", "ericsson/iot/resources/key.pem", proxy))
	log.Fatal(http.ListenAndServe(port, proxy))
}

