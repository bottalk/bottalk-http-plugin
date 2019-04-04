package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	bottalk "github.com/bottalk/go-plugin"
)

type btRequest struct {
	Token   string          `json:"token"`
	UserID  string          `json:"user"`
	Input   json.RawMessage `json:"input"`
	URL     string          `json:"url"`
	Payload string          `json:"payload"`
	Headers []string        `json:"headers"`
}

func errorResponse(message string) string {
	return "{\"result\": \"fail\",\"message\":\"" + message + "\"}"
}

func doAction(method string, r *http.Request) string {

	var BTR btRequest
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&BTR)
	if err != nil {
		return errorResponse(err.Error())
	}

	log.Println(BTR.URL)

	if len(BTR.URL) < 5 {
		return errorResponse("Url is not defined")
	}

	req, err := http.NewRequest(strings.ToUpper(method), BTR.URL, bytes.NewBuffer([]byte(BTR.Payload)))

	req.Header.Set("Content-Type", "application/json")
	for _, hd := range BTR.Headers {
		if strings.Contains(hd, ":") {
			req.Header.Set(strings.Split(hd, ":")[0], strings.Split(hd, ":")[1])
		}
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return errorResponse(err.Error())
	}
	output, _ := ioutil.ReadAll(res.Body)
	log.Println(string(output))

	return "{\"result\": \"ok\",\"status\":" + fmt.Sprintf("%d", res.StatusCode) + ",\"response\":" + string(output) + "}"
}

func main() {

	plugin := bottalk.NewPlugin()
	plugin.Name = "HTTP Plugin"
	plugin.Description = "This plugin performs http queries with different methods"

	plugin.Actions = map[string]bottalk.Action{
		"get": bottalk.Action{
			Name:        "get",
			Description: "This performs GET-request to remote endpoint",
			Endpoint:    "/get",
			Action: func(r *http.Request) string {
				return doAction("get", r)
			},
			Params: map[string]string{"url": "Endpoint url to call", "headers": "Array of headers to send"},
		},
		"post": bottalk.Action{
			Name:        "post",
			Description: "This performs POST-request to remote endpoint",
			Endpoint:    "/post",
			Action: func(r *http.Request) string {
				return doAction("post", r)
			},
			Params: map[string]string{"url": "Endpoint url to call", "payload": "Payload to send in json format", "headers": "Array of headers to send"},
		},
		"patch": bottalk.Action{
			Name:        "patch",
			Description: "This performs PATCH-request to remote endpoint",
			Endpoint:    "/patch",
			Action: func(r *http.Request) string {
				return doAction("patch", r)
			},
			Params: map[string]string{"url": "Endpoint url to call", "payload": "Payload to send in json format", "headers": "Array of headers to send"},
		},
		"delete": bottalk.Action{
			Name:        "delete",
			Description: "This performs DELETE-request to remote endpoint",
			Endpoint:    "/delete",
			Action: func(r *http.Request) string {
				return doAction("delete", r)
			},
			Params: map[string]string{"url": "Endpoint url to call", "payload": "Payload to send in json format", "headers": "Array of headers to send"},
		},
	}

	plugin.Run(":9064")
}
