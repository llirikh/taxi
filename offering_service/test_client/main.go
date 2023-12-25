package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"offering_service/internal/api/responses"
	"offering_service/internal/models"
)

const url = "http://localhost:8080/"

func main() {
	req := models.Offer{From: models.Position{Lat: 15, Lng: 16}, To: models.Position{Lat: 1, Lng: 2}, Client_id: "kirill"}
	marshalled, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	resp, err := http.Post(url+"offers", "application/json", bytes.NewBuffer(marshalled))
	if err != nil {
		fmt.Println("request error:", err)
	}
	defer resp.Body.Close()
	var respo responses.CreateOfferResponse
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("reading error", err)
	}
	err = json.Unmarshal(result, &respo)
	fmt.Println(respo.Offer_id)
	jwwwt := respo.Offer_id

	resp, err = http.Get(url + "offers/" + jwwwt)
	if err != nil {
		fmt.Println("request error:", err)
	}
	defer resp.Body.Close()
	var respoo responses.ParseOfferResponse
	result, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("reading error", err)
	}
	err = json.Unmarshal(result, &respoo)
	fmt.Println(respoo)
}
