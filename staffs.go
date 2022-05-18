// staffs.go -- FernUniversität in Hagen StaffSearch consultation
// SPDX-License-Identifier: UNLICENSED
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/anaskhan96/soup"
)

const cName = "staffs.go" //if len(os.Args) > 0 { name = filepath.Base(os.Args[0]) }
const cVersion = "1.0.0"
const cCopyright = "Copyright (C) by XA, IX 2021. All rights reserved."

const urlEndpoint = `https://staffsearch.fernuni-hagen.de/`
const hdrUA = "Mozilla/5.0 (X11; U; Linux i686) Gecko/20071127 Firefox/2.0.0.11"

func PrintDataset(data, headers []string) {
	sep := ": "
	if headers == nil {
		headers = make([]string, len(data))
		for i := range headers {
			headers[i] = "•"
		}
		sep = " "
	}
	for i, d := range data {
		h := headers[i]
		fmt.Printf("%10s%s%s\n", h, sep, d)
	}
}

func preprocessData(data []byte) []byte {
	idxA := bytes.Index(data, []byte(`<p id="result"`))
	if idxA == -1 {
		return []byte{}
	}
	idxB := bytes.Index(data[idxA:], []byte(`<!--`))
	if idxB == -1 {
		idxB = len(data)
	} else {
		idxB += idxA
	}
	prefix := []byte("<root>\n")
	postfix := []byte("<root>\n")
	res := make([]byte, 0, idxB+len(prefix)+len(postfix))
	res = append(res, prefix...)
	res = append(res, data[idxA:idxB]...)
	res = append(res, postfix...)
	return res
}

func queryPerson(searchterm string, show_raw bool) {
	query := url.Values{}
	query.Add("s", searchterm)
	query.Add("p", "1")

	req, _ := http.NewRequest("GET", urlEndpoint+"?"+query.Encode(), nil)
	req.Header.Set("User-Agent", hdrUA)
	res, err := (&http.Client{}).Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Fprintln(os.Stderr, "- Error getting response.")
		log.Fatal(err)
		return
	}

	dataBytes, _ := io.ReadAll(res.Body)
	res.Body.Close()
	dataBytes = preprocessData(dataBytes)
	data := soup.HTMLParse(string(dataBytes))

	eltTable := data.Find("table")
	if eltTable.Error != nil {
		if eltTable.Error.(soup.Error).Type == soup.ErrElementNotFound {
			fmt.Fprintln(os.Stderr, "- No result found.")
		}
		return
	}
	eltsTR := eltTable.FindAll("tr")

	var dataHeaders []string
	for _, e := range eltsTR[0].FindAll("th") {
		dataHeaders = append(dataHeaders, e.FullText())
	}
	if show_raw {
		fmt.Fprintln(os.Stderr, dataHeaders)
	}

	dataResults := make([][]string, 0, 5) //, len(dataHeaders))
	for i, tr := range eltsTR[1:] {
		dataResults = append(dataResults, make([]string, 0, len(dataHeaders)))
		for _, e := range tr.FindAll("td") {
			dataResults[i] = append(dataResults[i], e.FullText())
		}
	}
	if show_raw {
		for _, res := range dataResults {
			fmt.Fprintln(os.Stderr, res)
		}
	}

	for _, res := range dataResults {
		PrintDataset(res, dataHeaders)
		fmt.Println()
	}

}

type argsSpec struct {
	Search_terms []string `arg:"positional,required" help:"the query terms"`
	Endpoint     string   `arg:"-e,--endpoint" help:"StaffSearch endpoint"`
	Raw          bool     `arg:"-r,--raw"`
}

func (argsSpec) Description() string {
	return cName + " " + cVersion +
		" – Query the StaffSearch service of the FernUni.\n" +
		cCopyright + "\n"
}

func parseArguments(dest *argsSpec) *arg.Parser {
	argParser, err := arg.NewParser(arg.Config{}, dest)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	var arguments []string
	if len(os.Args) > 0 { // discard Args[0] if Args not nil
		arguments = os.Args[1:]
	}
	err = argParser.Parse(arguments)
	switch {
	case err == arg.ErrHelp:
		argParser.WriteHelp(os.Stdout)
		os.Exit(0)
	case err == arg.ErrVersion:
		fmt.Printf("%s %s\n", cName, cVersion)
	case err != nil:
		argParser.Fail(err.Error())
		os.Exit(-1)
	}
	return argParser
}

func main() {

	var args argsSpec
	parseArguments(&args)

	for _, searchterms := range args.Search_terms {
		queryPerson(searchterms, args.Raw)
		fmt.Println()
	}
	os.Exit(0)
}
