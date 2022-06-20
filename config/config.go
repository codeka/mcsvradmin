// Package config contains our configuration logic.
package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

// ServerMode represents the "mode" of the server, forge, vanilla, etc.
type ServerMode int

const (
	// Vanilla is a vanilla, unmodified server. By default, we'll look for minecraft_server*.jar
	// to run.
	Vanilla ServerMode = 1

	// Forge is a forge server. By default, we'll look for forge-*.jar to run.
	Forge ServerMode = 2
)

// Config contains the confiuration for McSvrAdmin.
type Config struct {
	BaseDirectory  string
	Mode           ServerMode
	JarFilePattern string
	JavaPath       string
	ListenPort     int
}

// Load parses the command-line and returns a Config instance with the config data.
func Load() (*Config, error) {
	baseDirectory, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %v", err)
	}

	args := os.Args[1:]
	if len(args) == 1 {
		baseDirectory = args[0]
	} else if len(args) > 1 {
		return nil, fmt.Errorf("invalid command-line flags: %v", args)
	}

	configFileName := path.Join(baseDirectory, "mcsvradmin.properties")
	log.Printf("Loading '%s'", configFileName)

	file, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{
		BaseDirectory: baseDirectory,
	}

	scanner := bufio.NewScanner(file)
	lineNo := 0
	defJarFilePattern := ""
	for scanner.Scan() {
		fullLine := scanner.Text()
		lineNo++
		idx := strings.Index(fullLine, "#")
		line := fullLine
		if idx >= 0 {
			line = line[:idx]
		}
		if strings.TrimSpace(line) == "" {
			continue
		}

		idx = strings.Index(line, "=")
		if idx < 0 {
			return nil, fmt.Errorf("invalid line %d: [%s]", lineNo, fullLine)
		}
		name := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])

		switch strings.ToLower(name) {
		case "mode":
			switch strings.ToLower(value) {
			case "vanilla":
				cfg.Mode = Vanilla
				defJarFilePattern = "minecraft_server*.jar"
				break
			case "forge":
				cfg.Mode = Forge
				defJarFilePattern = "forge-*.jar"
				break
			default:
				return nil, fmt.Errorf("invalid 'mode': %s", value)
			}
			break
		case "jarname":
			cfg.JarFilePattern = value
			break
		case "javadir":
			cfg.JavaPath = value
		case "listen":
			if value == "" {
				break
			}
			port, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid 'listen' (expected integer): %s: %v", value, err)
			}
			cfg.ListenPort = port
			break
		default:
			return nil, fmt.Errorf("invalid name: %s", name)
		}
	}

	if cfg.JarFilePattern == "" {
		cfg.JarFilePattern = defJarFilePattern
	}
	if cfg.ListenPort == 0 {
		cfg.ListenPort = 8080
	}

	return cfg, nil
}
