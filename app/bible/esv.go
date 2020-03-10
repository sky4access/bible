package bible

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Esv struct {
	ConfigFile   string
	BaseUrl      string
	Input        Input
	MemoryVerses []string
	Verses       []string
}

func (e *Esv) Init() {
	e.BaseUrl = "https://api.esv.org/v3/passage/text/?include-headings=false&include-footnotes=false&include-short-copyright=false&q="

	source, err := ioutil.ReadFile(e.ConfigFile)
	check(err)

	err = yaml.Unmarshal(source, &e.Input)
	if err != nil {
		panic(err)
	}

}

func (e *Esv) Fetch() {
	type result struct {
		Query     string
		Canonical string
		Passages  []string
	}

	space := regexp.MustCompile(`\s+`)

	for i, v := range e.Input.Memories {
		v = strings.Replace(space.ReplaceAllString(v, " "), " ", "+", -1)
		e.Input.Memories[i] = v
	}

	param := strings.Join(e.Input.Memories, ",")

	response, err := callAPI(e.BaseUrl, param)
	if err != nil {
		panic(err)
	}

	var res result
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		panic(err)
	}
	passages := make([]string, 0)
	for _, v := range res.Passages {
		passages = append(passages, v)
	}
	e.MemoryVerses = passages

	for i, v := range e.Input.Verses {
		v = strings.Replace(space.ReplaceAllString(v, " "), " ", "+", -1)
		e.Input.Verses[i] = v
	}

	param = strings.Join(e.Input.Verses, ",")

	response, err = callAPI(e.BaseUrl, param)
	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		panic(err)
	}
	passages = make([]string, 0)
	for _, v := range res.Passages {
		passages = append(passages, v)
	}
	e.Verses = passages
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
