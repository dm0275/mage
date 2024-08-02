package magefile

import (
	"fmt"
	"github.com/dm0275/mage/utils"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

var Config = &ProjectConfig{
	OutputDir:  "bin",
	CgoEnabled: false,
	LdFlags:    map[string]string{},
	OsTypes:    []string{"linux"},
	ArchTypes:  []string{"amd64"},
}

type ProjectConfig struct {
	OutputDir  string
	CgoEnabled bool
	LdFlags    map[string]string
	OsTypes    []string
	ArchTypes  []string
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	buildCmd := "go build "

	// Configure LD flags
	ldFlags := ""
	for key, value := range Config.LdFlags {
		ldFlags += fmt.Sprintf("-X %s=%s ", key, value)
	}

	if ldFlags != "" {
		buildCmd += fmt.Sprintf("-ldflags='%s'", ldFlags)
	}

	output, err := utils.ExecCmd(utils.ExecConfig{
		Command: buildCmd,
	})

	// Print Output
	fmt.Println(output)

	return err
}

func Test() error {
	fmt.Println(Config.OutputDir)
	fmt.Println(Config.CgoEnabled)
	fmt.Println(Config.LdFlags)
	fmt.Println(Config.OsTypes)
	fmt.Println(Config.ArchTypes)
	return nil
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return os.Rename("./MyApp", "/usr/bin/MyApp")
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "get", "github.com/stretchr/piglatin")
	return cmd.Run()
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("MyApp")
}
