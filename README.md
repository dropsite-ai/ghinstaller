# ghinstaller

Install GitHub release binaries onto a server.

## Install

Download from [Releases](https://github.com/dropsite-ai/ghinstaller/releases):

```bash
tar -xzf ghinstaller_Darwin_arm64.tar.gz
chmod +x ghinstaller
sudo mv ghinstaller /usr/local/bin/
```

Or manually build and install:

```bash
git clone git@github.com:dropsite-ai/ghinstaller.git
cd ghinstaller
make install
```

## Usage

```bash
  -dest string
    	Destination directory for downloaded binaries. (default "./binaries")
  -key string
    	Path to the private key for SSH. (Required)
  -match string
    	Substring to filter assets by name during download. (Optional)
  -publicDNS string
    	Public DNS of the EC2 instance for deployment. (Required)
  -repos string
    	Comma-separated list of GitHub repositories in 'owner/repo' format. (Required)
  -token string
    	GitHub Personal Access Token. Defaults to GITHUB_TOKEN environment variable if not provided.
```

## Test

```bash
make test
```

## Release

```bash
make release
```