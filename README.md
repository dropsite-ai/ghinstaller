# ghinstaller

Copy GitHub release binaries onto a server via SCP.

## Installation

### Go Package

```bash
go get github.com/dropsite-ai/ghinstaller
```

### Homebrew (macOS or Compatible)

If you use Homebrew, install ghinstaller with:
```bash
brew tap dropsite-ai/homebrew-tap
brew install ghinstaller
```

### Download Binaries

Grab the latest pre-built binaries from the [GitHub Releases](https://github.com/dropsite-ai/ghinstaller/releases). Extract them, then run the `ghinstaller` executable directly.

### Build from Source

1. **Clone the repository**:
   ```bash
   git clone https://github.com/dropsite-ai/ghinstaller.git
   cd ghinstaller
   ```
2. **Build using Go**:
   ```bash
   go build -o ghinstaller cmd/main.go
   ```

## Usage

### Command-Line

Deploy your downloaded binaries to a remote server via SCP (with no-clobber enabled) using:

```bash
ghinstaller -repo owner/repo -repo anotherOwner/anotherRepo \
  -dl "./downloads" -uploads "/home/ec2-user/uploads" \
  -key /path/to/key.pem -publicdns example.com -sshuser ec2-user \
  -token YOUR_GITHUB_TOKEN -match "linux"
```

- **-repo**: Specify one repository per flag in the format `owner/repo`. Use multiple `-repo` flags for multiple repositories.
- **-dl**: Local directory where binaries have been downloaded (default: `./downloads`).
- **-uploads**: Remote upload directory (default: `/home/ec2-user/uploads`).
- **-key**: Path to the private key for SCP (required).
- **-publicdns**: Public DNS or IP address of the target server (required).
- **-sshuser**: SSH username for the target server (default: `ec2-user`).
- **-token**: GitHub Personal Access Token (or use `GITHUB_TOKEN` environment variable).
- **-match**: (Optional) Filter asset names during download.

### Programmatic Usage

Below is an example of using the installer within a Go application. This snippet downloads assets using `ghdownloader` and then uploads them using `ghinstaller`:

```go
package main

import (
    "fmt"
    "log"

    "github.com/dropsite-ai/ghdownloader"
    "github.com/dropsite-ai/ghinstaller"
)

func main() {
    token := "YOUR_GITHUB_TOKEN"   // Enter token or leave empty (unauthenticated).
    downloadDir := "./downloads"
    uploadDir := "/home/ec2-user/uploads"
    repos := []string{"owner/repo", "anotherOwner/anotherRepo"}
    publicDNS := "example.com"
    keyPath := "/path/to/key.pem"
    sshUser := "ec2-user"
    match := "linux"

    // Download the latest releases.
    downloader := ghdownloader.New(token, downloadDir)
    downloader.SetMatchFilter(match)
    binPaths, err := downloader.DownloadLatestReleases(repos)
    if err != nil {
        log.Fatalf("Download error: %v", err)
    }
    
    // Upload the binaries using SCP (with no-clobber).
    if err := ghinstaller.Run(publicDNS, keyPath, sshUser, uploadDir, binPaths); err != nil {
        log.Fatalf("Upload error: %v", err)
    }
    
    fmt.Println("Deployment completed successfully.")
}
```

## Test

```bash
make test
```

## Release

```bash
make release
```