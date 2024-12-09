package parser

import (
	"fmt"
	"os/exec"
	"strings"
)

var recursiveCurrentDir bool

func runAg(text string, query string) ([]string, error) {
	// Start with the base ag command and necessary flags
	cmdArgs := []string{"--nocolor", "--no-filename", "-o", query}

	// Determine the input source
	switch {
	case text != "":
		cmd := exec.Command("ag", cmdArgs...)
		cmd.Stdin = strings.NewReader(text)
		return captureAgOutput(cmd, query)
	case recursiveCurrentDir:
		// Append "." to the command arguments if we want to search the current directory
		cmdArgs = append(cmdArgs, ".")
		cmd := exec.Command("ag", cmdArgs...)
		return captureAgOutput(cmd, query)
	default:
		cmd := exec.Command("ag", cmdArgs...)
		cmd.Stdin = strings.NewReader(text)
		return captureAgOutput(cmd, query)
	}
}

// Helper function to capture output from the ag command
func captureAgOutput(cmd *exec.Cmd, query string) ([]string, error) {
	// Capture the output of ag
	output, err := cmd.CombinedOutput()

	// Split the output by newline to get individual search results
	lines := strings.Split(string(output), "\n")

	// Remove the last empty line if any
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Handle "no matches" case separately
	if err != nil {
		// Check if the error is due to "no matches" (exit code 1)
		exitError, ok := err.(*exec.ExitError)
		if ok && exitError.ExitCode() == 1 && len(lines) == 0 {
			// Treat this as a valid scenario and return an empty list
			// fmt.Println("Query not found: ", query)
			return []string{}, nil
		}
		// For other errors, return the error as usual
		return nil, fmt.Errorf("error: %v, output: %s", err, string(output))
	}
	return lines, nil
}
