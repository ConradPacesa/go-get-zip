package copier

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ConradPacesa/go-get-zip/downloader"
)

// CopyToGopath unzips the zip file from the tmp folder into the appropriate
// folder in the $GOPATH
func CopyToGopath(src, githubURL string) ([]string, error) {
	gopath := os.Getenv("GOPATH")
	var dest string

	parsedGithubURL := downloader.ParseGithubURL(githubURL)

	if parsedGithubURL.Source == "github" {
		dest = fmt.Sprintf("%v/src/%v", gopath, parsedGithubURL.URL)
	} else {
		dest = fmt.Sprintf("%v/src/gopkg.in/%v.%v", gopath, parsedGithubURL.RepoName, parsedGithubURL.Version)
	}

	fmt.Println(dest)

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		cleanName := strings.Split(f.Name, "/")[1:]
		fname := strings.Join(cleanName, "/")
		fpath := filepath.Join(dest, fname)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {

			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			outFile.Close()

			if err != nil {
				return filenames, err
			}
		}
	}

	return filenames, nil
}
