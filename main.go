package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)
type CheckKeyAPI struct {
	Url *string
	HttpClient *http.Client
	Condition *string
	InputKey *string
	InputValue *string
	ResultValue *string
}
func main() {
	cka := CheckKeyAPI{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
			},
		},
	}
	cka.run()
}
func (g *CheckKeyAPI) run() {
	g.prepare()
	g.takeBody()
	g.doChecking()
}
func (g *CheckKeyAPI) parseInput() {
	url := flag.String("url", os.Getenv("URL"), "an url!")
	keyInput := flag.String("key", os.Getenv("KEY"), "a json key!")
	valueInput := flag.String("value", os.Getenv("VALUE"), "a json value!")
	condition := flag.String("condition", os.Getenv("CONDITION"), "a condition: equal or unequal!")
	flag.Parse()

	if *url == "" {
		g.msgAndExit(fmt.Sprintf("Please input with -%s or set %s env", "url", "URL"))
	}
	if *keyInput == "" {
		g.msgAndExit(fmt.Sprintf("Please input with -%s or set %s env", "key", "KEY"))
	}
	if *valueInput == "" {
		g.msgAndExit(fmt.Sprintf("Please input with -%s or set %s env", "value", "VALUE"))
	}
	if *condition == "" {
		g.msgAndExit(fmt.Sprintf("Please input with -%s or set %s env", "condition", "CONDITION"))
	}
	g.setUrl(url)
	g.setKey(keyInput)
	g.setValue(valueInput)
	g.setCondition(condition)
}
func (g *CheckKeyAPI) msgAndExit(msg string) {
	log.Fatalf("Error: %s", msg)
}

func (g *CheckKeyAPI) setUrl(v *string) {
	g.Url = v
}

func (g *CheckKeyAPI) setKey(v *string) {
	g.InputKey = v
}

func (g *CheckKeyAPI) setValue(v *string) {
	g.InputValue = v
}
func (g *CheckKeyAPI) setCondition(v *string) {
	g.Condition = v
}
func (g *CheckKeyAPI) prepare() {
	g.parseInput()
	g.checkInput()
}
func (g *CheckKeyAPI) checkInput()  {
	if g.checkCondition() == false {
		g.msgAndExit("Condition must be only equal or unequal")
	}
}
func (g *CheckKeyAPI) checkCondition() bool {
	switch *g.Condition {
	case "equal","unequal":
		return true
	default:
		return false
	}
}
func (g *CheckKeyAPI) httpGet(url *string) *[]byte {
	req, err := http.NewRequest(http.MethodGet, *url, nil)
	g.failOnError(err, "Error setting http request")
	resp, err := g.HttpClient.Do(req)
	g.failOnError(err, "Error getting url")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	g.failOnError(err, "Error reading all")
	return &body
}

func (g *CheckKeyAPI) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (g *CheckKeyAPI) takeBody() {
	// Call API
	body := g.httpGet(g.Url)
	// Define a map to parse json
	m := map[string]interface{}{}
	err := json.Unmarshal(*body, &m)
	g.failOnError(err, "Error decode json")
	value := fmt.Sprintf("%v", m[*g.InputKey])
	g.ResultValue = &value
}
func (g *CheckKeyAPI) doChecking() {
	// Check equal
	equal := strings.EqualFold(*g.InputValue,*g.ResultValue)
	if equal {
		if *g.Condition == "equal" {
			// As expected
			log.Printf("The value of %s is equal with %s as you want", *g.InputKey, *g.ResultValue)
			os.Exit(0)
		} else if *g.Condition == "unequal" {
			// As unexpected
			log.Printf("The value of %s is equal with %s but you want %s", *g.InputKey, *g.ResultValue, *g.Condition)
			os.Exit(1)
		}
	} else {
		if *g.Condition == "equal" {
			// As unexpected
			log.Printf("The value of %s is unequal with %s but you want %s", *g.InputKey, *g.ResultValue, *g.Condition)
			os.Exit(1)
		} else if *g.Condition == "unequal" {
			// As expected
			log.Printf("The value of %s is unequal with %s as you want", *g.InputKey, *g.ResultValue)
			os.Exit(0)
		}
	}
}