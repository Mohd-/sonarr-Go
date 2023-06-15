package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func pushTitle(mainTitle map[string]interface{}, url string, api string, titleDir string) error {

	addOps := map[string]interface{}{
		"QualityProfileId": 4,
		"RootFolderPath":   titleDir,
		"addOptions": map[string]interface{}{"searchForMovie": true,
			"monitor": "movieOnly"},
	}

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

}