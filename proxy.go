package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	Debugf("Received URI " + r.URL.Path)
	if fileExists(r.URL.Path) {
		body, err := ioutil.ReadFile("/tmp/forge/tracywebtech/pip/tracywebtech-pip-1.3.4/metadata.json")
		if err != nil {
			Fatalf(err.Error())
		}
		fmt.Fprint(w, string(body))
	} else {
		Fatalf("Need to get " + r.URL.Path)
	}
}

// getMetadataForgeModule queries the configured Puppet Forge and return
func getMetadataForgeModule(uri string) {
	if len(config.ForgeUrl) > 0 {
	}
	url := config.ForgeUrl + "/v3/releases/" + uri
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "https://github.com/xorpaul/g10k/")
	req.Header.Set("Connection", "close")
	proxyURL, err := http.ProxyFromEnvironment(req)
	if err != nil {
		Fatalf("getMetadataForgeModule(): Error while getting http proxy with golang http.ProxyFromEnvironment()" + err.Error())
	}
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	before := time.Now()
	Debugf("GETing " + url)
	resp, err := client.Do(req)
	duration := time.Since(before).Seconds()
	Verbosef("GETing Forge metadata from " + url + " took " + strconv.FormatFloat(duration, 'f', 5, 64) + "s")
	mutex.Lock()
	syncForgeTime += duration
	mutex.Unlock()
	if err != nil {
		Fatalf("getMetadataForgeModule(): Error while querying " + url + ": " + err.Error())
	}
	defer resp.Body.Close()

	if resp.Status == "200 OK" {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			Fatalf("getMetadataForgeModule(): Error while reading response body for " + url + ": " + err.Error())
		}

		writer, err := os.Create(uri)

		if err != nil {
			Fatalf(funcName + "(): error while Create() " + uri + err.Error())
		}

		io.Copy(writer, body)

		if err != nil {
			Fatalf(funcName + "(): error while io.Copy() " + uri + err.Error())
		}

		writer.Close()

		mutex.Lock()
		forgeJsonParseTime += duration
		mutex.Unlock()

	} else {
		Fatalf("getMetadataForgeModule(): Unexpected response code while GETing " + url + " " + resp.Status)
	}
}
