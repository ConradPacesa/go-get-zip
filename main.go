package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ConradPacesa/go-get-zip/downloader"
	"github.com/ConradPacesa/go-get-zip/installer"
)

func main() {
	os.Getenv("http_proxy")
	os.Getenv("https_proxy")
	gopath := os.Getenv("GOPATH")

	if len(os.Args) == 1 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		p := fmt.Sprintf("%v/src/", gopath)
		gitHubRepo, err := filepath.Rel(p, dir)

		if runtime.GOOS == "windows" {
			wp := strings.Split(gitHubRepo, "\\")
			gitHubRepo = strings.Join(wp, "/")
		}

		imports, err := downloader.GetImports(dir, gitHubRepo)
		if err != nil {
			log.Fatal(err)
		}

		for _, im := range imports {
			_, err := installer.Install(im)
			if err != nil {
				log.Fatal(err)
			}
		}

	} else if len(os.Args) == 2 {
		githubRepo := os.Args[1]

		f, err := installer.Install(githubRepo)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("Installed %s\n", f)
		}

	} else {
		log.Fatalf("There were too many arguments passed in. Pass in either 1 or zero args.")
	}

	fmt.Println("All done!")
}
