# GitHub Deployment Checklist

Use this checklist to ensure a smooth deployment to GitHub.

---

## Pre-Deployment

- [ ] All code is committed locally
- [ ] Tests pass (if any)
- [ ] Documentation is complete
- [ ] .gitignore is configured
- [ ] LICENSE file exists
- [ ] README.md is updated
- [ ] Binary files are excluded from git

---

## GitHub Setup

- [ ] GitHub account created
- [ ] New repository created on GitHub
- [ ] Repository name: `watchup-agent`
- [ ] Repository is public (or private if preferred)
- [ ] Repository description added

---

## Initial Push

```bash
# Check current status
git status

# Initialize if needed
git init

# Add all files
git add .

# Commit
git commit -m "Initial commit: Watchup Server Agent v1.0.0"

# Add remote (replace YOUR_USERNAME)
git remote add origin https://github.com/YOUR_USERNAME/watchup-agent.git

# Verify remote
git remote -v

# Push to GitHub
git branch -M main
git push -u origin main
```

- [ ] Code pushed to GitHub
- [ ] All files visible on GitHub
- [ ] README displays correctly

---

## Repository Configuration

### Topics/Tags
Add these topics to your repository:
- [ ] `monitoring`
- [ ] `server-monitoring`
- [ ] `golang`
- [ ] `systemd`
- [ ] `metrics`
- [ ] `alerting`

### About Section
- [ ] Description: "Lightweight server monitoring agent for Watchup platform"
- [ ] Website: https://watchup.site
- [ ] Topics added

### Branch Protection (Optional)
- [ ] Protect main branch
- [ ] Require pull request reviews
- [ ] Require status checks

---

## Documentation

- [ ] README.md is comprehensive
- [ ] Installation instructions are clear
- [ ] VPS deployment guide included
- [ ] Configuration examples provided
- [ ] Troubleshooting section added
- [ ] API documentation included

---

## GitHub Actions

- [ ] `.github/workflows/build.yml` created
- [ ] `.github/workflows/release.yml` created
- [ ] Workflows tested
- [ ] Build badges added to README (optional)

---

## First Release

```bash
# Create tag
git tag -a v1.0.0 -m "Release v1.0.0: Initial production release"

# Push tag
git push origin v1.0.0
```

On GitHub:
- [ ] Go to Releases
- [ ] Click "Create a new release"
- [ ] Choose tag: v1.0.0
- [ ] Release title: "v1.0.0 - Initial Release"
- [ ] Description added
- [ ] Binaries uploaded (optional)
- [ ] Release published

---

## Post-Deployment

- [ ] Clone repository to verify
- [ ] Test installation from GitHub
- [ ] Update any external links
- [ ] Share repository link
- [ ] Star your own repository 😊

---

## Maintenance

### Regular Updates
- [ ] Keep dependencies updated: `go get -u && go mod tidy`
- [ ] Update documentation as needed
- [ ] Respond to issues
- [ ] Review pull requests

### Version Releases
- [ ] Update version in code
- [ ] Update CHANGELOG.md
- [ ] Create new tag
- [ ] Create GitHub release
- [ ] Build and upload binaries

---

## Verification Commands

```bash
# Clone and test
git clone https://github.com/YOUR_USERNAME/watchup-agent.git
cd watchup-agent
go mod tidy
go build -o watchup-agent cmd/agent/main.go
./watchup-agent --help
```

---

## Common Issues

### Push Rejected
```bash
# Pull first
git pull origin main --rebase
git push origin main
```

### Authentication Failed
```bash
# Use personal access token
# GitHub Settings → Developer settings → Personal access tokens
# Use token as password when prompted
```

### Large Files
```bash
# Remove from git
git rm --cached large-file
echo "large-file" >> .gitignore
git commit -m "Remove large file"
git push
```

---

## Resources

- [GitHub Docs](https://docs.github.com)
- [Git Basics](https://git-scm.com/book/en/v2/Getting-Started-Git-Basics)
- [GitHub Actions](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)

---

## Checklist Complete? 🎉

If all items are checked, your project is successfully deployed to GitHub!

Next: [Install on VPS](../DEPLOYMENT.md#part-2-installing-on-vps)
