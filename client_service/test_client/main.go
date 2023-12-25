package main

import (
	"bytes"
	"client_service/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "http://localhost:8080/"

func main() {
	req := models.Offer{Price: models.Price{Amount: 345345, Currency: "RUB"}, From: models.Position{Lat: 15, Lng: 16}, To: models.Position{Lat: 1, Lng: 2}, Client_id: "kirill"}
	marshalled, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	client := &http.Client{}
	reqq, err := http.NewRequest("POST", url+"/trips", bytes.NewBuffer(marshalled))
	if err != nil {
		fmt.Println(err)
		return
	}
	reqq.Header.Add("Content-Type", "application/json")
	reqq.Header.Add("user_id", "123")

	resp, err := client.Do(reqq)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}
