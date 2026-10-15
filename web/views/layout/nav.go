package layout

type navItem struct {
	Key, Label, Href string
}

// defaultDescription is the summary for pages that don't give their own.
const defaultDescription = "Wildlife spotted along the WRC Safari Rally Kenya route: Naivasha, Hell's Gate and Soysambu."

// navItems are the pages in the header, in order.
var navItems = []navItem{
	{Key: "home", Label: "Home", Href: "/"},
	{Key: "stats", Label: "Stats", Href: "/stats"},
	{Key: "gallery", Label: "Gallery", Href: "/gallery"},
}
