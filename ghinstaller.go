package ghinstaller

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"sync"
)

func Run(publicDNS, keyPath, sshUser, homeDir string, binPaths []string) error {
	log.Printf("Deploying binaries to server %s@%s", sshUser, publicDNS)

	// Set up home directory
	createBinDirCmd := fmt.Sprintf("mkdir -p %s/bin %s/bin.old && chmod 0700 %s %s/bin %s/bin.old", homeDir, homeDir, homeDir, homeDir, homeDir)
	if _, err := ExecuteSSHCommand(publicDNS, keyPath, sshUser, createBinDirCmd); err != nil {
		return fmt.Errorf("failed to create bin directories: %v", err)
	}

	// Upload binaries concurrently
	var wg sync.WaitGroup
	for _, binPath := range binPaths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if _, err := ExecuteSCPCommand(publicDNS, keyPath, sshUser, homeDir, p); err != nil {
				fmt.Printf("failed to transfer %s: %v\n", p, err)
			}
		}(binPath)
	}

	wg.Wait()
	log.Println("Successfully deployed binaries")
	return nil
}

func ExecuteSSHCommand(publicDNS, keyPath, sshUser, command string) (string, error) {
	cmd := fmt.Sprintf(`ssh -i %s %s@%s '%s'`, keyPath, sshUser, publicDNS, command)
	return ExecuteCommand(cmd)
}

func ExecuteSCPCommand(publicDNS, keyPath, sshUser, homeDir, binPath string) (string, error) {
	name := filepath.Base(binPath)
	cmd := fmt.Sprintf("scp -i %s %s %s@%s:%s/bin/%s", keyPath, binPath, sshUser, publicDNS, homeDir, name)
	return ExecuteCommand(cmd)
}

func ExecuteCommand(command string) (string, error) {
	log.Printf("Executing command: %s", command)
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %s, output: %s", err, string(output))
	}

	log.Printf("Command output: %s", string(output))
	return string(output), nil
}
