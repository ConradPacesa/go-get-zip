package downloader

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

// GetImports gets all of the Imports for a Github repo
func GetImports(zipFilepath, githubRepo string) ([]string, error) {
	githubRepoSlice := strings.Split(githubRepo, "/")

	cmd := exec.Command(`go`, `list`, `-f`, `{{ join .Deps "\n" }}`)
	cmd.Dir = zipFilepath

	stdout, err := cmd.CombinedOutput()
	outString := string(stdout)
	depList := strings.Split(outString, "\n")

	var importList []string

	for _, d := range depList {
		if strings.HasPrefix(d, "github.com") {
			fmt.Println(strings.Split(d, "/")[:2])
			fmt.Println(githubRepoSlice[:2])
		}
		if strings.HasPrefix(d, "github.com") && !reflect.DeepEqual(strings.Split(d, "/")[:3], githubRepoSlice[:3]) {
			importList = append(importList, d)
		}
	}

	return importList, err
}

// DownloadZip downloads the GitHub Repo as a zip file into the tmp directory
func DownloadZip(githubRepo string) (string, error) {
	splitURL := strings.Split(githubRepo, "/")

	// Create the zip folder
	zipFilepath := fmt.Sprintf("./tmp/%v.zip", splitURL[len(splitURL)-1])
	out, err := os.Create(zipFilepath)
	if err != nil {
		return "", err
	}

	url := createApiUrl(githubRepo)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		url = createAltUrl(githubRepo)
		resp, err = http.Get(url)
		if err != nil {
			return "", err
		}
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	out.Close()

	return zipFilepath, nil
}

func createApiUrl(githubRepo string) string {
	splitURL := strings.Split(githubRepo, "/")

	var baseURL bytes.Buffer
	baseURL.WriteString("https://api.github.com/repos")

	for _, u := range splitURL[1:] {
		baseURL.WriteString("/")
		baseURL.WriteString(u)
	}

	baseURL.WriteString("/zipball/master")

	fmt.Println(baseURL.String())

	return baseURL.String()
}

func createAltUrl(githubRepo string) string {
	splitURL := strings.Split(githubRepo, "/")

	var baseURL bytes.Buffer
	baseURL.WriteString("https://github.com")

	for _, u := range splitURL[1:] {
		baseURL.WriteString("/")
		baseURL.WriteString(u)
	}

	baseURL.WriteString("/archive/v1.zip")

	fmt.Println(baseURL.String())

	return baseURL.String()
}
