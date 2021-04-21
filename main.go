package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
	// "reflect"
)

type whoisInfo struct {
	// Values returned with Error
	Error  bool   `json:"error"`
	Reason string `json:"reason"`
	// Values returned with valid results
	Asn                string  `json:"asn"`
	City               string  `json:"city"`
	ContinentCode      string  `json:"continent_code"`
	Country            string  `json:"country"`
	CountryCallingCode string  `json:"country_calling_code"`
	CountryName        string  `json:"country_name"`
	Currency           string  `json:"currency"`
	InEu               bool    `json:"in_eu"`
	IP                 string  `json:"ip"`
	Languages          string  `json:"languages"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Org                string  `json:"org"`
	Postal             string  `json:"postal"`
	Region             string  `json:"region"`
	RegionCode         string  `json:"region_code"`
	Timezone           string  `json:"timezone"`
	UtcOffset          string  `json:"utc_offset"`
}

func myIP() string {
	getURL := "http://icanhazip.com/"
	res, err := http.Get(getURL)
	defer res.Body.Close()
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	currentIP := strings.Trim(string(body), "\n")

	// fmt.Println(currentIP)

	return currentIP
}

func ipWhoIs(testIP string) whoisInfo {
	if !checkIP(testIP) {
		logErr := fmt.Sprintf("%s is not a valid IP", testIP)
		log.Fatal(logErr)
	}
	whoIsURL := fmt.Sprintf("https://ipapi.co/%s/json/", testIP)

	req, err := http.NewRequest("GET", whoIsURL, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Golang_Geo_Check/3.0")

	client := &http.Client{Timeout: time.Second * 10}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	var ipinfoResult whoisInfo
	json.Unmarshal(body, &ipinfoResult)

	return ipinfoResult

}

func checkIP(testIP string) bool {
	ipConv := net.ParseIP(testIP)
	// fmt.Println(ipConv)
	if ipConv == nil {
		return false
	}
	return true
}

func main() {

	var inFlag bool
	flag.BoolVar(&inFlag, "h", false, "Display help")
	flag.Parse()
	if inFlag == true {
		fmt.Println("\nRun the program with geoip <IP ADDRESS>")
		fmt.Println("If run without an IP address, the program will determine your public IP and use that.")
		fmt.Println("Have a nice day....\n")
		return
	}

	arg := flag.Arg(0)
	if arg == "" {
		arg = myIP()
	}

	whoisResult := ipWhoIs(arg)

	// fmt.Println(whoisResult)

	if whoisResult.Error {
		log.Fatal(whoisResult.Reason)
	}

	ipHost, _ := net.LookupAddr(arg)
	var resolvedHost string
	if len(ipHost) == 0 {
		resolvedHost = ""
	} else {
		resolvedHost = ipHost[0]
	}

	// mapURL := fmt.Sprintf("https://wego.here.com/%v,%v", whoisResult.Latitude, whoisResult.Longitude)
	mapURL := fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%v,%v", whoisResult.Latitude, whoisResult.Longitude)
	space := " "

	var printOut string
	printOut = fmt.Sprintf("\nIP:%-10s%s\n", space, whoisResult.IP)
	printOut += fmt.Sprintf("Host:%8s%s\n", space, resolvedHost)
	printOut += fmt.Sprintf("Location:%-4s%s / %s / %s\n", space, whoisResult.City, whoisResult.Region, whoisResult.CountryName)
	printOut += fmt.Sprintf("Map:%9s%s\n", space, mapURL)
	printOut += fmt.Sprintf("Org.:%8s%s\n", space, whoisResult.Org)
	printOut += fmt.Sprintf("ASN:%9s%s\n\n", space, whoisResult.Asn)

	fmt.Println(printOut)

}
