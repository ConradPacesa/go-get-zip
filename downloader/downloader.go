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

// Repo is a github repo
type Repo struct {
	RepoName string
	URL      string
	Version  string
	Source   string
}

// GetImports gets all of the Imports for a Github repo
func GetImports(zipFilepath, githubRepo string) ([]string, error) {
	parseGithubURL := ParseGithubURL(githubRepo)

	githubRepoSlice := strings.Split(parseGithubURL.URL, "/")

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
	parsedGHRepo := ParseGithubURL(githubRepo)

	splitURL := strings.Split(parsedGHRepo.URL, "/")

	// Create the zip folder
	zipFilepath := fmt.Sprintf("./tmp/%v.zip", splitURL[len(splitURL)-1])

	out, err := os.Create(zipFilepath)
	if err != nil {
		return "", err
	}

	url := createAPIURL(parsedGHRepo)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		url = createAltURL(parsedGHRepo)
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

// ParseGithubURL takes a gopkg url and converts it to a github url
func ParseGithubURL(githubrepo string) Repo {
	githubRepoStruct := Repo{}

	var version string
	var repoName string
	var gitHubURL []string
	var gitHubURLString string

	if strings.Contains(githubrepo, "gopkg") {
		repoSlice := strings.Split(githubrepo, "/")

		gitHubURL = append(gitHubURL, "github.com")

		if len(repoSlice) == 3 {
			gitHubURL = append(gitHubURL, repoSlice[1])

			v := strings.Split(repoSlice[2], ".")

			gitHubURL = append(gitHubURL, v[0])
			gitHubURLString = strings.Join(gitHubURL, "/")
			repoName = v[0]
			version = v[1]

		} else {

			var usr []string
			usr = append(usr, "go-")

			v := strings.Split(repoSlice[1], ".")

			usr = append(usr, v[0])

			usrString := strings.Join(usr, "")

			gitHubURL = append(gitHubURL, usrString)
			gitHubURL = append(gitHubURL, v[0])
			gitHubURLString = strings.Join(gitHubURL, "/")
			repoName = v[0]
			version = v[1]

		}

		githubRepoStruct = Repo{
			RepoName: repoName,
			URL:      gitHubURLString,
			Version:  version,
			Source:   "gopkg",
		}

	} else {

		githubRepoStruct = Repo{
			RepoName: repoName,
			URL:      githubrepo,
			Version:  version,
			Source:   "github",
		}

	}

	return githubRepoStruct
}

func createAPIURL(githubRepo Repo) string {
	splitURL := strings.Split(githubRepo.URL, "/")

	var baseURL bytes.Buffer
	baseURL.WriteString("https://api.github.com/repos")

	for _, u := range splitURL[1:] {
		baseURL.WriteString("/")
		baseURL.WriteString(u)
	}

	baseURL.WriteString("/zipball/")
	baseURL.WriteString(githubRepo.Version)

	fmt.Println(baseURL.String())

	return baseURL.String()
}

func createAltURL(githubRepo Repo) string {
	splitURL := strings.Split(githubRepo.URL, "/")

	var baseURL bytes.Buffer
	baseURL.WriteString("https://github.com")

	for _, u := range splitURL[1:] {
		baseURL.WriteString("/")
		baseURL.WriteString(u)
	}

	baseURL.WriteString("/archive/")
	baseURL.WriteString(githubRepo.Version)
	baseURL.WriteString(".zip")

	fmt.Println(baseURL.String())

	return baseURL.String()
}
