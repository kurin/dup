package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kurin/dup/update"
)

var config = filepath.Join(os.Getenv("HOME"), ".dup")

var (
	remote   = flag.String("remote", "https://domains.google.com/nic/update", "The remote API endpoint.")
	host     = flag.String("host", "", "The hostname to update.")
	ip       = flag.String("ip", "", "The IP to set.  Leave blank to use the connection IP.")
	offline  = flag.Bool("offline", false, "Set the host 'offline' status.")
	username = flag.String("username", "", "The API username.")
	password = flag.String("password", "", "The API password.")
	save     = flag.Bool("save", false, "Save these settings.")
)

func loadJSON(file string) (*update.Config, error) {
	var c update.Config
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func saveJSON(file string) error {
	c := &update.Config{
		Remote:   *remote,
		Host:     *host,
		IP:       *ip,
		Offline:  *offline,
		Username: *username,
		Password: *password,
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, b, 0600)
}

func main() {
	flag.Parse()
	if *save {
		if err := saveJSON(config); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	c, err := loadJSON(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	c.Update()
}
