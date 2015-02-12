package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/cli"
	"github.com/mgutz/ansi"
	"os"
	. "strings"
)

var baseString = "https://play.google.com/store/apps/details?id="

var divider = fmt.Sprintf(ansi.Color(Repeat("-", 56)+"\n", "grey"))

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
			pkg := c.Args()[0]
			fmt.Printf(ansi.Color("Processing Results for \"%s\"\n", "green"), pkg)
			fmt.Println("\n")
			fmt.Println(Repeat("-", 56))
			fmt.Println("\n")
			GetData(c.Args()[0])
			fmt.Println("\n")
			fmt.Println(Repeat("-", 56))
		}
	}

	app.Run(os.Args)
}

func GetData(pkgName string) {
	// Using with file
	// f, err := os.Open("poweramp.html")
	// PanicIf(err)
	// defer f.Close()
	// doc, err := goquery.NewDocumentFromReader(f)
	doc, err := goquery.NewDocument(fmt.Sprintf("%s%s", baseString, pkgName))
	PanicIf(err)

	tmp := make(map[string]string)
	doc.Find(".meta-info").Each(func(i int, s *goquery.Selection) {
		fieldName := TrimSpace(s.Find(".title").Text())
		switch fieldName {
		case "Updated":
			tmp["updated"] = s.Find(".content").Text()
		case "Installs":
			tmp["Total Installs"] = s.Find(".content").Text()
		case "Size":
			tmp["Size"] = s.Find(".content").Text()
		case "Current Version":
			tmp["Current Version"] = s.Find(".content").Text()
		case "Requires Android":
			tmp["Requires Android"] = s.Find(".content").Text()
		case "Content Rating":
			tmp["Content Rating"] = s.Find(".content").Text()
		case "Developer":
			// Ugly hack
			s.Find(".dev-link").Each(func(i int, t *goquery.Selection) {
				nodeHref, _ := t.Attr("href")
				if Contains(nodeHref, "mailto:") {
					tmp["Email"] = Split(nodeHref, "mailto:")[1]
				} else {
					raw := Split(nodeHref, "&")[0]
					tmp["Website"] = Split(raw, "q=")[1]
				}
			})
		}
	})

	tmp["Category"] = TrimSpace(doc.Find(".category").First().Text())
	tmp["Price"] = TrimSpace(doc.Find(".price").First().Text())
	tmp["Description"] = doc.Find(`div[itemprop='description']`).First().Text()
	tmp["Title"] = doc.Find(`div[itemprop='name']`).First().Text()

	score := doc.Find(".score-container").First()
	if score != nil {
		tmp["Score"] = score.Find(".score").First().Text()
		node := doc.Find(`meta[itemprop='ratingCount']`)
		v, _ := node.Attr("content")
		tmp["Votes"] = v
	}

	tmp["Developer"] = TrimSpace((doc.Find(`div[itemprop='author']`).Find(".primary").Text()))
	tmp["What's New"] = doc.Find(".whatsnew .recent-change").Text()

	for x, y := range tmp {
		var rows string
		rows = fmt.Sprintf("%s %s | %s\n", x, buffer(x, 11), TrimSpace(y))
		fmt.Printf(rows)
	}
}

// I copied it from https://github.com/addyosmani/psi/blob/master/lib%2Futils.js#L36-L50
func buffer(msg string, length int) string {
	var ret = ""

	length = 44

	length = length - len(msg) - 1

	if length > 0 {
		ret = Repeat(" ", length)
	}

	return ret
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
