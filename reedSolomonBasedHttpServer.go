package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/klauspost/reedsolomon"
)

type PrivateKey struct {
	PrivateKey string
}

type ShardData struct {
	ShardData [][]byte
}

func Temp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Done"})
}

func generateAndStoreErasureCodeShards(privateKey string) error {
	if _, err := os.Stat("./ShardsAndParityBlockData"); err == nil {
		os.RemoveAll("./ShardsAndParityBlockData") // To Remove Old Shards File
	}
	os.Mkdir("./ShardsAndParityBlockData", 0777) // To Generate Old Shards Directory
	enc, err := reedsolomon.New(6, 3)
	if err != nil {
		return fmt.Errorf("error generated while creating reed solomon encoder: %v", err)
	}
	data := make([][]byte, 9)
	for i := range data {
		data[i] = make([]byte, 11)
	}
	curPos := 0
	for _, item := range data[:6] {
		for j := range item {
			item[j] = byte(privateKey[curPos])
			curPos += 1
		}
	}

	err = enc.Encode(data)
	if err != nil {
		return fmt.Errorf("error generated while creating encoding data: %v", err)
	}

	for i := 0; i < 6; i++ {
		fileName := fmt.Sprint("./ShardsAndParityBlockData/shard", i)
		file, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("error generated while creating file for shards: %v", err)
		}
		tempData := "["
		for i, item := range data[i] {
			if i == 10 {
				tempData += (strconv.Itoa(int(item)))
			} else {
				tempData += (strconv.Itoa(int(item)) + ", ")
			}
		}
		tempData += "]"
		if _, err = file.WriteString(tempData); err != nil {
			return fmt.Errorf("error generated while storing shards in shard file: %v", err)
		}
	}
	for i := 0; i < 3; i++ {
		fileName := fmt.Sprint("./ShardsAndParityBlockData/parityBlock", i)
		file, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("error generated while creating file for shards: %v", err)
		}
		log.Println(string(data[i+6]))
		tempData := "["
		for i, item := range data[i+6] {
			if i == 10 {
				tempData += (strconv.Itoa(int(item)))
			} else {
				tempData += (strconv.Itoa(int(item)) + ", ")
			}
		}
		tempData += "]"
		if _, err = file.WriteString(tempData); err != nil {
			return fmt.Errorf("error generated while storing shards in parityBlock file: %v", err)
		}
	}
	log.Println(data)

	return nil
}

func RecoverPrivateKey(c *gin.Context) {
	var shardData ShardData
	if err := c.BindJSON(&shardData); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid Input"})
		return
	}
	enc, err := reedsolomon.New(6, 3)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": fmt.Errorf("error generated while creating reed solomon encoder: %v", err)})
		return
	}

	data := shardData.ShardData
	log.Println(data)
	// c.JSON(http.StatusOK, gin.H{"message": "Done", "recoveredPrivateKey": data})
	err = enc.Reconstruct(data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Done", "recoveredPrivateKey": "Error While Regenerating Data"})
	}
	log.Println("\n\n", "Regenerated Private Key Data", data)
	recoveredPrivateKey := ""
	for _, item := range data[:6] {
		recoveredPrivateKey += string(item)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Done", "recoveredPrivateKey": recoveredPrivateKey})
}

func GenerateErasureCode(c *gin.Context) {
	var userPrivateKeyData PrivateKey
	if err := c.BindJSON(&userPrivateKeyData); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid Input"})
		return
	}
	log.Println(userPrivateKeyData)
	if len(userPrivateKeyData.PrivateKey) != 66 {
		c.JSON(http.StatusExpectationFailed, gin.H{"message": "Private Key Must Be Have 66 Characters"})
		return
	}
	if err := generateAndStoreErasureCodeShards(userPrivateKeyData.PrivateKey); err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Done"})
}

func main() {
	router := gin.Default()
	router.GET("/", Temp)
	router.POST("/generateErasureCode", GenerateErasureCode)
	router.POST("/recoverPrivateKey", RecoverPrivateKey)
	log.Println("Ok")
	router.Run("localhost:8000")

}

/*
To Generate Erasure Code Of Private Key,
http://localhost:8000/generateErasureCode
sampleData:-
{
  "privateKey":"0x0123456789012345678901234567890123456789012345678901234567890123"
}
To Recover Private Key,
http://localhost:8000/recoverPrivateKey
Note:- Share Data Should Be Passed In Order Manner,
You Can See Below Given Data Which Is In Order Manner Where As Data Of Index 3,4,5 Is Missing And Because Of That Empty Array Is Passed.
sampleData:-
{
  "ShardData":[
[48, 120, 48, 49, 50, 51, 52, 53, 54, 55, 56],
[57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57],
[48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48],
[],
[],
[],
[2, 213, 54, 39, 48, 49, 58, 31, 62, 63, 58],
[11, 153, 55, 32, 49, 58, 59, 26, 63, 58, 59],
[42, 2, 56, 194, 135, 84, 140, 121, 200, 155, 67]
  ]
}
*/
