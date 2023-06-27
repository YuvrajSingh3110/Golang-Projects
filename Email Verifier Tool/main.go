package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Enter a domanin: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		checkDomain(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error: ", err)
	}
}

func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var SPFrecords, DMARCrecords string

	MXrecord, err := net.LookupMX(domain)
	if err != nil {
		log.Print(err)
	}
	if len(MXrecord) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Print(err)
		SPFrecords = "nil"
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			SPFrecords = record
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Print(err)
		DMARCrecords = "nil"
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			DMARCrecords = record
		}
	}

	fmt.Println("Domain: ", domain)
	fmt.Println("has MX: ", hasMX)
	fmt.Println("has SPF: ", hasSPF)
	fmt.Println("SPF records: ", SPFrecords)
	fmt.Println("has DMARC: ", hasDMARC)
	fmt.Println("DMARC records: ", DMARCrecords)
}
