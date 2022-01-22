package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Commands

func Auth(domain, token string) {
	// Need to ensure token is valid
	requrl := fmt.Sprintf("https://%s/api/v4/user", domain)
	req, _ := http.NewRequest("GET", requrl, nil)
	req.Header.Set("PRIVATE-TOKEN", token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		println("Token validation request failed:")
		println(err.Error())
		os.Exit(-1)
	}
	if res.StatusCode != 200 {
		println("Can't authenticate with provided token")
		os.Exit(-1)
	}
	// Save credentials
	savecredentials(domain, token)
	// Message
	println("Successful authentication.")
	println(fmt.Sprintf("Credentials are storred into ~/.config/gitlab-secrets/%s", domain))
}

func KVList(values, meta bool) {
	// Extract repository url
	repourl := repourl()
	// Extract credentials
	domain := repourl.Host
	token := loadcredentials(domain)
	// Encode project path (will be used as id)
	projectid := url.PathEscape(strings.ReplaceAll(repourl.Path, ".git", "")[1:])
	// Extract variables from API
	data := []map[string]interface{}{}
	request(&requestparams{
		Method:     "GET",
		URL:        fmt.Sprintf("https://%s/api/v4/projects/%s/variables", domain, projectid),
		Token:      token,
		DefaultErr: "Error during variables extraction request.",
		Decode:     &data,
	})
	// Extract to key/value map
	kvmap := map[string]string{}
	kvmasked := map[string]bool{}
	kvprotected := map[string]bool{}
	for _, v := range data {
		kvmap[v["key"].(string)] = v["value"].(string)
		kvmasked[v["key"].(string)] = v["masked"].(bool)
		kvprotected[v["key"].(string)] = v["protected"].(bool)
	}
	// Format to env
	env := envencode(kvmap, !values)
	// Extend with meta
	if meta {
		outbuilder := strings.Builder{}
		for _, line := range strings.Split(env, "\n") {
			key := strings.Split(line, "=")[0]
			metadata := []string{}
			if kvmasked[key] {
				metadata = append(metadata, "masked")
			}
			if kvprotected[key] {
				metadata = append(metadata, "protected")
			}
			outbuilder.WriteString(fmt.Sprintf("%s (%s)\n", line, strings.Join(metadata, ",")))
		}
		env = outbuilder.String()
	}
	// Print
	println(strings.TrimSpace(env))
}

func KVGet(key string, keysilent, meta bool) {
	// Extract repository url
	repourl := repourl()
	// Extract credentials
	domain := repourl.Host
	token := loadcredentials(domain)
	// Encode project path (will be used as id)
	projectid := url.PathEscape(strings.ReplaceAll(repourl.Path, ".git", "")[1:])
	// Extract variable from API
	data := map[string]interface{}{}
	request(&requestparams{
		Method:     "GET",
		URL:        fmt.Sprintf("https://%s/api/v4/projects/%s/variables/%s", domain, projectid, key),
		Token:      token,
		DefaultErr: "Error during variable extraction request.",
		Decode:     &data,
	})
	// Output
	row := strings.Builder{}
	if !keysilent {
		row.WriteString(key)
		row.WriteString("=")
	}
	row.WriteString(`"` + data["value"].(string) + `"`)
	if meta {
		metadata := []string{}
		if data["protected"].(bool) {
			metadata = append(metadata, "protected")
		}
		if data["masked"].(bool) {
			metadata = append(metadata, "masked")
		}
		row.WriteString(fmt.Sprintf(" (%s)", strings.Join(metadata, ",")))
	}
	println(row.String())
}

func KVSet(key, value string, masked, protected bool) {
	// Extract repository url
	repourl := repourl()
	// Extract credentials
	domain := repourl.Host
	token := loadcredentials(domain)
	// Encode project path (will be used as id)
	projectid := url.PathEscape(strings.ReplaceAll(repourl.Path, ".git", "")[1:])
	// Try to create new variable with API
	reqbody := map[string]interface{}{}
	reqbody["key"] = key
	reqbody["value"] = value
	reqbody["masked"] = masked
	reqbody["protected"] = protected
	reqbodybts, _ := json.Marshal(reqbody)
	err := request(&requestparams{
		Method:        "POST",
		URL:           fmt.Sprintf("https://%s/api/v4/projects/%s/variables", domain, projectid),
		Token:         token,
		Body:          string(reqbodybts),
		DefaultErr:    "Error during variable setting request.",
		ContinueOnErr: true,
	})
	// In case of error, try to update existing variable
	if err != nil {
		err = request(&requestparams{
			Method:        "PUT",
			URL:           fmt.Sprintf("https://%s/api/v4/projects/%s/variables/%s", domain, projectid, key),
			Token:         token,
			Body:          string(reqbodybts),
			DefaultErr:    "Error during variable setting request.",
			ContinueOnErr: true,
		})
		if err != nil {
			println("Error during variable setting request.")
			println(err.Error())
			os.Exit(-1)
		}
	}
	// Output
	println("Variable was set.")
}

func KVDel(key string) {
	// Extract repository url
	repourl := repourl()
	// Extract credentials
	domain := repourl.Host
	token := loadcredentials(domain)
	// Encode project path (will be used as id)
	projectid := url.PathEscape(strings.ReplaceAll(repourl.Path, ".git", "")[1:])
	// Try to delete variable with API
	request(&requestparams{
		Method:     "DELETE",
		URL:        fmt.Sprintf("https://%s/api/v4/projects/%s/variables/%s", domain, projectid, key),
		Token:      token,
		DefaultErr: "Error during variable deletion request.",
	})
	println("Variable deleted.")
}

