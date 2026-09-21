package wordbanks

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed data/*.txt
var wordbankFS embed.FS

func CreateWordSlice(selectedWordbank string) []string {
	dat, err := wordbankFS.ReadFile(selectedWordbank)
	if err != nil {
		panic(err)
	}
	datString := strings.TrimSpace(string(dat))
	stringSlice := strings.Fields(datString)
	return stringSlice
}

func GetAllFilenames() (files []string, err error) {
	if err := fs.WalkDir(wordbankFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}
