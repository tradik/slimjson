# Release Process

This document describes how to create a new release of SlimJSON.

## Prerequisites

1. All tests must pass
2. Documentation is up to date
3. CHANGELOG.md is updated with release notes
4. You have push access to the repository

## Release Steps

### 1. Update Version Information

Update `CHANGELOG.md`:

```markdown
## [1.0.0] - 2024-01-15

### Added
- Feature 1
- Feature 2

### Changed
- Change 1

### Fixed
- Bug fix 1
```

### 2. Commit Changes

```bash
git add CHANGELOG.md
git commit -m "chore: prepare release v1.0.0"
git push origin main
```

### 3. Create and Push Tag

```bash
# Create annotated tag
git tag -a v1.0.0 -m "Release v1.0.0"

# Push tag to GitHub
git push origin v1.0.0
```

### 4. Automated Release Process

Once the tag is pushed, GitHub Actions will automatically:

1. ✅ Run all tests
2. ✅ Run golangci-lint
3. ✅ Build binaries for:
   - linux/amd64
   - linux/arm64
   - darwin/amd64
   - darwin/arm64
   - freebsd/amd64
   - freebsd/arm64
4. ✅ Create GitHub Release with:
   - Release notes from CHANGELOG.md
   - Binary artifacts
   - SHA256 checksums
5. ✅ Publish to pkg.go.dev
6. ✅ Build and push Docker images to ghcr.io

### 5. Verify Release

After the GitHub Action completes (usually 5-10 minutes):

#### Check GitHub Release

```bash
# Visit GitHub releases page
https://github.com/tradik/slimjson/releases
```

Verify:
- ✅ Release notes are correct
- ✅ All binary artifacts are present
- ✅ SHA256 checksums are generated

#### Check pkg.go.dev

```bash
# Visit pkg.go.dev
https://pkg.go.dev/github.com/tradik/slimjson@v1.0.0
```

Verify:
- ✅ Documentation is rendered correctly
- ✅ Version is listed
- ✅ Examples are visible

**Note:** pkg.go.dev may take 5-15 minutes to index the new version.

#### Check Docker Image

```bash
# Pull the image
docker pull ghcr.io/tradik/slimjson:v1.0.0
docker pull ghcr.io/tradik/slimjson:latest

# Test the image
docker run --rm ghcr.io/tradik/slimjson:v1.0.0 --help
```

#### Test Installation

```bash
# Create a test directory
mkdir /tmp/test-slimjson
cd /tmp/test-slimjson

# Initialize Go module
go mod init test

# Install the new version
go get github.com/tradik/slimjson@v1.0.0

# Verify it works
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "github.com/tradik/slimjson"
)

func main() {
    data := map[string]interface{}{
        "test": "data",
        "list": []interface{}{1, 2, 3, 4, 5},
    }
    
    cfg := slimjson.Config{
        MaxListLength: 3,
        StripEmpty:    true,
    }
    
    slimmer := slimjson.New(cfg)
    result := slimmer.Slim(data)
    
    out, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(out))
}
EOF

go run main.go
```

### 6. Announce Release

After verification, announce the release:

1. **GitHub Discussions** (if enabled)
2. **Twitter/Social Media**
3. **Go Community** (reddit.com/r/golang)
4. **Project Users** (if applicable)

## Version Numbering

SlimJSON follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version (v1.0.0 → v2.0.0): Incompatible API changes
- **MINOR** version (v1.0.0 → v1.1.0): New functionality, backwards compatible
- **PATCH** version (v1.0.0 → v1.0.1): Bug fixes, backwards compatible

## Pre-release Versions

For testing before official release:

```bash
# Create pre-release tag
git tag -a v1.0.0-rc.1 -m "Release candidate 1"
git push origin v1.0.0-rc.1
```

Pre-release versions:
- `v1.0.0-alpha.1` - Alpha release
- `v1.0.0-beta.1` - Beta release
- `v1.0.0-rc.1` - Release candidate

## Hotfix Process

For urgent bug fixes:

1. Create hotfix branch from tag:
   ```bash
   git checkout -b hotfix/v1.0.1 v1.0.0
   ```

2. Make fixes and commit:
   ```bash
   git commit -m "fix: critical bug"
   ```

3. Update CHANGELOG.md

4. Create tag:
   ```bash
   git tag -a v1.0.1 -m "Hotfix v1.0.1"
   ```

5. Push branch and tag:
   ```bash
   git push origin hotfix/v1.0.1
   git push origin v1.0.1
   ```

6. Merge back to main:
   ```bash
   git checkout main
   git merge hotfix/v1.0.1
   git push origin main
   ```

## Rollback

If a release has critical issues:

1. **Delete the tag** (if not widely used):
   ```bash
   git tag -d v1.0.0
   git push origin :refs/tags/v1.0.0
   ```

2. **Create a new patch release** (if already in use):
   ```bash
   # Fix the issue
   git commit -m "fix: critical issue"
   
   # Create new version
   git tag -a v1.0.1 -m "Fix for v1.0.0"
   git push origin v1.0.1
   ```

## Troubleshooting

### pkg.go.dev not updating

```bash
# Manually trigger update
curl "https://proxy.golang.org/github.com/tradik/slimjson/@v/v1.0.0.info"

# Wait 5-15 minutes and check again
```

### GitHub Action fails

1. Check the Actions tab on GitHub
2. Review the error logs
3. Fix the issue
4. Delete the tag and recreate:
   ```bash
   git tag -d v1.0.0
   git push origin :refs/tags/v1.0.0
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

### Binary artifacts missing

1. Check if the build job completed successfully
2. Verify the matrix configuration in `.github/workflows/release.yml`
3. Re-run the failed jobs from GitHub Actions UI

## Checklist

Before creating a release, verify:

- [ ] All tests pass locally (`go test ./...`)
- [ ] golangci-lint passes (`golangci-lint run`)
- [ ] Documentation is updated
  - [ ] README.md
  - [ ] LIBRARY_EXAMPLES.md
  - [ ] doc.go
  - [ ] api/README.md
  - [ ] api/swagger.yaml
- [ ] CHANGELOG.md has release notes
- [ ] Version number follows semver
- [ ] No uncommitted changes
- [ ] On main branch
- [ ] All CI checks pass on main

## Post-Release

After successful release:

- [ ] Verify GitHub Release
- [ ] Verify pkg.go.dev
- [ ] Verify Docker images
- [ ] Test installation
- [ ] Update project board (if applicable)
- [ ] Announce release
- [ ] Close related issues/PRs

## Resources

- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [Go Modules](https://go.dev/blog/publishing-go-modules)
- [pkg.go.dev](https://pkg.go.dev/about)
