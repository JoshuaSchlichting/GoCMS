package landing_page

type NavBarLink struct {
	URL  string
	Text string
}

type FeaturedItem struct {
	Title    string
	Subtitle string
	Body     string
	ImageURL string
}

type LandingPageModel struct {
	Title         string
	Brand         string
	SignInURL     string
	Body          string
	Heading       string
	Subheading    string
	NavBarLinks   []NavBarLink
	FeaturedItems []FeaturedItem
}
