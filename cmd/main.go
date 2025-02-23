package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dropsite-ai/ghdownloader"
	"github.com/dropsite-ai/ghinstaller"
)

// repoList implements flag.Value to allow multiple -repo flags.
type repoList []string

func (r *repoList) String() string {
	return strings.Join(*r, ",")
}

func (r *repoList) Set(value string) error {
	*r = append(*r, value)
	return nil
}

func main() {
	// Define command-line flags
	token := flag.String("token", "", "GitHub Personal Access Token. Defaults to GITHUB_TOKEN environment variable if not provided.")
	downloadDir := flag.String("dl", "./downloads", "Destination directory for downloaded binaries.")
	var repos repoList
	flag.Var(&repos, "repo", "Repository in 'owner/repo' format. Can be specified multiple times. (Required)")
	sshUser := flag.String("sshuser", "ec2-user", "SSH username for the target server. (Default: ec2-user)")
	uploadDir := flag.String("uploads", "/home/ec2-user/uploads", "Upload directory on the target server. (Default: /home/ec2-user/uploads)")
	publicDNS := flag.String("publicdns", "", "Public DNS of the EC2 instance for deployment. (Required)")
	keyPath := flag.String("key", "", "Path to the private key for SCP. (Required)")
	match := flag.String("match", "", "Substring to filter assets by name during download. (Optional)")

	flag.Parse()

	// Use GITHUB_TOKEN environment variable if token flag is empty
	if *token == "" {
		*token = os.Getenv("GITHUB_TOKEN")
	}

	// Validate required arguments
	if len(repos) == 0 {
		log.Fatal("Error: At least one repository is required. Specify repositories using the -repo flag.")
	}
	if *publicDNS == "" {
		log.Fatal("Error: Public DNS of the EC2 instance is required. Specify it using the -publicdns flag.")
	}
	if *keyPath == "" {
		log.Fatal("Error: Path to the private key for SCP is required. Specify it using the -key flag.")
	}

	// Initialize the downloader
	downloader := ghdownloader.New(*token, *downloadDir)
	downloader.SetMatchFilter(*match)

	// Download the latest releases
	fmt.Println("Starting download...")
	binPaths, err := downloader.DownloadLatestReleases(repos)
	if err != nil {
		log.Fatalf("Error downloading binaries: %v", err)
	}

	// Upload the binaries using SCP (with no-clobber)
	fmt.Println("Starting upload...")
	if err := ghinstaller.Run(*publicDNS, *keyPath, *sshUser, *uploadDir, binPaths); err != nil {
		log.Fatalf("Error uploading binaries: %v", err)
	}

	fmt.Println("Upload completed successfully.")
}
