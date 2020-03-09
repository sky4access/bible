package bible

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	//bible_db "github.com/sky4access/bible/pkg/bible"
	"gopkg.in/doug-martin/goqu.v4"
	"gopkg.in/yaml.v2"
)

type Krv struct {
	ConfigFile   string
	Bible        BibleDB
	Input        Input
	MemoryVerses []string
	Verses       []string
}

func (k *Krv) Init() {

	d, err := sql.Open("sqlite3", "kofirst.db")
	if err != nil {
		panic(err)
	}
	k.Bible = BibleDB{goqu.New("sqlite3", d)}

	source, err := ioutil.ReadFile(k.ConfigFile)
	check(err)

	err = yaml.Unmarshal(source, &k.Input)
	if err != nil {
		panic(err)
	}
}

func (k *Krv) Fetch() {

	space := regexp.MustCompile(`\s+`)

	verses := make([]string, 0)
	for i, v := range k.Input.Memories {
		v = strings.Replace(space.ReplaceAllString(v, " "), " ", "+", -1)
		k.Input.Memories[i] = v
		p := k.Bible.ParseVerses(v)
		verses = append(verses, p)
	}
	k.MemoryVerses = verses

	verses = make([]string, 0)
	for i, v := range k.Input.Verses {
		v = strings.Replace(space.ReplaceAllString(v, " "), " ", "+", -1)
		k.Input.Verses[i] = v
		p := k.Bible.ParseVerses(v)
		verses = append(verses, p)
	}
	k.Verses = verses
}

func (k Krv) Print() {
	fmt.Printf("#  %s\n", k.Input.Title)
	fmt.Println("## Memory Verses")
	for _, v := range k.MemoryVerses {
		v = strings.Replace(v, "\n\n", " ", -1)
		fmt.Printf("- %s\n", v)
	}
	fmt.Println("\n\n## Verses")
	for _, v := range k.Verses {
		v = strings.Replace(v, "\n\n", " ", -1)
		fmt.Printf("- %s\n", v)
	}
	fmt.Println("\n")
}