func Pull(file string, utrunc bool) {
	// Extract repository url
	repourl := repourl()
	// Extract credentials
	domain := repourl.Host
	token := loadcredentials(domain)
	// Encode project path (will be used as id)
	projectid := url.PathEscape(strings.ReplaceAll(repourl.Path, ".git", "")[1:])
	// Extract variables from API
	data := []map[string]interface{}{}
	request(&requestparams{
		Method:     "GET",
		URL:        fmt.Sprintf("https://%s/api/v4/projects/%s/variables", domain, projectid),
		Token:      token,
		DefaultErr: "Error during variables extraction request.",
		Decode:     &data,
	})
	// Extract to key/value map
	kvmap := map[string]string{}
	for _, v := range data {
		// Extract key/value
		key := v["key"].(string)
		value := v["value"].(string)
		// Truncate underscore for key
		if utrunc && strings.HasPrefix(key, "_") {
			key = key[1:]
		}
		// Add to map
		kvmap[key] = value
	}
	// Format to env
	env := envencode(kvmap, false)
	// Write to file
	err := ioutil.WriteFile(file, []byte(env), 0777)
	if err != nil {
		println("Error during file writing.")
		println(err.Error())
		os.Exit(-1)
	}
	// Output
	println("Variables pulled.")
}

func main() {
	// Help message
	if len(os.Args) == 1 {
		println("Commands:")
		println("  auth     - Authenticate with a private token")
		println("  kv:list  - List keys or key/value pairs")
		println("  kv:get   - Get a value with a key")
		println("  kv:set   - Set a key/value pair")
		println("  kv:del   - Delete a key/value pair")
		println("  pull     - Pull secrets from GitLab to specified location")
		println("  push     - Push secrets from specified location to GitLab")
		println("Usage:")
		println("  gitlab-secrets auth -token <token>")
		println("  gitlab-secrets kv:get -key <key>")
		println("  gitlab-secrets kv:set -key <key> -value <value>")
		println("  gitlab-secrets kv:del -key <key>")
		println("  gitlab-secrets pull -file <filepath>")
		println("  gitlab-secrets push -file <filepath>")
		return
	}

	// Auth command arguments
	fsauth := flag.NewFlagSet("auth", flag.ExitOnError)
	fsauthdomain := fsauth.String("domain", "gitlab.com", "Domain of GitLab instance")
	fsauthtoken := fsauth.String("token", "", "Authorization token")
	// Key/value list command arguments
	fskvlist := flag.NewFlagSet("kv:list", flag.ExitOnError)
	fskvlistvalues := fskvlist.Bool("values", false, "Reveal values")
	fskvlistmeta := fskvlist.Bool("meta", false, "Reveal meta data (masked, protected)")
	// Key/value get command arguments
	fskvget := flag.NewFlagSet("kv:get", flag.ExitOnError)
	fskvgetkey := fskvget.String("key", "", "Pair key")
	fskvgetksilent := fskvget.Bool("kslicent", false, "Remove key from output")
	fskvgetmeta := fskvget.Bool("meta", false, "Reveal meta data (masked, protected)")
	// Key/value set command arguments
	fskvset := flag.NewFlagSet("kv:set", flag.ExitOnError)
	fskvsetkey := fskvset.String("key", "", "Pair key")
	fskvsetval := fskvset.String("value", "", "Pair value")
	fskvsetmasked := fskvset.Bool("masked", false, "Mark as masked")
	fskvsetprotected := fskvset.Bool("protected", false, "Mark as protected")
	// Key/value del command arguments
	fskvdel := flag.NewFlagSet("kv:del", flag.ExitOnError)
	fskvdelkey := fskvdel.String("key", "", "Pair key")
	// Import command arguments
	fspull := flag.NewFlagSet("pull", flag.ExitOnError)
	fspullfile := fspull.String("file", ".env", "Pull target file")
	fspullutrunc := fspull.Bool("utrunc", false, "Truncate underscore on start for keys")
	// Export command arguments
	fspush := flag.NewFlagSet("push", flag.ExitOnError)
	fspushfile := fspush.String("file", ".env", "Push source file")
	fspushumirr := fspush.Bool("umirr", false, "Add underscored duplicate for each pair")
	// Determine command
	switch os.Args[1] {
	case "auth":
		fsauth.Parse(os.Args[2:])
		Auth(*fsauthdomain, *fsauthtoken)
	case "kv:list":
		fskvlist.Parse(os.Args[2:])
		KVList(*fskvlistvalues, *fskvlistmeta)
	case "kv:get":
		fskvget.Parse(os.Args[2:])
		KVGet(*fskvgetkey, *fskvgetksilent, *fskvgetmeta)
	case "kv:set":
		fskvset.Parse(os.Args[2:])
		KVSet(*fskvsetkey, *fskvsetval, *fskvsetmasked, *fskvsetprotected)
	case "kv:del":
		fskvdel.Parse(os.Args[2:])
		KVDel(*fskvdelkey)
	case "pull":
		fspull.Parse(os.Args[2:])
		Pull(*fspullfile, *fspullutrunc)
	}

	_ = fspushfile
	_ = fspushumirr
}
