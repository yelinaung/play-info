package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

var baseString = "https://play.google.com/store/apps/details?id="

func main() {

	testPackage := "com.google.android.gm"
	doc, err := goquery.NewDocument(fmt.Sprintf("%s%s", baseString, testPackage))
	PanicIf(err)

	doc.Find(".main-content").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".document-title").Text()
		desc := s.Find(".id-app-orig-desc").Text()
		fmt.Printf("%s - %s", title, desc)
	})
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
