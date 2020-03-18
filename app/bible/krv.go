package bible

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
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

const (
	DB_FILE = "kofirst.db"
)

func (k *Krv) Init() {

	d, err := sql.Open("sqlite3", DB_FILE)
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
	verses := make([]string, 0)
	// vs, err := toSearchableFormat(k.Input.Memories)
	// if err != nil {
	// 	panic(err)
	// }
	for _, v := range k.Input.Memories {
		p := k.Bible.ParseVerses(toValidParam(v))
		verses = append(verses, p)
	}
	k.MemoryVerses = verses

	verses = make([]string, 0)
	// vs, err = toSearchableFormat(k.Input.Verses)
	// if err != nil {
	// 	panic(err)
	// }

	for _, v := range k.Input.Verses {
		p := k.Bible.ParseVerses(toValidParam(v))
		verses = append(verses, p)
	}
	k.Verses = verses
}

func (k Krv) Generate() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("#  %s\n", k.Input.Title))

	if len(k.MemoryVerses) > 0 {
		b.WriteString(fmt.Sprintln("\n## Memory Verses"))
		for _, v := range k.MemoryVerses {
			v = strings.Replace(v, "\n\n", " ", -1)
			b.WriteString(fmt.Sprintf("- %s\n", v))
		}
	}
	b.WriteString(fmt.Sprintln("\n## Verses"))
	for _, v := range k.Verses {
		v = strings.Replace(v, "\n\n", " ", -1)
		b.WriteString(fmt.Sprintf("- %s\n", v))
	}

	return b.String()
}
