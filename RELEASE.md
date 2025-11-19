# Release Guide for SkyVault

This guide describes the process for creating a new release of SkyVault.

## Prerequisites

- Write access to the GitHub repository
- All tests passing on the main branch
- All planned features for the release completed

## Release Process

### 1. Prepare the Release

1. **Update version numbers** (if applicable):
   - `web/package.json` - Update version field
   - `README.md` - Update version badge

2. **Update CHANGELOG** (create if it doesn't exist):
   - Document all changes since the last release
   - Group changes by category (Features, Bug Fixes, Breaking Changes, etc.)
   - Include links to relevant PRs and issues

3. **Test the build locally**:
   ```bash
   # Clean build
   task nuke

   # Build everything
   task build

   # Run tests
   task test

   # Test Docker build
   docker build -t skyvault:test .
   ```

4. **Commit and push changes**:
   ```bash
   git add .
   git commit -m "Prepare release v1.0.0"
   git push origin main
   ```

### 2. Create the Release on GitHub

1. **Navigate to the repository** on GitHub

2. **Go to Releases**:
   - Click on "Releases" in the right sidebar
   - Click "Draft a new release"

3. **Create a new tag**:
   - Click "Choose a tag"
   - Enter the version number (e.g., `v1.0.0`)
   - Click "Create new tag: v1.0.0 on publish"

4. **Fill in release details**:
   - **Release title**: `SkyVault v1.0.0` (or appropriate version)
   - **Description**: Use this template:

   ```markdown
   ## What's New in v1.0.0

   ### Features
   - 🔐 JWT-based authentication system
   - 📁 Folder creation and navigation
   - 📤 File upload with chunked upload support (up to 10GB)
   - 📥 File download
   - 📱 Mobile-first responsive UI
   - 🐳 Docker-based deployment

   ### Installation

   Quick start with Docker:

   \`\`\`bash
   # Download configuration files
   wget https://raw.githubusercontent.com/yourusername/skyvault/v1.0.0/.env.example -O .env
   wget https://raw.githubusercontent.com/yourusername/skyvault/v1.0.0/docker-compose.prod.yml

   # Edit .env with your settings
   nano .env

   # Start SkyVault
   docker compose -f docker-compose.prod.yml up -d
   \`\`\`

   See the [README](https://github.com/yourusername/skyvault/blob/v1.0.0/README.md) for detailed installation and configuration instructions.

   ### Docker Image

   The Docker image is available at:
   ```
   ghcr.io/yourusername/skyvault:v1.0.0
   ```

   ### Full Changelog

   See [CHANGELOG.md](https://github.com/yourusername/skyvault/blob/v1.0.0/CHANGELOG.md) for complete details.
   ```

5. **Set as latest release**: Check the box "Set as the latest release"

6. **Publish release**: Click "Publish release"

### 3. Automated Build

Once you publish the release, GitHub Actions will automatically:

1. Build the Docker image for both `linux/amd64` and `linux/arm64`
2. Push the image to GitHub Container Registry with tags:
   - `ghcr.io/yourusername/skyvault:v1.0.0`
   - `ghcr.io/yourusername/skyvault:v1`
   - `ghcr.io/yourusername/skyvault:latest`
3. Generate build attestation for security

You can monitor the build progress:
- Go to the "Actions" tab in your GitHub repository
- Click on the "Build and Publish Release" workflow
- Watch the build progress

### 4. Verify the Release

1. **Wait for the build to complete** (usually 5-10 minutes)

2. **Test the Docker image**:
   ```bash
   # Pull the image
   docker pull ghcr.io/yourusername/skyvault:v1.0.0

   # Test with docker-compose
   docker compose -f docker-compose.prod.yml up -d

   # Check health
   curl http://localhost:8090/api/v1/pub/health

   # Test in browser
   open http://localhost:8090
   ```

3. **Verify the release on GitHub**:
   - Check that the release is marked as "Latest"
   - Verify all links in the release notes work
   - Confirm the Docker image appears in the Packages section

### 5. Announce the Release

1. **Update social media** (if applicable):
   - Twitter/X
   - Reddit (r/selfhosted)
   - Discord communities
   - Dev.to or Hashnode

2. **Notify users**:
   - GitHub Discussions announcement
   - Newsletter (if you have one)
   - Documentation site update

## Hotfix Releases

For urgent bug fixes:

1. Create a branch from the release tag:
   ```bash
   git checkout -b hotfix/v1.0.1 v1.0.0
   ```

2. Make the fix and commit:
   ```bash
   git commit -m "Fix critical bug in authentication"
   ```

3. Merge back to main:
   ```bash
   git checkout main
   git merge hotfix/v1.0.1
   git push origin main
   ```

4. Create the hotfix release following steps above

## Version Numbering

SkyVault follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version (v2.0.0): Incompatible API changes
- **MINOR** version (v1.1.0): New functionality, backwards compatible
- **PATCH** version (v1.0.1): Backwards compatible bug fixes

## Rollback Procedure

If you need to roll back a release:

1. **Mark the release as pre-release** on GitHub (don't delete it)

2. **Create a new release** with the previous stable version

3. **Notify users** about the rollback and provide instructions

4. **Fix the issues** and create a new release

## Checklist

Use this checklist for each release:

- [ ] All tests passing
- [ ] Version numbers updated
- [ ] CHANGELOG updated
- [ ] Local build tested
- [ ] Docker build tested locally
- [ ] Changes committed and pushed
- [ ] GitHub release created
- [ ] Release notes written
- [ ] Release published
- [ ] GitHub Actions build successful
- [ ] Docker image tested
- [ ] Release announced

## Troubleshooting

### Build fails in GitHub Actions

1. Check the Actions logs for error messages
2. Test the build locally: `docker build -t skyvault:test .`
3. Ensure all files are committed
4. Verify the Dockerfile is correct

### Docker image not appearing in Packages

1. Check that the workflow completed successfully
2. Verify repository settings allow package publishing
3. Check that GITHUB_TOKEN has write:packages permission

### Users can't pull the image

1. Verify the package visibility is set to "Public"
2. Go to the package settings
3. Change visibility from "Private" to "Public"

## Support

If you encounter issues with the release process:

1. Check GitHub Actions logs
2. Review this guide
3. Open an issue in the repository
4. Contact the maintainers
