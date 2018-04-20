package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func downloadZip(githubRepo string) string {
	splitURL := strings.Split(githubRepo, "/")

	// Create the zip folder
	zipFilepath := fmt.Sprintf("./tmp/%v.zip", strings.Join(splitURL, "-"))
	out, err := os.Create(zipFilepath)
	if err != nil {
		fmt.Println(err)
	}

	var baseURL bytes.Buffer
	baseURL.WriteString("https://api.github.com/repos")

	for _, u := range splitURL[1:] {
		baseURL.WriteString("/")
		baseURL.WriteString(u)
	}

	baseURL.WriteString("/zipball/master")

	resp, err := http.Get(baseURL.String())
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	out.Close()

	return zipFilepath
}

func copyToGopath(src string, githubURL string) ([]string, error) {
	gopath := os.Getenv("GOPATH")

	dest := fmt.Sprintf("%v/src/%v", gopath, githubURL)

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

		fpath := filepath.Join(dest, f.Name)
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

func main() {
	os.Getenv("http_proxy")
	os.Getenv("https_proxy")

	githubRepo := os.Args[1]

	zipFilepath := downloadZip(githubRepo)

	filenames, err := copyToGopath(zipFilepath, githubRepo)
	if err != nil {
		fmt.Printf("There was an error copying the files over %v", err)
	}

	fmt.Printf("The following files were copied: %v", filenames)
	// Close the outfile and delete it

	// err = os.Remove(zipFilepath)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
