package handler

import "fmt"

type Pagebutton struct {
	Class string
	Href  string
	Page  int
}

// newPagebuttons returns the html information needed to create pagination buttons
func newPagebuttons(activePage int, totalPages int) []Pagebutton {
	var pages []Pagebutton

	for i := 1; i <= totalPages; i++ {
		pages = append(pages, Pagebutton{
			Class: getPageButtonClass(i == activePage),
			Href:  fmt.Sprintf("?page=%d", i),
			Page:  i,
		})
	}

	return pages
}

// getPageButtonClass returns the appropriate button html class depending on if its the currently acitve page
func getPageButtonClass(isActivePage bool) string {
	if isActivePage {
		return "active"
	}
	return "waves-effect"
}
