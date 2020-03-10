package bible

import (
	"fmt"
	"net/http"
	"os"
)

const (
	API_KEY = "ESV_API_KEY"
)

type Input struct {
	Title    string
	Memories []string
	Verses   []string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func callAPI(url string, param string) (*http.Response, error) {
	key := os.Getenv(API_KEY)
	if len(key) == 0 {
		panic(fmt.Errorf("No Environment variable: %s", API_KEY))
	}

	path := fmt.Sprintf("%s%s", url, param)

	req, err := http.NewRequest("GET", path, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", key))
	client := &http.Client{}
	response, err := client.Do(req)
	return response, err
}

type Bible interface {
	Init()
	Fetch()
	Generate() string
}
