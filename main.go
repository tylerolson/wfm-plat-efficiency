package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const API = "https://api.warframe.market/v2"

func main() {
	resp, err := http.Get(fmt.Sprintf("%v/items", API))
	if err != nil {
		log.Fatalln(err)
	}

	decoder := json.NewDecoder(resp.Body)
	var response ItemsResponse
	err = decoder.Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.ApiVersion)
}
