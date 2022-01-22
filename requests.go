package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type requestparams struct {
	Method        string
	URL           string
	Body          string
	Token         string
	DefaultErr    string
	ContinueOnErr bool
	Decode        interface{}
}

func request(params *requestparams) error {
	req, _ := http.NewRequest(params.Method, params.URL, strings.NewReader(params.Body))
	if params.Method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Private-Token", params.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if params.ContinueOnErr {
			return err
		}
		println(params.DefaultErr)
		println(err.Error())
		os.Exit(-1)
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		if params.ContinueOnErr {
			return err
		}
		println("Not found.")
		os.Exit(-1)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		if params.ContinueOnErr {
			return err
		}
		println(params.DefaultErr)
		println(res.StatusCode)
		os.Exit(-1)
	}
	if params.Decode != nil {
		err := json.NewDecoder(res.Body).Decode(params.Decode)
		if err != nil {
			if params.ContinueOnErr {
				return err
			}
			println(params.DefaultErr)
			println(err.Error())
			os.Exit(-1)
		}
	}
	return nil
}
