package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func BuildProject(id string) error {

	projectPath := filepath.Join("./output", id)
	cmdInstall := exec.Command("npm", "install")
	cmdInstall.Dir = projectPath
	cmdInstall.Stdout = os.Stdout
	cmdInstall.Stderr = os.Stderr
	fmt.Println("Running npm install .... ")
	if err := cmdInstall.Run(); err != nil {
		return fmt.Errorf("npm install failed : %w", err)
	}
	cmdBuild := exec.Command("npm", "run", "build")
	cmdBuild.Dir = projectPath
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	fmt.Println("Running npm build .... ")
	if err := cmdBuild.Run(); err != nil {
		return fmt.Errorf("npm build failed : %w", err)
	}
	fmt.Println("Build completed successfully for ", id)
	return nil
}
