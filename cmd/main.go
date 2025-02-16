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

func main() {
	// Define command-line flags
	token := flag.String("token", "", "GitHub Personal Access Token. Defaults to GITHUB_TOKEN environment variable if not provided.")
	destDir := flag.String("dest", "./binaries", "Destination directory for downloaded binaries.")
	repos := flag.String("repos", "", "Comma-separated list of GitHub repositories in 'owner/repo' format. (Required)")
	sshUser := flag.String("sshuser", "ec2-user", "SSH username for the target server. (Default: ec2-user)")
	homeDir := flag.String("homedir", "/home/ec2-user", "Home directory on the target server. (Default: /home/ec2-user)")
	publicDNS := flag.String("publicdns", "", "Public DNS of the EC2 instance for deployment. (Required)")
	keyPath := flag.String("key", "", "Path to the private key for SSH. (Required)")
	match := flag.String("match", "", "Substring to filter assets by name during download. (Optional)")

	flag.Parse()

	// Use GITHUB_TOKEN environment variable if token flag is empty
	if *token == "" {
		*token = os.Getenv("GITHUB_TOKEN")
	}

	// Validate required arguments
	if *repos == "" {
		log.Fatal("Error: At least one repository is required. Specify repositories using the -repos flag.")
	}
	if *publicDNS == "" {
		log.Fatal("Error: Public DNS of the EC2 instance is required. Specify it using the -publicDNS flag.")
	}
	if *keyPath == "" {
		log.Fatal("Error: Path to the private key for SSH is required. Specify it using the -key flag.")
	}

	// Parse the repositories
	userRepos := strings.Split(*repos, ",")

	// Initialize the downloader
	downloader := ghdownloader.New(*token, *destDir)
	downloader.SetMatchFilter(*match) // Set the match filter for asset names

	// Download the latest releases
	fmt.Println("Starting download...")
	binPaths, err := downloader.DownloadLatestReleases(userRepos)
	if err != nil {
		log.Fatalf("Error downloading binaries: %v", err)
	}

	// Deploy the binaries
	fmt.Println("Starting deployment...")
	if err := ghinstaller.Run(*publicDNS, *keyPath, *sshUser, *homeDir, binPaths); err != nil {
		log.Fatalf("Error deploying binaries: %v", err)
	}

	fmt.Println("Deployment completed successfully.")
}
