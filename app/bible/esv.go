package bible

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Esv struct {
	ConfigFile   string
	Input        Input
	MemoryVerses []string
	Verses       []string
}

func (e *Esv) Init() {
	source, err := ioutil.ReadFile(e.ConfigFile)
	check(err)

	err = yaml.Unmarshal(source, &e.Input)
	if err != nil {
		panic(err)
	}

}

func fetchVerse(input []string) ([]string, error) {

	type result struct {
		Query     string
		Canonical string
		Passages  []string
	}

	passages := make([]string, 0)
	for _, param := range input {
		response, err := callAPI(ESV_API_URL, toValidParam(param))
		if err != nil {
			return nil, err
		}

		var res result
		err = json.NewDecoder(response.Body).Decode(&res)
		if err != nil {
			panic(err)
		}

		if len(res.Passages) > 0 {
			passages = append(passages, res.Passages[0])
		}

		time.Sleep(30 * time.Millisecond)

	}
	return passages, nil
}

func (e *Esv) Fetch() {

	fetchedVerse, err := fetchVerse(e.Input.Memories)
	if err != nil {
		panic(err)
	}
	e.MemoryVerses = fetchedVerse

	fetchedVerse, err = fetchVerse(e.Input.Verses)
	if err != nil {
		panic(err)
	}
	e.Verses = fetchedVerse
}

func (e Esv) Generate() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("#  %s\n", e.Input.Title))
	if len(e.MemoryVerses) > 0 {
		b.WriteString(fmt.Sprintln("\n## Memory Verses"))
		for _, v := range e.MemoryVerses {
			v = strings.Replace(v, "\n\n", " ", -1)
			b.WriteString(fmt.Sprintf("- %s\n", v))
		}
	}
	b.WriteString(fmt.Sprintln("\n## Verses"))
	for _, v := range e.Verses {
		v = strings.Replace(v, "\n\n", " ", -1)
		b.WriteString(fmt.Sprintf("- %s\n", v))
	}

	return b.String()
}
