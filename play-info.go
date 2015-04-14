package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/cli"
	"github.com/mgutz/ansi"
	"os"
	"sort"
	. "strings"
)

var baseString = "https://play.google.com/store/apps/details?id="

// for my own
const debug = false

var divider = fmt.Sprintf(ansi.Color(Repeat("-", 56)+"\n", "yellow"))

func main() {
	app := cli.NewApp()
	app.Name = "Play Go"
	app.Usage = "Get app info via commandline"
	app.Version = "0.1.0"

	app.Action = func(c *cli.Context) {
		arg := c.Args()
		if len(arg) == 0 {
			fmt.Println(ansi.Color("Error : Please enter package name", "red"))
		} else {
			fmt.Println("\n")
			pkg := c.Args()[0]
			fmt.Printf(ansi.Color("Processing Results for \"%s\"\n", "green"), pkg)
			fmt.Println(divider)
			GetData(c.Args()[0])
			fmt.Println("\n")
			fmt.Println(divider)
		}
	}

	app.Run(os.Args)
}

func GetData(pkgName string) {
	// Using with file
	var doc *goquery.Document
	var err error
	if debug {
		f, err := os.Open("karrency.html")
		PanicIf(err)
		defer f.Close()
		doc, err = goquery.NewDocumentFromReader(f)
	} else {
		doc, err = goquery.NewDocument(fmt.Sprintf("%s%s", baseString, pkgName))
		PanicIf(err)
	}

	tmp := make(map[TitleMap]string)

	tmp[TitleMap{1, "Title"}] = doc.Find(`div[itemprop='name']`).First().Text()
	tmp[TitleMap{2, "Category"}] = TrimSpace(doc.Find(".category").First().Text())

	price := TrimSpace(doc.Find(".price").First().Text())
	if price == "Install" {
		tmp[TitleMap{3, "Price"}] = "Free"
	} else {
		tmp[TitleMap{3, "Price"}] = TrimSpace(doc.Find(".price").First().Text())
	}

	doc.Find(".meta-info").Each(func(i int, s *goquery.Selection) {
		fieldName := TrimSpace(s.Find(".title").Text())
		switch fieldName {
		case "Updated":
			tmp[TitleMap{4, "Updated"}] = s.Find(".content").Text()
		case "Installs":
			tmp[TitleMap{5, "Total Installs"}] = s.Find(".content").Text()
		case "Size":
			tmp[TitleMap{6, "Size"}] = s.Find(".content").Text()
		case "Current Version":
			tmp[TitleMap{7, "Current Version"}] = s.Find(".content").Text()
		case "Requires Android":
			tmp[TitleMap{8, "Requires Android"}] = s.Find(".content").Text()
		case "Content Rating":
			tmp[TitleMap{9, "Content Rating"}] = s.Find(".content").Text()
		case "Developer":
			// Ugly hack
			s.Find(".dev-link").Each(func(i int, t *goquery.Selection) {
				nodeHref, _ := t.Attr("href")
				if Contains(nodeHref, "mailto:") {
					tmp[TitleMap{10, "Email"}] = Split(nodeHref, "mailto:")[1]
				} else {
					raw := Split(nodeHref, "&")[0]
					tmp[TitleMap{11, "Website"}] = Split(raw, "q=")[1]
				}
			})
		}
	})

	score := doc.Find(".score-container").First()
	if score != nil {
		tmp[TitleMap{12, "Score"}] = score.Find(".score").First().Text()
		node := doc.Find(`meta[itemprop='ratingCount']`)
		v, _ := node.Attr("content")
		tmp[TitleMap{13, "Votes"}] = v
	}

	tmp[TitleMap{14, "Developer"}] = TrimSpace((doc.Find(`div[itemprop='author']`).Find(".primary").Text()))
	tmp[TitleMap{15, "What's New"}] = doc.Find(".whatsnew .recent-change").Text()

	tmp[TitleMap{16, "Description"}] = doc.Find(`div[itemprop='description']`).First().Text()

	tmp[TitleMap{17, "App Id"}] = pkgName
	tmp[TitleMap{18, "Icon Url"}], _ = doc.Find(".cover-image").Attr("src")

	// Go iteration order is randomzies
	// https://blog.golang.org/go-maps-in-action#TOC_7.

	var keys ByIndex

	for k := range tmp {
		keys = append(keys, k)
	}

	// sort the keys
	sort.Sort(keys)

	for _, k := range keys {
		var rows string
		rows = fmt.Sprintf("%s %s | %s\n", k.Title, buffer(k.Title, 11), TrimSpace(tmp[k]))
		fmt.Printf(rows)
	}
}

// I copied it from https://github.com/addyosmani/psi/blob/master/lib%2Futils.js#L36-L50
func buffer(msg string, length int) string {
	var ret = ""

	length = 24

	length = length - len(msg) - 1

	if length > 0 {
		ret = Repeat(" ", length)
	}

	return ret
}

type ByIndex []TitleMap

func (a ByIndex) Len() int           { return len(a) }
func (a ByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type TitleMap struct {
	Index int
	Title string
}
