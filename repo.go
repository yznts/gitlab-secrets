package main

import (
	"net/url"
	"os"
	"os/exec"
	"strings"
)

// Function to resolve current repository url
func repourl() *url.URL {
	// execute "git remote show origin" to extract remote info
	cmd := exec.Command("sh", "-c", "git remote show origin")
	outputbts, err := cmd.CombinedOutput()
	if err != nil {
		println("can't get repository remote info with git.")
		println(string(outputbts))
		println(err.Error())
		os.Exit(-1)
	}
	// Find and return
	for _, line := range strings.Split(string(outputbts), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "Fetch URL:") {
			// Cleanup
			line = strings.TrimSpace(line)
			line = strings.ReplaceAll(line, "Fetch URL: ", "")
			// Parse
			repourl, err := url.Parse(line)
			if err != nil {
				println("Can't parse repository url.")
				println(err.Error())
				os.Exit(-1)
			}
			// Return
			return repourl
		}
	}
	println("Can't find repository url.")
	os.Exit(-1)
	return nil
}
