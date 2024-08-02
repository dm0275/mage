package magefile

import (
	"fmt"
	"github.com/dm0275/mage/utils"
	"os"
	"runtime/debug"
)

const (
	Perm0664 = 0o664
	Perm0755 = 0o755
)

var Config = &ProjectConfig{
	ProjectName: "",
	OutputDir:   "bin",
	CgoEnabled:  false,
	LdFlags:     map[string]string{},
	OsTypes:     []string{"linux"},
	ArchTypes:   []string{"amd64"},
}

type ProjectConfig struct {
	ProjectName string
	OutputDir   string
	CgoEnabled  bool
	LdFlags     map[string]string
	OsTypes     []string
	ArchTypes   []string
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	if Config.ProjectName == "" {
		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {
			return fmt.Errorf("Failed to read build info")
		}

		Config.ProjectName = buildInfo.Main.Path
	}

	goCmd := "go"
	buildCmd := []string{"build", "-v"}

	// Configure LD flags
	ldFlags := ""
	for key, value := range Config.LdFlags {
		ldFlags += fmt.Sprintf("-X %s=%s ", key, value)
	}

	if ldFlags != "" {
		buildCmd = append(buildCmd, fmt.Sprintf("-ldflags='%s'", ldFlags))
	}

	// Set output dir
	err := os.MkdirAll(Config.OutputDir, Perm0755)
	if err != nil {
		return fmt.Errorf("unable to create output directory %s. ERROR: %s", Config.OutputDir, err)
	}

	buildCmd = append(buildCmd, "-o",
		fmt.Sprintf("%s/%s", Config.OutputDir, Config.ProjectName),
		".",
	)

	// Run the build command
	output, err := utils.ExecCmd(utils.ExecConfig{
		Command: goCmd,
		Args:    buildCmd,
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

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("MyApp")
}
