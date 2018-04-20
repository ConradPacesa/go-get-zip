package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func downloadZip(githubRepo string) (string, error) {
	splitURL := strings.Split(githubRepo, "/")

	// Create the zip folder
	zipFilepath := fmt.Sprintf("./tmp/%v.zip", splitURL[len(splitURL)-1])
	out, err := os.Create(zipFilepath)
	if err != nil {
		return "", err
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
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	out.Close()

	return zipFilepath, nil
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

func install(filepath string) {
	cmd := exec.Command("go", "install", filepath)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("There was an error installing the file: %v", err)
	}
	fmt.Println(stdout)
}

func main() {
	os.Getenv("http_proxy")
	os.Getenv("https_proxy")

	githubRepo := os.Args[1]

	fmt.Println("Downloading Zip...")
	zipFilepath, err := downloadZip(githubRepo)
	if err != nil {
		fmt.Printf("There was an error downloading the file %v", err)
	}

	fmt.Println("Unzipping file into $GOPATH")
	_, err = copyToGopath(zipFilepath, githubRepo)
	if err != nil {
		fmt.Printf("There was an error copying the files over %v", err)
	}

	fmt.Println("Installing Go packages...")

	install(githubRepo)

	// Close the outfile and delete it
	err = os.Remove(zipFilepath)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("All done!")
}
