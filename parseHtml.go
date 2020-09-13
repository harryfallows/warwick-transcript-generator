package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

type postRequest struct {
	Name       string            `json:"name"`
	Course     string            `json:"course"`
	Grad       string            `json:"grad"`
	Logo       string            `json:"logo"`
	Disclaimer bool              `json:"disclaimer"`
	Reverse    bool              `json:"reverse"`
	Files      map[string]string `json:"files"`
}

//Retrieves all information about all modules
func scrape(htmlReader io.Reader) (modules map[string]map[string]string, yearInfo map[string]string) {

	htmlTok := html.NewTokenizer(htmlReader)

	var curModCode string

	yearInfo = make(map[string]string)
	modules = make(map[string]map[string]string)

	for {

		curTok := htmlTok.Next()

		switch {

		case curTok == html.ErrorToken:
			return modules, yearInfo

		case curTok == html.StartTagToken:

			cT := htmlTok.Token()
			isSpan := cT.Data == "span"
			isH3 := cT.Data == "h3"

			if isSpan {

				tokClass := cT.Attr[0].Val

				if tokClass == "mod-code" {

					curTok = htmlTok.Next()
					curModCode = htmlTok.Token().Data
					modules[curModCode] = make(map[string]string, 5)

				} else if tokClass == "mod-name" {

					curTok = htmlTok.Next()
					modules[curModCode]["Name"] = htmlTok.Token().Data

				} else if tokClass == "mod-reg-summary-item" {

					htmlTok.Next()
					htmlTok.Next()
					modAttr := strings.TrimSuffix(htmlTok.Token().Data, ":")
					htmlTok.Next()
					htmlTok.Next()
					modAttrVal := strings.TrimPrefix(htmlTok.Token().Data, " ")
					modules[curModCode][modAttr] = modAttrVal

				}

			} else if isH3 {

				if cT.Attr != nil {

					continue

				}

				htmlTok.Next()

				for i := 1; i < 4; i++ {
					htmlTok.Next()
					htmlTok.Next()
					//fmt.Println(htmlTok.Token().Data)
					yearAttr := strings.TrimSuffix(htmlTok.Token().Data, ":")
					htmlTok.Next()
					htmlTok.Next()
					//fmt.Println(htmlTok.Token().Data)
					yearVal := strings.TrimPrefix(strings.TrimSuffix(htmlTok.Token().Data, " "), " ")
					yearInfo[yearAttr] = yearVal
				}

			}

		}

	}

}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var pReq postRequest
	err = json.Unmarshal(body, &pReq)

	row_template, _ := ioutil.ReadFile("resources/row_template.tex")

	var tables string
	years := make([]string, 0)
	for y, _ := range pReq.Files {
		years = append(years, y)
	}
	sort.Strings(years)
	if pReq.Reverse {
		for i, j := 0, len(years)-1; i < j; i, j = i+1, j-1 {
			years[i], years[j] = years[j], years[i]
		}
	}
	for _, y := range years {
		doc := strings.NewReader(pReq.Files[y])
		modules, yearInfo := scrape(doc)
		tables += GenerateTable(string(row_template), y, modules, yearInfo)
	}

	var info string
	if pReq.Name != "" {
		info += `\textbf{Name:} ` + pReq.Name + `\newline`
	}
	if pReq.Course != "" {
		info += `\textbf{Course Name:} ` + pReq.Course + `\newline`
	}
	if pReq.Grad != "" {
		info += `\textbf{Graduation Year:} ` + pReq.Grad + `\newline`
	}

	template, _ := ioutil.ReadFile("resources/template.tex")
	file := strings.Replace(string(template), "{tables}", tables, 1)
	file = strings.Replace(file, "{info}", info, 1)
	file = strings.Replace(file, "{logo}", pReq.Logo, 1)
	ioutil.WriteFile("output/temp.pdf", compileLatex(file), 0644)
}

// main function, mainly used to test atm
func main() {

	http.HandleFunc("/post", handleRequest)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
