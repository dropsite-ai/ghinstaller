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
    	Destination directory for downloaded binaries (default "./downloads")
  -match string
    	Substring to filter assets by name (optional)
  -repos string
    	Comma-separated list of GitHub repositories in 'owner/repo' format (required)
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