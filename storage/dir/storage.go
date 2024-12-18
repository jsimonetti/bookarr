package dir

import (
	"bookarr/storage"
	"bufio"
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"
	"sort"
)

type fileStore struct {
	rootDir string
}

func NewFileStore(rootDir string) storage.Store {
	rootDir, err := absoluteCanonicalPath(rootDir)
	if err != nil {
		log.Fatal(err)
	}
	return &fileStore{rootDir: rootDir}
}

func (fs *fileStore) File(path string) *storage.File {
	fPath := filepath.Join(fs.rootDir, path)

	safePath, err := verifyPath(fPath, fs.rootDir)
	if err != nil {
		log.Printf("File verifyPath err: %s", err)
		return nil
	}

	f, err := os.Open(safePath)
	if err != nil {
		log.Printf("File os.Open err: %s", err)
		return nil
	}

	stat, err := f.Stat()
	if err != nil {
		log.Printf("File f.Stat() err: %s", err)
		return nil
	}

	return &storage.File{
		Reader:        f,
		ContentType:   mime.TypeByExtension(filepath.Ext(safePath)),
		ContentLength: stat.Size(),
	}
}

func (fs *fileStore) List(path string) ([]storage.Entry, error) {
	fPath := filepath.Join(fs.rootDir, path)

	safePath, err := verifyPath(fPath, fs.rootDir)
	if err != nil {
		log.Printf("File verifyPath err: %s", err)
		return nil, err
	}

	entries := []storage.Entry{}

	dirEntries, _ := os.ReadDir(safePath)
	for _, entry := range dirEntries {
		if !entry.IsDir() && fileShouldBeIgnored(entry.Name()) {
			continue
		}
		filename := filepath.Join(safePath, entry.Name())
		pathType := getPathType(filename)
		info, err := entry.Info()
		if err != nil {
			continue
		}
		e := storage.Entry{
			Name:       entry.Name(),
			Type:       getMimeType(entry.Name(), pathType),
			Aquisition: getRel(entry.Name(), pathType),
			Updated:    info.ModTime(),
		}

		if err := addMetadata(&e, filename); err != nil {
			if err != errNotRecognised {
				log.Printf("addMetadata err: %s", err)
				continue
			}
		}

		entries = append(entries, e)
	}

	sortEntries(&entries)
	return entries, nil
}

func sortEntries(entries *[]storage.Entry) {
	sortByName(entries)
	sortByAuthor(entries)
}

func sortByName(entries *[]storage.Entry) {
	sortEntriesBy(entries, func(a, b storage.Entry) bool {
		return a.Name < b.Name
	})
}

func sortByAuthor(entries *[]storage.Entry) {
	sortEntriesBy(entries, func(a, b storage.Entry) bool {
		as := a.Metadata.GetCreator()
		bs := b.Metadata.GetCreator()
		return as < bs
	})
}

func sortEntriesBy(entries *[]storage.Entry, less func(a, b storage.Entry) bool) {
	sorter := &entrySorter{
		entries: *entries,
		less:    less,
	}
	sort.Sort(sorter)
}

type entrySorter struct {
	entries []storage.Entry
	less    func(a, b storage.Entry) bool
}

func (s *entrySorter) Len() int {
	return len(s.entries)
}

func (s *entrySorter) Swap(i, j int) {
	s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
}

func (s *entrySorter) Less(i, j int) bool {
	return s.less(s.entries[i], s.entries[j])
}

func (fs *fileStore) Cover(path string) *storage.File {
	fPath := filepath.Join(fs.rootDir, path)
	ext := filepath.Ext(fPath)
	if _, ok := supportedBookExtensions[ext]; !ok {
		return nil
	}

	safePath, err := verifyPath(fPath, fs.rootDir)
	if err != nil {
		log.Printf("File verifyPath err: %s", err)
		return nil
	}

	var cover *image.Image
	switch ext {
	case ".epub":
		cover = getEpubCover(safePath)
	}

	if cover == nil {
		return nil
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	if err := jpeg.Encode(w, *cover, nil); err != nil {
		return nil
	}

	r := bufio.NewReader(&b)
	rc := io.NopCloser(r)
	return &storage.File{
		Reader:        rc,
		ContentType:   mime.TypeByExtension(".jpeg"), // "image/jpeg
		ContentLength: int64(b.Len()),
	}
}

func (fs *fileStore) Thumbnail(path string) *storage.File {
	return nil
}

func (fs *fileStore) PathType(path string) storage.PathType {
	fPath := filepath.Join(fs.rootDir, path)

	safePath, err := verifyPath(fPath, fs.rootDir)
	if err != nil {
		log.Printf("File verifyPath err: %s", err)
		return storage.PathTypeNotExists
	}

	return getPathType(safePath)
}
