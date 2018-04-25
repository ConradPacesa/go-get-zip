package installer

import (
	"fmt"
	"os/exec"
)

// Install installs the source code from the downloaded GitHub repo
func Install(filepath string) {
	cmd := exec.Command("go", "install", filepath)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("There was an error installing the file: %v", err)
	}
	fmt.Println(string(stdout))
}
