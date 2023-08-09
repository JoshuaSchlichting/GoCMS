package blog

import (
	"regexp"

	"github.com/joshuaschlichting/gocms/db"
)

type NavBarLink struct {
	URL  string
	Text string
}

type Post struct {
	ID          int
	Title       string
	Body        string
	Author      string
	PublishedAt string
}

type Page struct {
	Title       string
	Brand       string
	SignInURL   string
	Body        string
	Heading     string
	Subheading  string
	NavBarLinks []NavBarLink
	Posts       []db.BlogPost
	SideWidget  SideWidget
}

type SideWidget struct {
	Title string
	Body  string
}

// GetFirstImageURI returns the first image URI from a markdown blog post's body.
func GetFirstImageURI(body string) string {
	// Regular expression pattern for markdown image link with specific extensions
	re := regexp.MustCompile(`!\[.*?\]\((.*?\.(jpg|png|svg|bmp))\)`)
	matches := re.FindStringSubmatch(body)

	// if matches were found
	if len(matches) > 1 {
		// return first image URI
		return matches[1]
	}
	// return empty string if no image URI found
	return ""
}
