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
	fmt.Printf("dialTLS - 1\n")
	
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	fmt.Printf("dialTLS - 2\n")
	
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{ServerName: host}
	fmt.Printf("dialTLS - 3\n")
	
	tlsConn := tls.Client(conn, cfg)
	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		return nil, err
	}
	fmt.Printf("dialTLS - 4\n")
	
	cs := tlsConn.ConnectionState()
	cert := cs.PeerCertificates[0]
	fmt.Printf("dialTLS - 5\n")
	
	// Verify here
	cert.VerifyHostname(host)
	log.Println(cert.Subject)
	fmt.Printf("dialTLS - 6\n")
	
	return tlsConn, nil
}


func main() {
	// Getting Configuration from json file
	//config := GetConfig()

	fmt.Printf("1\n")

	
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		//Host:   config.AppiotURL,
		Host:   "kddiappiot.sensbysigma.com",
	})
	fmt.Printf("2\n")
	
	// Set a custom DialTLS to access the TLS connection state
	proxy.Transport = &http.Transport{DialTLS: dialTLS}
	fmt.Printf("3\n")
	
	// Change req.Host so badssl.com host check is passed
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = req.URL.Host
	}
	fmt.Printf("4\n")
	return;
	
	//logFile := "testlogfile"
	port := "8888"
	//port += config.Port
    if os.Getenv("HTTP_PLATFORM_PORT") != "" {
        //logFile = "D:\\home\\site\\wwwroot\\testlogfile"
        port = os.Getenv("HTTP_PLATFORM_PORT")
	}
	fmt.Printf("5\n")
	
	fmt.Printf("starting server\n")

	//log.Fatal(http.ListenAndServeTLS(port, "ericsson/iot/resources/certificate.pem", "ericsson/iot/resources/key.pem", proxy))
	log.Fatal(http.ListenAndServe(":"+port, proxy))
	fmt.Printf("6\n")
}

