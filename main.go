package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// returning JSON data of searched title -> joining a predetermined addOptions then -> fed into Post request
// 1st part
func searchTitle(name string, url string, api string) (map[string]interface{}, error) {

	req, err := http.NewRequest("GET", url+"lookup?term="+name, nil)
	if err != nil {
		fmt.Printf("New request for title [%s] is invalid: %s\n", name, err)
		return nil, err
	}

	req.Header.Set("X-Api-Key", api)
	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		fmt.Printf("Lookup request for title [%s] did not connect: %s\n", name, err)
		return nil, err
	}

	Body := res.Body
	resBody, err := io.ReadAll(Body)
	if err != nil {
		fmt.Printf("Couldn't parse search results from []bytes: %s\n", err)
		return nil, err
	}

	var lookup []map[string]interface{}

	err = json.Unmarshal(resBody, &lookup) //says []byte is redundant so removed it
	if err != nil {
		fmt.Printf("JSON unmarshal error: %s", err)
		return nil, err
	}

	mainTitle := lookup[0] //first search result

	return mainTitle, nil
}

func pushTitle(mainTitle map[string]interface{}, url string, api string, addOps map[string]interface{}) error {

	for key, value := range addOps {
		mainTitle[key] = value
	}
	payload, err := json.Marshal(&mainTitle)
	if err != nil {
		return err
	}

	req2, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req2.Header.Set("X-Api-Key", api)
	req2.Header.Set("Content-Type", "application/json")
	cli := &http.Client{}
	res2, err := cli.Do(req2)
	if err != nil {
		return err
	}

	fmt.Println(res2.StatusCode)

	return nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please provide the title as a command-line argument.")
		return
	}

	title := os.Args[1]
	titleDir := "/movies/"
	url := "http://localhost:7878/api/v3/movie/"
	api := "769ba9a660bf419b9c1ad868f7149b06"

	addOps := map[string]interface{}{
		"qualityProfileId": 1,
		"RootFolderPath":   titleDir,
		"monitored":        true,
		"addOptions": map[string]interface{}{"searchForMovie": true,
			"monitor": "movieOnly"},
	}

	results, _ := searchTitle(title, url, api)
	err := pushTitle(results, url, api, addOps)
	if err != nil {
		fmt.Printf("%s", err)
	}
}
