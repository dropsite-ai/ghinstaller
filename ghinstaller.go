package ghinstaller

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func Run(publicDNS, keyPath, sshUser, homeDir string, binPaths []string) error {
	log.Printf("Deploying binaries to server %s@%s", sshUser, publicDNS)

	// Set up home and bin directories
	setupDirsCmd := fmt.Sprintf(`
		mkdir -p %s/bin %s/bin.old &&
		chmod 0700 %s %s/bin %s/bin.old`,
		homeDir, homeDir, homeDir, homeDir, homeDir,
	)
	if _, err := ExecuteSSHCommand(publicDNS, keyPath, sshUser, setupDirsCmd); err != nil {
		return fmt.Errorf("failed to create and set up bin directories: %v", err)
	}

	// Upload tar.gz files concurrently
	var wg sync.WaitGroup
	for _, binPath := range binPaths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if _, err := ExecuteSCPCommand(publicDNS, keyPath, sshUser, homeDir, p); err != nil {
				fmt.Printf("Failed to transfer %s: %v\n", p, err)
			}
		}(binPath)
	}
	wg.Wait()

	// Extract tar.gz files and install any resulting binaries
	if err := InstallBinaries(publicDNS, keyPath, sshUser, homeDir, binPaths); err != nil {
		return fmt.Errorf("failed to install binaries: %v", err)
	}

	log.Println("Successfully deployed and installed binaries")
	return nil
}

func InstallBinaries(publicDNS, keyPath, sshUser, homeDir string, binPaths []string) error {
	for _, tarPath := range binPaths {
		tarFileName := filepath.Base(tarPath)
		tarFileNameNoExt := strings.TrimSuffix(tarFileName, ".tar.gz")

		installCmd := fmt.Sprintf(`
			set -e
			SOURCE_BIN_DIR=%s/bin
			DEST_BIN_DIR=/usr/local/bin
			TAR_FILE=%s
			EXTRACT_DIR=%s/%s

			# Extract tar.gz file
			mkdir -p "$EXTRACT_DIR"
			tar -xzf "$SOURCE_BIN_DIR/$TAR_FILE" -C "$EXTRACT_DIR"

			# Find and install binaries
			for file in "$EXTRACT_DIR"/*; do
				if [ -x "$file" ] && [ ! -d "$file" ]; then
					BINARY_NAME=$(basename "$file")
					# Backup existing binary if it exists
					if [ -f "$DEST_BIN_DIR/$BINARY_NAME" ]; then
						mv "$DEST_BIN_DIR/$BINARY_NAME" "%s/bin.old/$BINARY_NAME"
						chmod 0700 "%s/bin.old/$BINARY_NAME"
					fi

					# Move new binary to destination
					cp "$file" "$DEST_BIN_DIR/$BINARY_NAME"
					chown root:root "$DEST_BIN_DIR/$BINARY_NAME"
					chmod 0755 "$DEST_BIN_DIR/$BINARY_NAME"
				fi
			done

			# Cleanup
			rm -rf "$EXTRACT_DIR"
		`, homeDir, tarFileName, homeDir, tarFileNameNoExt, homeDir, homeDir)

		if _, err := ExecuteSSHCommand(publicDNS, keyPath, sshUser, installCmd); err != nil {
			return fmt.Errorf("failed to install binaries from %s: %v", tarFileName, err)
		}
	}

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
