package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hyturing/compose-pulse/internal/compose"
	"github.com/hyturing/compose-pulse/internal/dag"
	"github.com/hyturing/compose-pulse/internal/docker"
	"github.com/hyturing/compose-pulse/internal/ui"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "cpulse: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	file := flag.String("file", "", "Path to docker-compose.yml (default: auto-detect from CWD)")
	ver := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *ver {
		fmt.Printf("cpulse %s\n", version)
		return nil
	}

	composePath, err := compose.Locate(*file)
	if err != nil {
		return err
	}

	cfg, err := compose.Parse(composePath)
	if err != nil {
		return fmt.Errorf("error parsing compose file: %w", err)
	}

	graph, err := dag.Build(cfg)
	if err != nil {
		return err
	}

	dockerClient, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("error connecting to Docker daemon: %w", err)
	}
	defer func() { _ = dockerClient.Close() }()

	model := ui.New(graph, dockerClient)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
