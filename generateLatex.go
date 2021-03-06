package main

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/rwestlund/gotex"
)

// generates tables based on scraped info
func GenerateTable(template string, year string, modules map[string]map[string]string, yearInfo map[string]string) (table string) {

	var rows string

	var keys []string
	for k, _ := range modules {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, code := range keys {

		info := modules[code]
		rows += code + "&" + strings.Replace(info["Name"], "&", `\&`, -1) + "&" + info["CATS"] + "&" + info["Mark"] + `\%&` + info["Grade"] + `\\` + `\hline `

	}

	var yearInfoStr string

	for yearAttr, yearVal := range yearInfo {

		yearInfoStr += `\multicolumn{5}{|c|}{\large{` + yearAttr + ": " + strings.Replace(yearVal, `%`, `\%`, -1) + `}} \\ \hline`

	}

	template = strings.Replace(template, "{year}", "Year "+year, 1)
	table = strings.Replace(template, "{rows}", rows, 1)
	table = strings.Replace(table, "{rows2}", yearInfoStr, 1)

	return table

}

// run latex compilation
func compileLatex(document string) (pdf []byte) {

	var err error
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	pdf, err = gotex.Render(document, gotex.Options{
		Command:   "",
		Runs:      1,
		Texinputs: path})

	if err != nil {
		log.Println("render failed ", err)
	}

	return pdf
}

/* // GenerateLatex ties together all latex generation related code
func GenerateTable(modules map[string]map[string]string, yearInfo map[string]string) string {



	return table

} */
