package blog

import (
	"html/template"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type NavBarLink struct {
	URL  string
	Text string
}

type Post struct {
	ID                 uuid.UUID
	Title              string
	Subtitle           string
	FeaturedImageUri   string
	Body               string
	Author             string
	PublishedTimestamp time.Time
	UpdatedTimestamp   time.Time
}

type Page struct {
	Title        string
	Brand        string
	SignInURL    string
	Heading      string
	Subheading   string
	FeaturedPost Post
	NavBarLinks  []NavBarLink
	Posts        []Post
	Body         template.HTML
	SideWidget   SideWidget
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
