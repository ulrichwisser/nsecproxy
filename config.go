package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/spf13/pflag"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	Verbose       int
	UpstreamNSEC  string
	UpstreamNSEC3 string
	IPlist        []string
}

func getConfig() *Configuration {
	var config Configuration
	var conffilename string

	// define and parse command line arguments
	pflag.StringVar(&conffilename, "conf", "", "Filename to read configuration from")
	pflag.CountVarP(&config.Verbose, "verbose", "v", "print more information while running")
	pflag.StringSliceVarP(&config.IPlist, "ip", "ip", []string{}, "IPv4 or IPv6 address(es) to listen on")
	pflag.Parse()

	var confFromFile *Configuration
	if conffilename != "" {
		var err error
		confFromFile, err = readConfigFile(conffilename)
		if err != nil {
			panic(err)
		}
	}

	defaultConfig := readDefaultConfigFiles()
	return checkConfiguration(joinConfig(defaultConfig, joinConfig(confFromFile, &config)))
}

func readConfigFile(filename string) (*Configuration, error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &Configuration{}
	err = yaml.Unmarshal(source, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func readDefaultConfigFiles() (config *Configuration) {

	// config in user home directory
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	fileconfig, err := readConfigFile(path.Join(usr.HomeDir, ".nsecproxy"))
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	config = joinConfig(config, fileconfig)

	// done
	return
}

func joinConfig(oldConf *Configuration, newConf *Configuration) (config *Configuration) {
	if oldConf == nil && newConf == nil {
		return nil
	}
	if oldConf != nil && newConf == nil {
		return oldConf
	}
	if oldConf == nil && newConf != nil {
		return newConf
	}

	// we have two configs, join them
	config = &Configuration{}
	config.Verbose = newConf.Verbose
	if newConf.UpstreamNSEC != "" {
		config.UpstreamNSEC = newConf.UpstreamNSEC
	} else {
		config.UpstreamNSEC = oldConf.UpstreamNSEC
	}
	if newConf.UpstreamNSEC3 != "" {
		config.UpstreamNSEC3 = newConf.UpstreamNSEC3
	} else {
		config.UpstreamNSEC3 = oldConf.UpstreamNSEC3
	}
	if len(newConf.IPlist) > 0 {
		config.IPlist = newConf.IPlist
	} else {
		config.IPlist = oldConf.IPlist
	}

	// Done
	return config
}

func checkConfiguration(config *Configuration) *Configuration {
	if len(config.UpstreamNSEC) == 0 {
		log.Fatal("NSEC servers must be given.")
	}
	if len(config.UpstreamNSEC3) == 0 {
		log.Fatal("NSEC3 Servers must be given.")
	}
	if len(config.IPlist) == 0 {
		log.Fatal("IP addresses must be given.")
	}

	// Done
	return config
}
