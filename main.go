package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"encoding/json"
)

type Configuration struct {
	Token    string
}

func readConf(conf_fn string) Configuration  {
	file, _ := os.Open(conf_fn)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	return configuration
}

type TelegramResponse struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date     int    `json:"date"`
			Text     string `json:"text"`
			Entities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
			} `json:"entities"`
		} `json:"message"`
	} `json:"result"`
}

var clients []int = []int{}

func contains(slice []int, item int) bool {
	set := make(map[int]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func parse_response(content string)  {
	var response TelegramResponse

	if err := json.Unmarshal([]byte(content), &response); err != nil {
		panic(err)
	}

	fmt.Println(response.Ok)

	for i := 0; i<len(response.Result) ; i++ {
		client_id := response.Result[i].Message.From.ID

		if response.Result[i].Message.Text == "/start" {
			if !contains(clients, client_id) {
				clients = append(clients, client_id)
			}
		}
	}

	fmt.Println(clients)
}

func main() {
	config := readConf("conf.json")
	fmt.Println("Token:", config.Token)

	var url string = "https://api.telegram.org/bot" + config.Token + "/getUpdates"

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(contents))

		parse_response(string(contents))
	}
}