package main

import (
	"context"
	"fmt"

	"os"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
)

func main() {
	ctx := context.Background()

	// Inline Compose file definition
	composeFile := []byte(`
services:
  hello:
    build:
      context: .
      dockerfile: Dockerfile
    command: ["echo", "Hello from Compose SDK!"]
`)
	// Create a Compose project from inline YAML
	// Write the inline YAML to a temporary file
	tmpFile, _ := os.CreateTemp("", "compose-*.yaml")

	defer os.Remove(tmpFile.Name())
	tmpFile.Write(composeFile)
	tmpFile.Close()

	projectOpts, _ := cli.NewProjectOptions([]string{}, func(opts *cli.ProjectOptions) error {
		opts.WorkingDir = "."
		opts.ConfigPaths = []string{tmpFile.Name()}
		return nil
	})
	project, _ := cli.ProjectFromOptions(ctx, projectOpts)
	dockerCli, _ := command.NewDockerCli()
	dockerCli.Initialize(flags.NewClientOptions())

	composeService := compose.NewComposeService(dockerCli)

	serviceNames := []string{"hello"}

	// The Error is raised from this build step:
	// failed to build: no builder "desktop-linux" found
	err := composeService.Build(ctx, project, api.BuildOptions{
		Services: serviceNames,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building project: %v\n", err)
		return
	}

	fmt.Println("Build completed successfully!")
}
