package main

import (
	"fmt"
	"os"

	"github.com/ConradPacesa/go-get-zip/copier"
	"github.com/ConradPacesa/go-get-zip/downloader"
	"github.com/ConradPacesa/go-get-zip/installer"
)

func main() {
	os.Getenv("http_proxy")
	os.Getenv("https_proxy")

	githubRepo := os.Args[1]

	fmt.Println("Downloading Zip...")
	zipFilepath, err := downloader.DownloadZip(githubRepo)
	if err != nil {
		fmt.Printf("There was an error downloading the file %v", err)
	}

	fmt.Println("Unzipping file into $GOPATH")
	_, err = copier.CopyToGopath(zipFilepath, githubRepo)
	if err != nil {
		fmt.Printf("There was an error copying the files over %v", err)
	}

	fmt.Println("Installing Go packages...")

	installer.Install(githubRepo)

	// Close the outfile and delete it
	err = os.Remove(zipFilepath)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("All done!")
}
