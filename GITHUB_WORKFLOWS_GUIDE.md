# GitHub Actions Workflows Guide

This document explains what each GitHub Actions workflow file does and the commands I
learned or wrestled through to use

---

## My Workflow Files

I have 5 workflow files, each with a specific job:

### 1. `test.yml` - Test & Build
Code works and builds successfully. Every time you push code or create a pull request,
it is invoked

**Steps**
1. **Sets up a test environment** with Go and PostgreSQL database
2. **Downloads your code** from GitHub
3. **Installs dependencies** (like downloading libraries your code needs)
4. **Generates database code** using sqlc (converts SQL to Go code)
5. **Checks code quality** with `go vet` (finds potential bugs)
6. **Checks formatting** with `gofmt` (ensures consistent code style)
7. **Builds the application** (compiles Go code into an executable)
8. **Tests database connection** (makes sure the app can talk to the database)
9. **Runs tests** (if you have any test files)
10. **Saves the built application** as an "artifact" (a file other workflows can use)

**Commands I took from other projects**

- `go mod download` = Download all the libraries your project needs
- `go vet ./...` = Check all folders (`./...`) for potential bugs
- `gofmt -s -l .` = Check if code is formatted properly (`-s` = simplify, `-l` = list unformatted files)
- `timeout 10s ./telemetry-api` = Run the app for maximum 10 seconds then stop it

---

### 2. `security.yml` - Security Scan
Looks for security vulnerabilities in your code (also due to dependencies)

**Runs:** 
- Every push/pull request
- **Daily at 2 AM UTC** (automatic security check)

**Steps:**
1. **Scans your Go code** for security issues using Gosec
2. **Checks dependencies** for known vulnerabilities using govulncheck
3. **Uploads results** to GitHub's Security tab where you can see issues

**Commands I took from other projects**

- `gosec -fmt sarif -out gosec.sarif ./...` = Scan all code for security issues and save results in SARIF format (a standard format for security results)
- `govulncheck ./...` = Check if any libraries you use have known security problems
- **SARIF** = Security Analysis Results Interchange Format (a way to store security scan results)


---

### 3. `quality.yml` - Code Quality
Ensures your code follows best practices and measures test coverage, done every 
pull/push request

**Steps:**
1. **Linting** - Checks code style and finds potential issues
2. **Dependency check** - Makes sure `go.mod` file is up to date
3. **Code coverage** - Measures how much of your code is tested

**Commands I took from other projects**
- `golangci-lint` = A tool that runs multiple code quality checks at once
- `go mod tidy` = Clean up the `go.mod` file (remove unused dependencies, add missing ones)
- `go test -v -race -coverprofile=coverage.out` = Run tests with:
  - `-v` = verbose (show detailed output)
  - `-race` = detect race conditions (when multiple parts of code access same data simultaneously)
  - `-coverprofile` = measure how much code is tested
- **Race conditions** = Bugs that happen when multiple parts of your program try to use the same data at the same time

**Code coverage** Percentage of your code is tested.

---

### 4. `deploy.yml` - Deployment
Automatically deploys your application to staging/production servers

**Runs:**

- After successful tests and security scans
- Only on `main` branch (production) or `develop` branch (staging)


1. **Downloads the built application** from the test workflow
2. **Creates a GitHub release** (tags your code with a version number)
3. **Deploys to servers** (copies your app to the production server)
4. **Sends notifications** (optional - can notify Slack, email, etc.)

**Commands I took from other projects**
- `workflow_run` trigger = This workflow waits for other workflows to finish successfully
- `environment: production` = Uses GitHub's environment protection (can require approvals before deploying)
- `github.run_number` = Automatic version number that increases with each run
- `scp telemetry-api user@server:/opt/telemetry-api/` = Copy file to server using SSH

**Environments.** GitHub environments let you set up protection rules (like requiring someone to approve deployments to production).

---

### 5. `database.yml` - Database Management
Manages database backups, migrations, and seeding

**Runs:**
- **Daily at 3 AM UTC** (automatic backups)
- **Manually** when you trigger it from GitHub Actions tab

**What happens:**

1. **Backup** - Creates a copy of your database
2. **Migrate** - Updates database structure (adds new tables/columns)
3. **Seed** - Fills database with sample data

