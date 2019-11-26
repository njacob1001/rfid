package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	rfid "github.com/firmom/go-rfid-rc522/rfid"
	rc522 "github.com/firmom/go-rfid-rc522/rfid/rc522"
)

const triggerPin int8 = 7
const echoPin int8 = 11

func main() {

	// RFID sensor
	product := ""
	allowNew := true
	products := []string{"test"}

	type PostBody struct {
		Articles []string `json:"articles"`
	}

	type ResponseAPI struct {
		Ok      bool   `json:"ok"`
		Message string `json:"message"`
	}

	fmt.Println("RFID STARTED")
	reader, err := rc522.NewRfidReader()
	if err != nil {
		fmt.Println(err)
		return
	}
	readerChan, err := rfid.NewReaderChan(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	contador := 0
	rfidChan := readerChan.GetChan()
	for product == "" {
		if contador > 5 && len(products) > 1 {
			allowNew = false
			resp := PostBody{
				Articles: products,
			}
			js, err := json.Marshal(resp)
			if err != nil {
				fmt.Println(err)
				allowNew = true
				products = []string{"test"}
				contador = 0
				continue
			}
			// Get domain information from SSLLabs API
			hostInfo, err := http.Post("http://13.59.72.139:80/api/user/sale", "application/json", bytes.NewBuffer(js))
			if err != nil {
				fmt.Println(err)
				allowNew = true
				products = []string{"test"}
				contador = 0
				continue
			}
			defer hostInfo.Body.Close()

			var hostResponse ResponseAPI
			if err := json.NewDecoder(hostInfo.Body).Decode(&hostResponse); err != nil {
				fmt.Println(err)
			}
			if hostResponse.Ok {
				fmt.Println("Compra exitosa!")
			} else {
				fmt.Println("Compra fracasada")
			}
			contador = 0
			products = []string{"test"}
			allowNew = true
		}
		if allowNew {
			select {
			case id := <-rfidChan:
				// product = id
				products = append(products, id)
			default:
				contador++
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}

	fmt.Println(product)
}
