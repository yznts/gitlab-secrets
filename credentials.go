package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

// Function to save credentials into config
func savecredentials(domain, token string) {
	// Extract current user
	usr, _ := user.Current()
	// Create config directory
	if err := os.MkdirAll(fmt.Sprintf("%s/.config/gitlab-secrets/", usr.HomeDir), 0777); err != nil {
		println("Can't create config directory.")
		println(err.Error())
		os.Exit(-1)
	}
	// Store credentials
	if err := ioutil.WriteFile(fmt.Sprintf("%s/.config/gitlab-secrets/%s", usr.HomeDir, domain), []byte(token), 0777); err != nil {
		println("Can't save token.")
		println(err.Error())
		os.Exit(-1)
	}
}

// Function to load credentials from config
func loadcredentials(domain string) string {
	// Extract current user
	usr, _ := user.Current()
	// Try to read credentials file
	tokenbts, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/gitlab-secrets/%s", usr.HomeDir, domain))
	if err != nil {
		println("Can't read token.")
		println(err.Error())
		os.Exit(-1)
	}
	return string(tokenbts)
}
