package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/docopt/docopt-go"
	"io/ioutil"
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
	req.Header.Set("Content-Type", "text/plain; charset=UTF-8")
	req.SetBasicAuth(dmUser, dmPasswd)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(Exit{1})
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		os.Exit(2)
	}
	if result["status"] != "success" {
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

func main() {
	usage := `Domain Master DNS01 acme exec provider

Usage:
	dm-dns01 present <fqdn> <txt>
	dm-dns01 cleanup <fqdn> <txt>

Options:
	-h --help     Show this screen.
`
	args, _ := docopt.ParseDoc(usage)

	defer handleExit()
	var fqdn = strings.Split(args["<fqdn>"].(string), ".")
	var name = strings.Join(fqdn[:len(fqdn)-3], ".")
	var domain = strings.Join(fqdn[len(fqdn)-3:len(fqdn)-1], ".")

	if args["present"] == true {
		addTxtRecord(name, domain, args["<txt>"].(string))
	} else if args["cleanup"] == true {
		delTxtRecord(name, domain)
	}
}
