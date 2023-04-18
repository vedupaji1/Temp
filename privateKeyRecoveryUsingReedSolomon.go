package main

import (
	"log"

	"github.com/klauspost/reedsolomon"
)

func main() {
	userPrivateKey := "0x0123456789012345678901234567890123456789012345678901234567890123"
	enc, err := reedsolomon.New(6, 3)
	if err != nil {
		log.Panicln("Error While Creating Reedsolomon Encoder", err)
	}
	data := make([][]byte, 9)
	for i := range data {
		data[i] = make([]byte, 11)
	}
	curPos := 0
	for _, item := range data[:6] {
		for j := range item {
			item[j] = byte(userPrivateKey[curPos])
			curPos += 1
		}
	}
	log.Println("\n\n", "Private Key Data", data)

	// Encoding Data, Here Parity Blocks Will Be Added In Data
	err = enc.Encode(data)
	if err != nil {
		log.Panicln("Error While Encoding Data", err)
	}
	log.Println("\n\n", "Encoded Private Key Data", data)

	// Verifying Encoded Data
	ok, err := enc.Verify(data)
	if err != nil {
		log.Panicln("Error While Verifying Data", err)
	}
	if ok {
		log.Println("Data Verification Done")
	} else {
		log.Panicln("Data Verification Failed")
	}

	// Removing 2 Data Shards
	data[1] = nil
	// data[5] = nil
	data[6] = nil
	data[7] = nil
	log.Println("\n\n", "Incomplete Private Key Data", data)

	// Regenerating Removed Or Destroyed Data
	err = enc.Reconstruct(data)
	if err != nil {
		log.Panicln("Error While Regenerating Data", err)
	}
	log.Println("\n\n", "Regenerated Private Key Data", data)
	recoveredPrivateKey := ""
	for _, item := range data[:6] {
		recoveredPrivateKey += string(item)
	}
	log.Println("User Original Private Key:- ", userPrivateKey)
	log.Println("User Recovered Private Key:- ", recoveredPrivateKey)
}
