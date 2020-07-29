package webserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Pages - Pages struct that holds a Page struct.
type Pages struct {
	Path     string
	Template string
	Page     Page
}

// PagesSlice - A type that represents a Pages (struct) slice.
type PagesSlice []Pages

// Page - Page struct.
type Page struct {
	Title string
	Body  []byte
}

func loadpages(pages *PagesSlice) {
	// Loop through all pages and load the body.
	for _, page := range *pages {
		body, _ := ioutil.ReadFile("templates/" + page.Template)

		// Assign body.
		page.Page.Body = body
	}
}

func loadpage(template string) *Page {
	path := "templates/" + template

	body, _ := ioutil.ReadFile(path)

	return &Page{Title: template, Body: body}
}

func addroute(path string, template string, title string, pages *PagesSlice) {
	*pages = append(*pages, Pages{Path: path, Template: template, Page: Page{Title: title}})
}

// ServePage - Serves each request to the HTTP server.
func (pages *PagesSlice) ServePage(w http.ResponseWriter, r *http.Request) {
	// Get current path.
	path := r.URL.Path

	// Initialize page and load not found as default.
	page := loadpage("notfound.html")

	// Loop through each page and see if we have a match based off of path.
	for _, tpage := range *pages {
		if path == tpage.Path {
			page = loadpage(tpage.Template)

			// Since we found a page, let's just assign the title manually.
			page.Title = tpage.Page.Title
		}
	}

	// Write response.
	fmt.Fprintf(w, string(page.Body))
}

// StartServer - Adds necessary routes and start an HTTP server.
func StartServer(ip string, port int) {
	// Setup pages struct
	var pages PagesSlice

	// Add necessary routes.
	addroute("/", "index.html", "Home", &pages)
	addroute("/chatbox", "chatbox.html", "Chat Box", &pages)

	// Load pages.
	loadpages(&pages)

	// Output routes (debug).
	for _, page := range pages {
		fmt.Println("Loading " + page.Path + " -- " + page.Template + " (" + page.Page.Title + ")")
	}

	// Setup web server.
	srv := &http.Server{
		Addr:         ip + ":" + strconv.Itoa(port),
		Handler:      http.HandlerFunc(pages.ServePage),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start listening.
	srv.ListenAndServe()
}
