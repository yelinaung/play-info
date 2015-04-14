package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/cli"
	"github.com/mgutz/ansi"
	"os"
	"sort"
	str "strings"
)

var baseString = "https://play.google.com/store/apps/details?id="
var baseYouTubeURL = "https://www.youtube.com/watch?v="

// for my own
const debug = false

var divider = fmt.Sprintf(ansi.Color(str.Repeat("-", 56)+"\n", "yellow"))

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
			getData(c.Args()[0])
			fmt.Println("\n")
			fmt.Println(divider)
		}
	}

	app.Run(os.Args)
}

func getData(pkgName string) {
	// Using with file
	var doc *goquery.Document
	var err error
	if debug {
		f, err := os.Open("yoteshin.html")
		panicIf(err)
		defer f.Close()
		doc, err = goquery.NewDocumentFromReader(f)
	} else {
		doc, err = goquery.NewDocument(fmt.Sprintf("%s%s", baseString, pkgName))
		panicIf(err)
	}

	tmp := make(map[TitleMap]string)

	tmp[TitleMap{1, "Title"}] = doc.Find(`div[itemprop='name']`).First().Text()
	tmp[TitleMap{2, "Category"}] = str.TrimSpace(doc.Find(".category").First().Text())

	price := str.TrimSpace(doc.Find(".price").First().Text())
	if price == "Install" {
		tmp[TitleMap{3, "Price"}] = "Free"
	} else {
		tmp[titleMap{3, "Price"}] = str.TrimSpace(doc.Find(".price").First().Text())
	}

	doc.Find(".meta-info").Each(func(i int, s *goquery.Selection) {
		fieldName := str.TrimSpace(s.Find(".title").Text())
		switch fieldName {
		case "Updated":
			tmp[titleMap{4, "Updated"}] = s.Find(".content").Text()
		case "Installs":
			tmp[titleMap{5, "Total Installs"}] = s.Find(".content").Text()
		case "Size":
			tmp[titleMap{6, "Size"}] = s.Find(".content").Text()
		case "Current Version":
			tmp[titleMap{7, "Current Version"}] = s.Find(".content").Text()
		case "Requires Android":
			tmp[titleMap{8, "Requires Android"}] = s.Find(".content").Text()
		case "Content Rating":
			tmp[titleMap{9, "Content Rating"}] = s.Find(".content").Text()
		case "Developer":
			// Ugly hack
			s.Find(".dev-link").Each(func(i int, t *goquery.Selection) {
				nodeHref, _ := t.Attr("href")
				if str.Contains(nodeHref, "mailto:") {
					tmp[titleMap{10, "Email"}] = str.Split(nodeHref, "mailto:")[1]
				} else {
					raw := str.Split(nodeHref, "&")[0]
					tmp[titleMap{11, "Website"}] = str.Split(raw, "q=")[1]
				}
			})
		}
	})

	score := doc.Find(".score-container").First()
	if score != nil {
		tmp[titleMap{12, "Score"}] = score.Find(".score").First().Text()
		node := doc.Find(`meta[itemprop='ratingCount']`)
		v, _ := node.Attr("content")
		tmp[titleMap{13, "Votes"}] = v
	}

	tmp[titleMap{14, "Developer"}] = str.TrimSpace((doc.Find(`div[itemprop='author']`).Find(".primary").Text()))
	tmp[titleMap{15, "What's New"}] = doc.Find(".whatsnew .recent-change").Text()

	tmp[titleMap{16, "Description"}] = doc.Find(`div[itemprop='description']`).First().Text()

	tmp[titleMap{17, "App Id"}] = pkgName
	tmp[titleMap{18, "Icon Url"}], _ = doc.Find(".cover-image").Attr("src")

	// Should we make slice or just an string ?
	var imgSlice []string
	doc.Find(".full-screenshot").Each(func(i int, s *goquery.Selection) {
		fsLinks, _ := s.Attr("src")
		imgSlice = append(imgSlice, fsLinks)
	})

	for _, imgSliceLinks := range imgSlice {
		tmp[titleMap{19, "Full Screenshot"}] += fmt.Sprintf("%s\n", imgSliceLinks)
	}

	tmp[titleMap{20, "Market Url"}] = marketURL(pkgName)

	doc.Find(".recommendation").Find(".rec-cluster").Each(func(i int, recommended *goquery.Selection) {
		header := str.TrimSpace(recommended.Find(".heading").First().Text())
		if header == "Similar" {
			recommended.Find(".card").Each(func(j int, card *goquery.Selection) {
				similarAppIds, _ := card.Attr("data-docid")
				tmp[titleMap{21, "Related App"}] += fmt.Sprintf("%s\n", similarAppIds)
			})
		} else {
			recommended.Find(".card").Each(func(j int, card *goquery.Selection) {
				similarAppIds, _ := card.Attr("data-docid")
				tmp[titleMap{22, "More from Developer"}] += fmt.Sprintf("%s\n", similarAppIds)
			})
		}
	})

	doc.Find(".play-action-container").Each(func(i int, node *goquery.Selection) {
		url, _ := node.Attr("data-video-url")
		videoID := str.Split(str.Split(url, "embed/")[1], "?")[0]
		if len(url) > 0 {
			tmp[titleMap{23, "YouTube Url"}] = fmt.Sprintf("%s%s", baseYouTubeURL, videoID)
		}
	})

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
		rows = fmt.Sprintf("%s %s | %s\n", k.Title, buffer(k.Title), str.TrimSpace(tmp[k]))
		fmt.Printf(rows)
	}
}

func marketURL(pkgName string) string {
	return fmt.Sprintf("%s%s&hl=en", baseString, pkgName)
}

// I copied it from https://github.com/addyosmani/psi/blob/master/lib%2Futils.js#L36-L50
func buffer(msg string) string {
	var ret = ""
	length := 24
	length = length - len(msg) - 1

	if length > 0 {
		ret = str.Repeat(" ", length)
	}

	return ret
}

type byIndex []titleMap

func (a byIndex) Len() int           { return len(a) }
func (a byIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type titleMap struct {
	Index int
	Title string
}
