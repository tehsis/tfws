package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var command string
	ws := ""

	switch len(os.Args) {
	case 2:
		command = os.Args[1]
	case 3:
		ws = os.Args[1]
		command = os.Args[2]
	default:
		command = ""
	}

	switch command {
	case "list":
		listWorkspaces()
	default:
		fmt.Printf("No command specified %s", ws)
	}
}

func listWorkspaces() {
	workspacesRawList, err := runTFCommand("", "workspace", "list")

	if err != nil {
		log.Fatal("Error getting workspaces ", workspacesRawList)
	}

	// For some reason terraform workspace list returns workspaces multiple times
	workspacesList := dedup(strings.Split(workspacesRawList, "\n"))

	for _, o := range workspacesList {
		ws := strings.TrimSpace(strings.Trim(o, "*"))

		if ws != "" {
			_, err := runTFCommand(ws, "output")

			if err == nil {
				fmt.Printf("✅ %s\n", ws)
			} else {
				fmt.Printf("❌ %s\n", ws)
			}
		}
	}
}

func runTFCommand(ws string, tfCmd ...string) (string, error) {
	cmd := exec.Command("terraform", tfCmd...)

	if ws != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("TF_WORKSPACE=%s", ws))
	}

	out, err := cmd.CombinedOutput()

	return string(out), err
}

func dedup(elements []string) []string {
	found := map[string]bool{}
	output := []string{}

	for _, e := range elements {
		// The current workspace contains an * on the output
		se := strings.TrimSpace(strings.Replace(e, "*", "", -1))
		if found[se] {
		} else {
			found[se] = true
			output = append(output, se)
		}
	}

	return output
}
