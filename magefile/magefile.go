package magefile

import (
	"fmt"
	"github.com/dm0275/mage/utils"
	l "log"
	"os"
)

const (
	Perm0664 = 0o664
	Perm0755 = 0o755
)

var debugEnabled = false
var logger = l.New(os.Stdout, "", 0)

var Config = &ProjectConfig{
	ProjectName: "",
	OutputDir:   "bin",
	CgoEnabled:  false,
	LdFlags:     map[string]string{},
	OsTypes:     []string{"linux"},
	ArchTypes:   []string{"amd64"},
}

type ProjectConfig struct {
	ProjectName  string
	OutputBinary string
	OutputDir    string
	CgoEnabled   bool
	LdFlags      map[string]string
	OsTypes      []string
	ArchTypes    []string
}

func configureDebugSettings() {
	if os.Getenv("MAGEFILE_DEBUG") == "1" || os.Getenv("MAGEFILE_VERBOSE") == "1" {
		debugEnabled = true
	}
}

func init() {
	// Configure debug settings
	configureDebugSettings()
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	if Config.ProjectName == "" {
		logger.Fatalf("no ProjectName defined")
	}

	goCmd := "go"
	buildCmd := []string{"build", "-v"}

	// Configure LD flags
	ldFlags := ""
	for key, value := range Config.LdFlags {
		ldFlags += fmt.Sprintf("-X %s=%s ", key, value)
	}

	if ldFlags != "" {
		buildCmd = append(buildCmd, fmt.Sprintf("-ldflags=%s", ldFlags))
	}

	// Set output dir
	err := os.MkdirAll(Config.OutputDir, Perm0755)
	if err != nil {
		return fmt.Errorf("unable to create output directory %s. ERROR: %s", Config.OutputDir, err)
	}

	var output string

	// Run the build command
	for _, osType := range Config.OsTypes {
		for _, archType := range Config.ArchTypes {
			logger.Printf("Building %s-%s-%s...", Config.ProjectName, osType, archType)

			buildOsCmd := append(buildCmd, "-o",
				fmt.Sprintf("%s/%s-%s-%s", Config.OutputDir, Config.ProjectName, osType, archType),
				".",
			)

			if debugEnabled {
				logger.Printf("Executing binary [%s] with args: %s", goCmd, buildOsCmd)
			}
			output, err = utils.ExecCmd(utils.ExecConfig{
				Environment: []string{
					fmt.Sprintf("GOOS=%s", osType),
					fmt.Sprintf("GOARCH=%s", archType),
					fmt.Sprintf("CGO_ENABLED=%s", Config.CgoEnabled),
				},
				Command: goCmd,
				Args:    buildOsCmd,
			})

			if err != nil {
				return err
			}

			if debugEnabled {
				logger.Printf("Build Output: %s", output)
			}

			fmt.Println("Done")
		}
	}

	return err
}

func Test() error {
	// TODO: add test task
	return nil
}

// Clean up
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll(Config.OutputDir)
}
