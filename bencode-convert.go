package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	bencode "github.com/jackpal/bencode-go"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	format = kingpin.Flag("format", "Output format, either JSON or YAML (default).").Default("yaml").String()
	file   = kingpin.Arg("file", "Input file to parse.").Required().String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	dat, err := os.Open(*file)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"File":  *file,
		}).Fatal("Error opening file.")
	}

	decoded, err := bencode.Decode(dat)
	dat.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"File":  *file,
		}).Fatal("Could not decode.")
	}

	switch strings.ToLower(*format) {
	case "json":
		json, err := json.MarshalIndent(decoded, "", "  ")
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Fatal("Could not convert to JSON.")
		}
		fmt.Printf("%s\n", json)
	case "yaml":
		yaml, err := yaml.Marshal(decoded)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Fatal("Could not convert to YAML.")
		}
		fmt.Printf("%s\n", yaml)
	default:
		log.WithFields(log.Fields{
			"Format": *format,
		}).Fatal("Unknown format.")
	}
}
