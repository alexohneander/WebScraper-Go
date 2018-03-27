package main

import (
	"fmt"
  "regexp"
  "os"
  "encoding/json"
  //"encoding/csv"
  "io"
  "log"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var path = "result.txt"
var data string

func main() {
  // define Global vars
  var counter int
  var domain string
  var search string

  // Create Result File
  deleteFile()
  createFile()

  // Course stores information about a coursera course
  type Site struct {
    Title       string
    URL         string
  }

  // Get Domain
  fmt.Print("Domain eingeben (bollants.de): ")
  fmt.Scanln(&domain)
  fmt.Print(domain)

  // Get Search request
  fmt.Print("Was wird gesucht (Sterne): ")
  fmt.Scanln(&search)
  fmt.Print(search)

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains(domain, "www." + domain),

    // Cache responses to prevent multiple download of pages
    // even if the collector is restarted
    colly.CacheDir("./cache"),
	)

	extensions.RandomUserAgent(c)

  sites := make([]Site, 0, 4000)


	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})



  c.OnHTML("div[id=content]", func(e *colly.HTMLElement){
      title := e.ChildText("h1")
      description := e.ChildText("h2")
      content := e.ChildText("p")
      //span := e.ChildText("span")

      //match1, _ :=
      if regexp.MustCompile(search).MatchString(title) == true{
        fmt.Printf("Content Found: %s %s\n",  title, description)
        site := Site{
    			Title:       title,
    			URL:         e.Request.URL.String(),
    		}
        writeFile(e.Request.URL.String())
        sites = append(sites, site)
      }

      //match2, _ := regexp.MustCompile(`STERNE`).MatchString(description)
      if regexp.MustCompile(search).MatchString(description) == true{
        fmt.Printf("Content Found: %s %s\n",  title, description)
        site := Site{
          Title:       title,
          URL:         e.Request.URL.String(),
        }
        writeFile(e.Request.URL.String())
        sites = append(sites, site)
      }

      //match3, _ := regexp.MustCompile(`STERNE`).MatchString(content)
      if regexp.MustCompile(search).MatchString(content) == true{
        fmt.Printf("Content Found: %s %s\n",  title, description)
        site := Site{
          Title:       title,
          URL:         e.Request.URL.String(),
        }

        writeFile(e.Request.URL.String())
        sites = append(sites, site)
      }

  })

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL.String())
     counter ++;
     fmt.Println("counter:", counter)
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www." + domain)

  enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(sites)

}


func createFile() {
  // detect if file exists
  var _, err = os.Stat(path)

  // create file if not exists
  if os.IsNotExist(err) {
    var file, err = os.Create(path)
    if isError(err) { return }
    defer file.Close()
  }

  fmt.Println("==> done creating file", path)
}



func writeFile(data string) {

  // open file using READ & WRITE permission
  var file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
  if isError(err) { return }
  defer file.Close()

  readFile() 

  // write some text line-by-line to file
  _, err = file.WriteString(data + "\n")
  if isError(err) { return }

  // save changes
  err = file.Sync()
  if isError(err) { return }

  fmt.Println("==> done writing to file")
}



func readFile() {
  // re-open file
  var file, err = os.OpenFile(path, os.O_RDWR, 0644)
  if isError(err) { return }
  defer file.Close()

  // read file, line by line
  var text = make([]byte, 1024)
  for {
    _, err = file.Read(text)
    
    // break if finally arrived at end of file
    if err == io.EOF {
      break
    }
    
    // break if error occured
    if err != nil && err != io.EOF {
      isError(err)
      break
    }
  }
  
  fmt.Println("==> done reading from file")
  fmt.Println(string(text))

}

func deleteFile() {
  // delete file
  var err = os.Remove(path)
  if isError(err) { return }

  fmt.Println("==> done deleting file")
}

func checkError(message string, err error) {
    if err != nil {
        log.Fatal(message, err)
    }
}

func isError(err error) bool {
  if err != nil {
    fmt.Println(err.Error())
  }

  return (err != nil)
}