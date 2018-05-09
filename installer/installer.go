package installer

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ConradPacesa/go-get-zip/copier"
	"github.com/ConradPacesa/go-get-zip/downloader"
)

// Install recursively finds all the dependencies and installs them
func Install(githubRepo string) (string, error) {

	fmt.Println("Downloading Zip...")

	// Download the zipfile into the tmp folder
	zipFilepath, err := downloader.DownloadZip(githubRepo)
	if err != nil {
		fmt.Printf("There was an error downloading the file %v", err)
	}

	fmt.Println("Unzipping file into $GOPATH")

	// Unzip the contents from the zipfile into the $GOPATH
	file, err := copier.CopyToGopath(zipFilepath, githubRepo)
	if err != nil {
		fmt.Printf("There was an error copying the files over %v", err)
	}

	// Get the root of unzipped filepath
	fp := file[0]

	// List all the dependencies of the downloaded package
	imports, err := downloader.GetImports(fp, githubRepo)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(imports)

	for _, im := range imports {
		_, err := Install(im)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Installing Go packages...")

	// Install the unzipped files
	f, err := InstallPackage(githubRepo)
	if err != nil {
		fmt.Printf("There was an error installing the package %s: %s", githubRepo, err)
	} else {
		fmt.Println(f)
	}

	// Close the outfile and delete it
	err = os.Remove(zipFilepath)
	if err != nil {
		fmt.Println(err)
	}

	return githubRepo, nil
}

// InstallPackage installs the source code from the downloaded GitHub repo
func InstallPackage(filepath string) (string, error) {
	cmd := exec.Command("go", "install", filepath)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}
