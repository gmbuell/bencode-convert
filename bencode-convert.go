package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	bencode "github.com/jackpal/bencode-go"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	format = kingpin.Flag("format", "Output format, either JSON or YAML (default).").Default("yaml").String()
	path   = kingpin.Arg("path", "Input file or url to parse.").Required().String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	var dat io.ReadCloser
	// Try and parse path as a url
	url, err := url.Parse(*path)
	if err != nil || url.Scheme == "" {
		log.WithFields(log.Fields{
			"Error": err,
			"File":  *path,
		}).Info("Could not parse as a url.")

		// Open a local file
		dat, err = os.Open(*path)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
				"File":  *path,
			}).Fatal("Error opening file.")
		}
	} else {
		var client http.Client
		resp, err := client.Get(*path)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
				"URL":   url,
			}).Fatal("Error getting url.")
		}
		dat = resp.Body
	}

	decoded, err := bencode.Decode(dat)
	dat.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"File":  *path,
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
