package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/docopt/docopt-go"
	"golang.org/x/net/publicsuffix"
	"io"
	"net/http"
	"os"
	"strings"
)

const dmApi = "https://www.domainmaster.cz/masterapi/server.php"

var dmUser = os.Getenv("DM_API_USER")
var dmPasswd = os.Getenv("DM_API_PASSWD")

type Exit struct{ Code int }

func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok == true {
			os.Exit(exit.Code)
		}
		panic(e) // not an Exit, bubble up
	}
}

func sendCommand(command string, params string) map[string]interface{} {
	var jsonStr = []byte(`{"command":"` + command + `","params":` + params + `}`)
	req, err := http.NewRequest("POST", dmApi, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Fprintln(os.Stderr, "dm-dns01: build request:", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "text/plain; charset=UTF-8")
	req.SetBasicAuth(dmUser, dmPasswd)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dm-dns01: request failed:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Fprintln(os.Stderr, "response Status:", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dm-dns01: read response:", err)
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Fprintln(os.Stderr, "dm-dns01: bad JSON response:", string(body))
		os.Exit(2)
	}
	if result["status"] != "success" {
		fmt.Fprintln(os.Stderr, "dm-dns01: API error:", string(body))
		os.Exit(3)
	}
	return result
}

func addTxtRecord(name string, domain string, txt string) {
	var params = `{"domain":"` + domain + `","name":"` + name + `","type":"TXT","data":"` + txt + `"}`
	sendCommand("create dns record", params)
}

func delTxtRecord(name string, domain string) {
	var params = `{"domain":"` + domain + `"}`
	var result map[string]interface{}
	var id = ""

	result = sendCommand("list dns records", params)
	var records = result["data"].([]interface{})
	for _, item := range records {
		item := item.(map[string]interface{})
		if item["name"] == name {
			id = item["id"].(string)
			break
		}
	}

	if id != "" {
		var params = `{"domain":"` + domain + `","id":"` + id + `"}`
		sendCommand("delete dns record", params)
	}
}

// splitFQDN splits a lego-supplied challenge FQDN (which carries a trailing dot,
// e.g. "_acme-challenge.grafana.kralovi.net.") into the record name and the
// registrable domain. It uses the public suffix list so multi-label suffixes
// (e.g. example.co.uk) resolve to the correct registrable domain.
func splitFQDN(fqdn string) (name string, domain string) {
	fqdn = strings.TrimSuffix(fqdn, ".")
	domain, err := publicsuffix.EffectiveTLDPlusOne(fqdn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dm-dns01: cannot determine registrable domain:", err)
		os.Exit(65)
	}
	name = strings.TrimSuffix(strings.TrimSuffix(fqdn, domain), ".")
	return name, domain
}

func main() {
	if dmUser == "" || dmPasswd == "" {
		fmt.Fprintln(os.Stderr, "dm-dns01: DM_API_USER and DM_API_PASSWD must be set")
		os.Exit(64)
	}

	usage := `Domain Master DNS01 acme exec provider

Usage:
	dm-dns01 present <fqdn> <txt>
	dm-dns01 cleanup <fqdn> <txt>

Options:
	-h --help     Show this screen.
`
	args, _ := docopt.ParseDoc(usage)

	defer handleExit()
	name, domain := splitFQDN(args["<fqdn>"].(string))

	if args["present"] == true {
		addTxtRecord(name, domain, args["<txt>"].(string))
	} else if args["cleanup"] == true {
		delTxtRecord(name, domain)
	}
}
