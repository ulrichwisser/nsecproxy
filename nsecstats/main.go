package main

import (
	"log"
	"os"
)

func main() {

	// Command line parameters
	getConfig()

	// calculate all stats
	makeDir()
	makeRcode()
}

func makeDir() {
	if _, err := os.Stat(config.statsdir); os.IsNotExist(err) {
		err = os.MkdirAll(config.statsdir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

}
