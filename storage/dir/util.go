package dir

import (
	"bookarr/storage"
	"errors"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

func fileShouldBeIgnored(filename string) bool {
	const (
		ignoreFile       = true
		includeFile      = false
		hiddenFilePrefix = "."
	)

	if strings.HasPrefix(filename, hiddenFilePrefix) {
		return ignoreFile
	}

	ext := filepath.Ext(filename)
	if _, ok := supportedBookExtensions[ext]; ok {
		return includeFile
	}
	if _, ok := supportedImageExtensions[ext]; ok {
		return includeFile
	}

	return ignoreFile
}

func getRel(filename string, pathType storage.PathType) string {
	if pathType == storage.PathTypeAquisition || pathType == storage.PathTypeNavigation {
		return "subsection"
	}

	ext := filepath.Ext(filename)
	if _, ok := supportedImageExtensions[ext]; ok {
		return "http://opds-spec.org/image/thumbnail"
	}

	// mobi, epub, etc
	return "http://opds-spec.org/acquisition"
}

func getMimeType(name string, pathType storage.PathType) string {
	switch pathType {
	case storage.PathTypeFile:
		return mime.TypeByExtension(filepath.Ext(name))
	case storage.PathTypeAquisition:
		return "application/atom+xml;profile=opds-catalog;kind=acquisition"
	case storage.PathTypeNavigation:
		return "application/atom+xml;profile=opds-catalog;kind=navigation"
	default:
		return mime.TypeByExtension("xml")
	}
}
func getPathType(dirpath string) storage.PathType {
	fi, err := os.Stat(dirpath)
	if err != nil {
		log.Printf("getPathType os.Stat err: %s", err)
		return storage.PathTypeNotExists
	}

	if !fi.IsDir() {
		return storage.PathTypeFile
	}

	dirEntries, err := os.ReadDir(dirpath)
	if err != nil {
		log.Printf("getPathType: readDir err: %s", err)
		return storage.PathTypeNotExists
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			return storage.PathTypeNavigation
		}
	}

	// Directory of directories
	return storage.PathTypeAquisition
}

// verify path use a trustedRoot to avoid http transversal
// from https://www.stackhawk.com/blog/golang-path-traversal-guide-examples-and-prevention/
func sanitisePath(path, trustedRoot string) (string, error) {
	// clean is already used upstream but leaving this
	// to keep the functionality of the function as close as possible to the blog.
	c := filepath.Clean(path)

	// get the canonical path
	r, err := filepath.EvalSymlinks(c)
	if err != nil {
		fmt.Println("Error " + err.Error())
		return c, errors.New("unsafe or invalid path specified")
	}

	if !inTrustedRoot(r, trustedRoot) {
		return r, errors.New("unsafe or invalid path specified")
	}

	return r, nil
}

func inTrustedRoot(path string, trustedRoot string) bool {
	return strings.HasPrefix(path, trustedRoot)
}

// absoluteCanonicalPath returns the canonical path of the absolute path that was passed
func absoluteCanonicalPath(aPath string) (string, error) {
	// get absolute path
	aPath, err := filepath.Abs(aPath)
	if err != nil {
		return "/doesNotExist", fmt.Errorf("get absolute path %s: %w", aPath, err)
	}

	// get canonical path
	aPath, err = filepath.EvalSymlinks(aPath)
	if err != nil {
		return "/doesNotExist", fmt.Errorf("get connonical path from absolute path %s: %w", aPath, err)
	}

	return aPath, nil
}
