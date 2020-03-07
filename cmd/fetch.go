package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	API_KEY = "ESV_API_KEY"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type BibleInput struct {
	Title  string
	Url    string
	Verses []string
}

type BibleOutput struct {
	Query     string
	Canonical string
	Passages  []string
}

func main() {
	source, err := ioutil.ReadFile("config.yaml")
	check(err)

	var bible BibleInput
	err = yaml.Unmarshal(source, &bible)
	if err != nil {
		panic(err)
	}

	space := regexp.MustCompile(`\s+`)

	passages := make([]string, 0)
	for _, v := range bible.Verses {
		v = space.ReplaceAllString(v, " ")
		v = strings.Replace(v, " ", "+", -1)

		key := os.Getenv(API_KEY)
		path := fmt.Sprintf("%s?q=%s", bible.Url, v)

		req, err := http.NewRequest("GET", path, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", key))
		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		var result BibleOutput
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			panic(err)
		}
		for _, v := range result.Passages {
			passages = append(passages, v)
		}
	}
	fmt.Println(passages)
}
