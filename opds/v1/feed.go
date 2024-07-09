package opdsv1

import (
	"encoding/xml"
	"net/url"
	"time"
)

const (

	// Feed link types
	//FeedAcquisitionLinkType = "application/atom+xml;profile=opds-catalog;kind=acquisition"
	FeedNavigationLinkType = "application/atom+xml;profile=opds-catalog;kind=navigation"
	//FeedSearchLinkType      = "application/opensearchdescription+xml"
)

type Feed struct {
	XMLName      xml.Name `xml:"feed"`
	Xmlns        string   `xml:"xmlns,attr"`
	XmlnsDC      string   `xml:"xmlns:dc,attr,omitempty"`
	XmlnsOS      string   `xml:"xmlns:opensearch,attr,omitempty"`
	XmlnsOPDS    string   `xml:"xmlns:opds,attr,omitempty"`
	Title        string   `xml:"title"`
	ID           string   `xml:"id"`
	Updated      TimeStr  `xml:"updated"`
	Link         []Link   `xml:"link"`
	Author       []Author `xml:"author,omitempty"`
	Entry        []*Entry `xml:"entry"`
	Category     string   `xml:"category,omitempty"`
	Icon         string   `xml:"icon,omitempty"`
	Logo         string   `xml:"logo,omitempty"`
	Content      string   `xml:"content,omitempty"`
	Subtitle     string   `xml:"subtitle,omitempty"`
	SearchResult uint     `xml:"opensearch:totalResults,omitempty"`
}

type Link struct {
	XMLName xml.Name `xml:"link"`
	Type    string   `xml:"type,attr,omitempty"`
	Title   string   `xml:"title,attr,omitempty"`
	Href    string   `xml:"href,attr,omitempty"`
	Rel     string   `xml:"rel,attr,omitempty"`
	Length  string   `xml:"length,attr,omitempty"`
}

type Author struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name,omitempty"`
	Uri     string   `xml:"uri,omitempty"`
}

type Summary struct {
	XMLName xml.Name `xml:"summary"`
	Content string   `xml:",chardata"`
	Type    string   `xml:"type,attr"`
}

type Content struct {
	XMLName xml.Name `xml:"content"`
	Content string   `xml:",chardata"`
	Type    string   `xml:"type,attr"`
}

type TimeStr string

func (f *Feed) Time(t time.Time) TimeStr {
	return TimeStr(t.Format("2006-01-02T15:04:05-07:00"))
}

func NewFeed(title, subtitle string, self, baseUrl *url.URL) *Feed {
	if baseUrl == nil {
		return nil
	}
	if self == nil {
		return nil
	}

	f := &Feed{
		Xmlns:     "http://www.w3.org/2005/Atom",
		XmlnsDC:   "http://purl.org/dc/terms/",
		XmlnsOS:   "http://a9.com/-/spec/opensearch/1.1/",
		XmlnsOPDS: "http://opds-spec.org/2010/catalog",
		Title:     title,
		ID:        self.String(),
		Link: []Link{
			//{Rel: "search", Href: baseUrl.JoinPath(self.String(), "search?q={searchTerms}").String(), Type: FeedSearchLinkType, Title: "Search on catalog"},
			{Rel: "start", Href: baseUrl.String() + "/", Type: FeedNavigationLinkType},
			{Rel: "self", Href: baseUrl.JoinPath(self.String()).String(), Type: FeedNavigationLinkType},
		},
		Subtitle: subtitle,
	}
	f.Updated = f.Time(time.Now())
	return f
}
