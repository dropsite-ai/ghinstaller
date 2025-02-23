package ghinstaller

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"sync"
)

// Run uploads the binaries to the remote uploadDir using scp with no-clobber (-n).
func Run(publicDNS, keyPath, sshUser, uploadDir string, binPaths []string) error {
	log.Printf("Uploading binaries to server %s@%s into directory %s", sshUser, publicDNS, uploadDir)

	// Upload files concurrently using scp with no-clobber.
	var wg sync.WaitGroup
	for _, binPath := range binPaths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if _, err := ExecuteSCPCommand(publicDNS, keyPath, sshUser, uploadDir, p); err != nil {
				log.Printf("Failed to transfer %s: %v", p, err)
			}
		}(binPath)
	}
	wg.Wait()

	log.Println("Successfully uploaded binaries")
	return nil
}

// ExecuteSCPCommand uploads a file to the remote server using scp with the no-clobber flag (-n).
func ExecuteSCPCommand(publicDNS, keyPath, sshUser, uploadDir, binPath string) (string, error) {
	name := filepath.Base(binPath)
	cmd := fmt.Sprintf("scp -n -i %s %s %s@%s:%s/%s", keyPath, binPath, sshUser, publicDNS, uploadDir, name)
	return ExecuteCommand(cmd)
}

// ExecuteCommand runs a shell command.
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
