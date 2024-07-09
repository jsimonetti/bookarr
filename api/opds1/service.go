// package service provides a http handler that reads the path in the request.url and returns
// an xml document that follows the OPDS 1.1 standard
// https://specs.opds.io/opds-1.1.html
package opds1

import (
	"encoding/xml"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	opdsv1 "bookarr/opds/v1"
	"bookarr/storage"

	"github.com/gin-gonic/gin"
)

func init() {
	_ = mime.AddExtensionType(".mobi", "application/x-mobipocket-ebook")
	_ = mime.AddExtensionType(".epub", "application/epub+zip")
	_ = mime.AddExtensionType(".cbz", "application/x-cbz")
	_ = mime.AddExtensionType(".cbr", "application/x-cbr")
	_ = mime.AddExtensionType(".fb2", "text/fb2+xml")
	_ = mime.AddExtensionType(".pdf", "application/pdf")
}

const xmlHeader = `<?xml version="1.0" encoding="utf-8"?>`

type opdsv1Handler struct {
	baseURL string
	storage storage.Store
}

func New(baseUrl string, store storage.Store) *opdsv1Handler {
	return &opdsv1Handler{
		baseURL: baseUrl,
		storage: store,
	}
}

// Handler serve the content of a book file or
// returns an Acquisition Feed when the entries are documents or
// returns an Navegation Feed when the entries are other folders
func (h *opdsv1Handler) Handler(c *gin.Context) {
	if strings.HasSuffix(c.Request.RequestURI, "/cover") {
		h.GetCover(c)
		return
	}
	if strings.HasSuffix(c.Request.RequestURI, "/thumbnail") {
		h.GetThumbnail(c)
		return
	}

	urlPath := c.Param("path")

	var feed *opdsv1.Feed
	var contentType storage.PathType

	switch h.storage.PathType(urlPath) {
	case storage.PathTypeFile:
		file := h.storage.File(urlPath)
		if file != nil {
			defer file.Reader.Close()
			c.DataFromReader(http.StatusOK, file.ContentLength, file.ContentType, file.Reader, nil)
			return
		}
		c.Writer.Write([]byte(xmlHeader))
		c.XML(http.StatusNotFound, nil)
		return
	case storage.PathTypeAquisition:
		feed = h.makeFeed(c)
		contentType = storage.PathTypeAquisition
	case storage.PathTypeNavigation:
		feed = h.makeFeed(c)
		contentType = storage.PathTypeNavigation
	case storage.PathTypeNotExists:
		c.Writer.Write([]byte(xmlHeader))
		c.XML(http.StatusNotFound, nil)
		return
	}

	c.Writer.Write([]byte(xmlHeader))

	content, err := xml.Marshal(feed)
	if err != nil {
		log.Printf("xml.Marshal err: %s", err)
		c.XML(http.StatusInternalServerError, nil)
	}
	c.Data(http.StatusOK, string(contentType), content)

}

func (h *opdsv1Handler) GetCover(c *gin.Context) {
	urlPath := c.Param("path")
	path := strings.TrimSuffix(urlPath, "/cover")

	cover := h.storage.Cover(path)
	if cover != nil {
		defer cover.Reader.Close()
		c.DataFromReader(http.StatusOK, cover.ContentLength, cover.ContentType, cover.Reader, nil)
		return
	}
	c.XML(http.StatusNotFound, nil)
}

func (h *opdsv1Handler) GetThumbnail(c *gin.Context) {
	urlPath := c.Param("path")
	path := strings.TrimSuffix(urlPath, "/thumbnail")

	cover := h.storage.Thumbnail(path)
	if cover != nil {
		defer cover.Reader.Close()
		c.DataFromReader(http.StatusOK, cover.ContentLength, cover.ContentType, cover.Reader, nil)
		return
	}
	c.XML(http.StatusNotFound, nil)
}

func (h opdsv1Handler) makeFeed(c *gin.Context) *opdsv1.Feed {
	urlPath := c.Param("path")

	selfUrl := &url.URL{Path: urlPath}
	baseUrl := &url.URL{Path: h.baseURL}
	feed := opdsv1.NewFeed("Catalog in "+urlPath, "", selfUrl, baseUrl)

	dirEntries, _ := h.storage.List(urlPath)
	for _, entry := range dirEntries {

		originalName := entry.Name
		if entry.Metadata.GetTitle() != "" {
			entry.Name = entry.Metadata.GetTitle()
		}

		e := &opdsv1.Entry{
			Title:   entry.Name,
			ID:      filepath.Join(c.Request.RequestURI, url.PathEscape(originalName)),
			Updated: feed.Time(entry.Updated),
			Link: []opdsv1.Link{
				{
					Type:  entry.Type,
					Title: entry.Name,
					Href:  baseUrl.JoinPath(urlPath, originalName).String(),
					Rel:   entry.Aquisition,
				},
			},
		}
		if entry.Metadata.HasCover() {
			e.Link = append(e.Link, opdsv1.Link{
				Type: "image/jpeg",
				Rel:  "http://opds-spec.org/image",
				Href: baseUrl.JoinPath(urlPath, originalName, "cover").String(),
			})

		}
		if entry.Metadata.HasThumbnail() {
			e.Link = append(e.Link, opdsv1.Link{
				Type: "image/jpeg",
				Rel:  "http://opds-spec.org/image/thumbnail",
				Href: baseUrl.JoinPath(urlPath, originalName, "thumbnail").String(),
			})
		}

		if entry.Metadata.GetCreator() != "" {
			e.Authors = append(e.Authors, opdsv1.Author{
				Name: entry.Metadata.GetCreator(),
			})
		}
		if entry.Metadata.GetSubject() != "" {
			e.Summary = &opdsv1.Summary{
				Content: entry.Metadata.GetSubject(),
				Type:    "text",
			}
		}
		if entry.Metadata.GetDescription() != "" {
			e.Summary = &opdsv1.Summary{
				Content: truncateText(entry.Metadata.GetDescription()) + "...",
				Type:    "text",
			}
		}
		if entry.Metadata.GetLanguage() != "" {
			e.Language = entry.Metadata.GetLanguage()
		}
		feed.AddEntry(e)
	}

	return feed
}

const maxSummaryLength = 450

func truncateText(s string) string {
	if maxSummaryLength > len(s) {
		return s
	}
	return s[:strings.LastIndexAny(s[:maxSummaryLength], " .,:;-")]
}