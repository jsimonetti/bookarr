package opdsv1

import "encoding/xml"

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	// Xmlns     string   `xml:"xmlns,attr,omitempty"`
	Title     string   `xml:"title"`
	ID        string   `xml:"id"`
	Link      []Link   `xml:"link"`
	Published string   `xml:"published,omitempty"`
	Updated   TimeStr  `xml:"updated"`
	Category  string   `xml:"category,omitempty"`
	Authors   []Author `xml:"author"`
	Summary   *Summary `xml:"summary"`
	Content   *Content `xml:"content"`
	Rights    string   `xml:"rights,omitempty"`
	Source    string   `xml:"source,omitempty"`

	// Extensions
	Language string `xml:"dc:language,omitempty"`
}

func NewEntry(title, id string, updated TimeStr) *Entry {
	return &Entry{
		Title:   title,
		ID:      id,
		Updated: updated,
	}
}

func (f *Feed) AddEntry(e *Entry) {
	f.Entry = append(f.Entry, e)
}