**Commands I took from other projects**
- `workflow_dispatch` = Allows manual triggering with options
- `pg_dump $DATABASE_URL > backup.sql` = Create a backup file of PostgreSQL database
- `cron: '0 3 * * *'` = Run at 3:00 AM every day (cron is a time scheduling format)
- **Cron format**: `minute hour day month day-of-week`

**Manual triggers:** You can go to GitHub Actions → Database Management → "Run workflow" and choose what to do (backup, migrate, or seed).

---

## Common Workflow Concepts

### **Artifacts**
Files saved by one workflow that other workflows can use.

### **Services**

External programs (like databases) that workflows need. GitHub provides these temporarily during the workflow run.

### **Cache**

Saves downloaded files between workflow runs to make them faster. Like keeping frequently used tools in an easily accessible toolbox.

### **Matrix Strategy**

Running the same workflow with different settings (like testing on different Go versions). We don't use this, but you might see it elsewhere.

---

## Common Issues & Solutions I faced fixing the workflows

### **"Workflow failed"**

1. Check the logs in GitHub Actions tab
2. Look for red marks to see which step failed
3. Common fixes:
   - Code formatting: Run `gofmt -w .` locally
   - Missing dependencies: Run `go mod tidy`
   - Test failures: Fix the failing tests

### **"Artifact not found"**
- The test workflow might have failed, so no artifact was created
- Check that the test workflow completed successfully first
- Artifacts are only available within the same workflow run by default. If the deploy workflow runs separately from the test workflow, it can't find the artifact.

**fix:** 
- Artifacts from different workflow runs need special handling
- Use `run-id` parameter to specify which workflow run to get artifacts from
- Or combine workflows so they run together

```yaml
- name: Download build artifact
  uses: actions/download-artifact@v4
  with:
    name: telemetry-api-${{ github.sha }}
    run-id: ${{ github.event.workflow_run.id }}
    github-token: ${{ secrets.GITHUB_TOKEN }}
```

### **"Linting errors" (errcheck, unused functions)**

Common linting issues and how to fix them:

**Unchecked errors (`errcheck`):**
```go
// Bad - ignoring error
data, _ := json.Marshal(obj)

// Good - checking error  
data, err := json.Marshal(obj)
if err != nil {
    return fmt.Errorf("failed to marshal: %w", err)
}
```

**JSON Unmarshal errors (also `errcheck`):**
```go
// Bad - ignoring unmarshal error
var data map[string]interface{}
json.Unmarshal(jsonBytes, &data)

// Good - checking unmarshal error with fallback
var data map[string]interface{}
if err := json.Unmarshal(jsonBytes, &data); err != nil {
    log.Printf("Failed to unmarshal: %v", err)
    data = make(map[string]interface{}) // provide fallback
}
```

**Unused functions (`unused`):**
- Remove functions that aren't called anywhere
- Or add `//nolint:unused` comment if you plan to use them later

**Printf format errors (`govet`):**
```go
// Bad - wrong format for string ID
log.Printf("Failed for item %d: %v", item.ID, err) // ID is string, not int

// Good - correct format for string ID  
log.Printf("Failed for item %s: %v", item.ID, err) // %s for strings
```

Common format specifiers:
- `%s` for strings
- `%d` for integers  
- `%v` for any value (safe fallback)

**fixes:**
- Always check error return values from functions
- Remove or comment unused code
- Use `continue` in loops to skip failed items instead of stopping everything

### **"Package not found" or "Module not found"**

This happens when package URLs change or repositories move:

- `github.com/securecodewarrior/gosec` → `github.com/securego/gosec` (gosec moved)
- Old package versions that no longer exist

**fix**
1. Search for the current package location
2. Update the `go install` command with the correct path
3. Check the official documentation for the latest install instructions

```yaml
# Old (broken)
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# New (working)
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

### **"Permission denied"**

- Usually means GitHub Secrets (like database passwords) aren't set up

### **"Version compatibility issues"**

Always use matching versions of upload and download artifact actions:
- `actions/upload-artifact@v4` with `actions/download-artifact@v4`
- `actions/upload-artifact@v3` with `actions/download-artifact@v3`

---

## Workflow Dependencies

Our workflows have dependencies:

```
Push to main branch
    ↓
test.yml runs ──→ Artifact created
    ↓
security.yml runs
    ↓
Both complete successfully
    ↓
deploy.yml runs ──→ Tries to download artifact
```

If any step fails, the chain breaks and deployment won't happen.