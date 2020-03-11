package bible

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"
)

const (
	API_KEY     = "ESV_API_KEY"
	ESV_API_URL = "https://api.esv.org/v3/passage/text/?include-headings=false&include-footnotes=false&include-short-copyright=false&q="
)

type Bible interface {
	Init()
	Fetch()
	Generate() string
}

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

// {"john 1:1", "mark 2:1,5"} => {"john+1:1", "mark+2:1", "mark+2:5"}
func toSearchableFormat(verses []string) ([]string, error) {
	result := make([]string, 0)
	space := regexp.MustCompile(`\s+`)

	for _, v := range verses {
		v = strings.Replace(space.ReplaceAllString(v, " "), " ", "+", 1)
		arr := strings.Split(v, "+")
		if len(arr) < 2 {
			return nil, fmt.Errorf("%s is not valid book chapter:verse", v)
		}
		b, cv := arr[0], arr[1]
		cv = stripSpaces(cv)
		arr = strings.Split(cv, ":")
		if len(arr) < 2 {
			return nil, fmt.Errorf("%s is not valid chapter:verse", cv)
		}
		c, vs := arr[0], arr[1]

		vs = stripSpaces(vs)
		arr = strings.Split(vs, ",")

		for _, v1 := range arr {
			canonical := fmt.Sprintf("%s+%s:%s", b, c, v1)
			result = append(result, canonical)
		}

	}
	return result, nil
}

func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}
