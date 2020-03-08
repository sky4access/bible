package main

import (
	"flag"

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

	b := bible.Esv{ConfigFile: configFile}
	b.Init()
	b.Fetch()
	b.Print()
}
