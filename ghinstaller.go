package ghinstaller

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"sync"
)

func Run(publicDNS, keyPath string, binPaths []string) error {
	log.Printf("Deploying binaries to EC2 instance %s", publicDNS)

	// Set up home directory
	createBinDirCmd := "mkdir -p /home/ec2-user/bin /home/ec2-user/bin.old && chmod 0700 /home/ec2-user /home/ec2-user/bin /home/ec2-user/bin.old"
	if _, err := ExecuteSSHCommand(publicDNS, keyPath, createBinDirCmd); err != nil {
		return fmt.Errorf("failed to create bin directories: %v", err)
	}

	// Upload binaries concurrently
	var wg sync.WaitGroup
	for _, binPath := range binPaths {
		wg.Add(1) // Add a counter for each goroutine
		go func(p string) {
			defer wg.Done() // Decrease the counter when the goroutine completes
			if _, err := ExecuteSCPCommand(publicDNS, keyPath, p); err != nil {
				fmt.Printf("failed to transfer %s: %v\n", p, err)
			}
		}(binPath)
	}

	wg.Wait() // Wait for all goroutines to complete

	log.Println("Successfully deployed binaries")
	return nil
}

func ExecuteSSHCommand(publicDNS, keyPath, command string) (string, error) {
	return ExecuteCommand(fmt.Sprintf(`ssh -i %s ec2-user@%s '%s'`, keyPath, publicDNS, command))
}

func ExecuteSCPCommand(publicDNS, keyPath, binPath string) (string, error) {
	name := filepath.Base(binPath)
	return ExecuteCommand(fmt.Sprintf("scp -i %s %s ec2-user@%s:/home/ec2-user/bin/%s", keyPath, binPath, publicDNS, name))
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
