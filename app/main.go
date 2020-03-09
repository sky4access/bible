package main

import (
	"flag"

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
	flag.StringVar(&configFile, "config", "config.yaml", "config yaml file")
	flag.Parse()

	esv := bible.Esv{ConfigFile: configFile}
	esv.Init()
	esv.Fetch()
	esv.Print()

	krv := bible.Krv{ConfigFile: configFile}
	krv.Init()
	krv.Fetch()
	krv.Print()

}
