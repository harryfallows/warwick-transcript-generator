package main

import (
    "golang.org/x/net/html"
    "io"
	"strings"
	"io/ioutil"
	"log"
	"fmt"
)

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

				}	else if tokClass == "mod-name" {

					curTok = htmlTok.Next()
					modules[curModCode]["Name"] = htmlTok.Token().Data

				}	else if tokClass == "mod-reg-summary-item" {

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

				fmt.Println(yearInfo)

			}
			

			
		}

		
	}

}

func main() {

	htm, err := ioutil.ReadFile("test_data/test.htm")

    if err != nil {
        log.Fatal(err)
    }

    text := string(htm)
    doc := strings.NewReader(text)
	modules, yearInfo := scrape(doc)
	GenerateLatex(modules, yearInfo)
    
}
