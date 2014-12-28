package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	str "strings"
)

var baseString = "https://play.google.com/store/apps/details?id="

func main() {

	// testPackage := "com.zwenexsys.yoteshin"

	// Using with file
	f, err := os.Open("ys.html")
	PanicIf(err)
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	// doc, err := goquery.NewDocument(fmt.Sprintf("%s%s", baseString, testPackage))
	PanicIf(err)

	//	doc.Find(".main-content").Each(func(i int, s *goquery.Selection) {
	//		title := s.Find(".document-title").Text()
	//		desc := s.Find(".id-app-orig-desc").Text()
	//		fmt.Printf("%s - %s", title, desc)
	//
	//	})

	tmp := make(map[string]string)
	doc.Find(".meta-info").Each(func(i int, s *goquery.Selection) {
		fieldName := str.TrimSpace(s.Find(".title").Text())
		switch fieldName {
		case "Updated":
			tmp["updated"] = s.Find(".content").Text()
		case "Installs":
			tmp["installs"] = s.Find(".content").Text()
		case "Size":
			tmp["size"] = s.Find(".content").Text()
		case "Current Version":
			tmp["currentVersion"] = s.Find(".content").Text()
		case "Requires Android":
			tmp["requiresAndroid"] = s.Find(".content").Text()
		case "Content Rating":
			tmp["contentRating"] = s.Find(".content").Text()
		case "Developer":
			// Ugly hack
			s.Find(".dev-link").Each(func(i int, t *goquery.Selection) {
				nodeHref, _ := t.Attr("href")
				if str.Contains(nodeHref, "mailto:") {
					tmp["email"] = str.Split(nodeHref, "mailto:")[1]
				} else {
					raw := str.Split(nodeHref, "&")[0]
					tmp["websiteURL"] = str.Split(raw, "q=")[1]
				}
			})
		}
	})

	tmp["category"] = str.TrimSpace(doc.Find(".category").First().Text())
	for x, y := range tmp {
		fmt.Printf("%s - %s\n", x, str.TrimSpace(y))
	}

}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
