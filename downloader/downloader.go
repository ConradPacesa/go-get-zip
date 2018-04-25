package downloader

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// DownloadZip downloads the GitHub Repo as a zip file into the tmp directory
func DownloadZip(githubRepo string) (string, error) {
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
