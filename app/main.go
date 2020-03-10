package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sky4access/bible/app/bible"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	configFile := ""
	lang := ""

	flag.StringVar(&configFile, "config", "config.yaml", "config yaml file")
	flag.StringVar(&lang, "lang", "eng", "language: eng or kor")
	flag.Parse()

	if lang != "eng" && lang != "kor" {
		fmt.Println("language must be either eng or kor")
		os.Exit(1)
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("%s does not exist\n", configFile)
		os.Exit(1)
	}

	if _, err := os.Stat(bible.DB_FILE); os.IsExist(err) {
		fmt.Printf("database file %s does not exist\n", bible.DB_FILE)
	}

	var b bible.Bible

	if lang == "eng" {
		b = &bible.Esv{ConfigFile: configFile}
	} else {
		b = &bible.Krv{ConfigFile: configFile}
	}

	b.Init()
	b.Fetch()
	fmt.Print(b.Generate())

}
