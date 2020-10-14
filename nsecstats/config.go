package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/spf13/pflag"
)

type configuration struct {
	verbose  int
	statsdir string
	data     []measurement
}

var config configuration

func getConfig() {
	var datafile string
	var err error

	// define and parse command line arguments
	pflag.CountVarP(&config.verbose, "verbose", "v", "print more information while running")
	pflag.StringVarP(&config.statsdir, "statsdir", "s", "", "directory to save statistics")
	pflag.StringVarP(&datafile, "data", "d", "", "file (json) with nsec3 test results")
	pflag.Parse()

	if len(config.statsdir) == 0 {
		log.Println("command line option statsdir must be given")
	}
	config.statsdir, err = filepath.Abs(config.statsdir)
	if err != nil {
		log.Fatal(err)
	}
	if config.verbose > 0 {
		log.Println("Statsdir ", config.statsdir)
	}
	if len(datafile) == 0 {
		log.Println("command line option data must be given")
	}

	if config.verbose > 0 {
		log.Println("Reading test results from ", datafile)
	}
	dat, err := ioutil.ReadFile(datafile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(dat, &config.data)
	if err != nil {
		log.Fatal(err)
	}
	if config.verbose > 0 {
		log.Println("Done getConfig()")
	}
}
