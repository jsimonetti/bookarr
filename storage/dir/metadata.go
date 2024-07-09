package dir

import (
	"bookarr/storage"
	"bookarr/storage/epub"
	"errors"
	"image"
	"log"
)

var errNotRecognised = errors.New("not recognised")

func addMetadata(e *storage.Entry, filename string) error {
	switch e.Type {
	case "application/epub+zip":
		return addEpubMetadata(e, filename)
	}

	e.Metadata = &storage.NOOPMetadata{}
	return errNotRecognised
}

func addEpubMetadata(e *storage.Entry, filename string) error {
	e.Metadata = getEpubMetadata(filename)
	if e.Metadata == nil {
		return errNotRecognised
	}
	return nil
}

func getEpubMetadata(filename string) storage.Metadata {
	metadata, err := epub.OpenReader(filename)
	if err != nil {
		return nil
	}
	defer metadata.Close()
	return metadata.Rootfiles[0]
}

func getEpubCover(filename string) *image.Image {
	epub, err := epub.OpenReader(filename)
	if err != nil {
		log.Printf("File OpenReader err: %s", err)
		return nil
	}
	defer epub.Close()
	p := epub.Rootfiles[0]

	cover := ""
	for _, meta := range p.Meta {
		if meta.Name == "cover" {
			cover = meta.Content
			break
		}
	}
	if cover != "" {
		for _, item := range p.Manifest.Items {
			if item.ID == cover {
				f, err := item.Open()
				if err != nil {
					log.Printf("File item.Open err: %s", err)
					return nil
				}

				img, _, err := image.Decode(f)
				if err != nil {
					log.Printf("imageDecode err: %s", err)
					return nil
				}
				return &img
			}
		}
	}
	return nil
}

var supportedBookExtensions = map[string]struct{}{
	".mobi": {},
	".epub": {},
	".pdf":  {},
	".cbz":  {},
	".cbr":  {},
	".fb2":  {},
}

var supportedImageExtensions = map[string]struct{}{
	".png":  {},
	".jpg":  {},
	".jpeg": {},
	".gif":  {},
}
