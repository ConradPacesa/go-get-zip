package main

import (
	"fmt"
	"os"

	"github.com/ConradPacesa/go-get-zip/installer"
)

func main() {
	os.Getenv("http_proxy")
	os.Getenv("https_proxy")

	githubRepo := os.Args[1]

	f, err := installer.Install(githubRepo)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Installed %s\n", f)
	}

	fmt.Println("All done!")
}
