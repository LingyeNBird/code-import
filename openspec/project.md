# Project Context

## Purpose
CNB Code Import is a batch code repository migration tool designed to facilitate large-scale repository transfers from multiple source platforms to CNB (a centralized code hosting platform).

**Key Goals:**
- Automate migration from 10+ VCS platforms (CODING, GitHub, GitLab, Gitee, Gitea, Alibaba Cloud Codeup, Huawei Cloud CodeArts Repo, CNB, Tencent Git, and generic platforms)
- Preserve repository structure using `<CNB root org>/<source repo path>` hierarchy
- Automatically handle large files (>256 MiB) by converting them to Git LFS objects
- Support CODING-specific features: map project display names to CNB sub-organization aliases and project descriptions to organization descriptions
- Track migration state to skip already-migrated repositories (via `successful.log`)
- Enable concurrent migrations (up to 10 simultaneous repositories) for efficiency

## Tech Stack
- **Language:** Go 1.23.0+ (with toolchain 1.23.4)
- **CLI Framework:** Cobra v1.8.0 - Command-line interface builder
- **Configuration:** Viper v1.18.2 - Multi-source configuration management (YAML, env vars, flags)
- **VCS Integration:**
  - `google/go-github` v66 - GitHub API v3 client
  - `xanzy/go-gitlab` v0.108.0 - GitLab API client
  - `huaweicloud/huaweicloud-sdk-go-v3` - Huawei Cloud CodeArts Repo SDK
- **Logging:** Uber Zap v1.21.0 - High-performance structured logging
- **Concurrency:** `golang.org/x/sync` v0.11.0 - Semaphore for rate limiting
- **HTTP/Auth:** `golang.org/x/oauth2` v0.27.0 - OAuth2 authentication
- **Serialization:** `gopkg.in/yaml.v3` v3.0.1 - YAML parsing
- **Build/Runtime:** Docker (multi-stage: Golang builder + Alpine runtime)
- **Version Control:** Git + Git LFS for large file handling

## Project Conventions

### Code Style
- **Language:** Go with idiomatic patterns; Chinese comments for internal code, English for global/exported APIs
- **Naming Conventions:**
  - Package-level variables: PascalCase or camelCase
  - Constants: UPPER_CASE with underscores
  - Exported types/functions: PascalCase
  - Private types/functions: camelCase
- **Error Handling:** Explicit error returns with context-rich messages; errors logged with structured fields
- **Logging:** Use `uber/zap` with appropriate levels (Info, Debug, Warn, Error); mask sensitive data (tokens, passwords) in logs
- **Comments:** Primarily Chinese language for internal documentation; use English for public-facing documentation
- **File Organization:** Group related functionality in packages (`cmd/` for CLI, `pkg/` for core logic)
- **No Formal Linters:** Project follows Go idioms informally; consider adding `.golangci.yml` for consistency

### Architecture Patterns
- **Plugin Architecture:** VCS platforms implement a common `VCS` interface (`pkg/vcs/interface.go`) enabling extensibility
- **Factory Pattern:** `NewVcs(platform)` returns platform-specific implementations (coding.go, github.go, gitlab.go, etc.)
- **Repository Pattern:** API clients in `pkg/api/*` abstract platform-specific REST/SDK calls
- **Concurrency Model:** Semaphore-based rate limiting (max 10 concurrent migrations) to avoid overwhelming APIs
- **Command Pattern:** Cobra CLI commands in `cmd/` orchestrate high-level operations
- **Wrapper Pattern:** `pkg/git/git.go` wraps Git CLI commands with retry logic and error handling
- **Configuration Hierarchy:** YAML file → Environment variables → CLI flags (increasing precedence)
- **Atomic Operations:** Use atomic counters for thread-safe statistics tracking (success/failure/skip counts)

### Testing Strategy
- **Framework:** Go standard `testing` package (no external frameworks like testify)
- **Test Coverage:** 8 test files covering critical paths:
  - Configuration parsing (`config_test.go`, `url_test.go`)
  - VCS platform implementations (`coding_test.go`, `gitee_test.go`)
  - Migration logic (`migrate_test.go`)
  - Git operations (`git_test.go`)
  - API clients (`api_test.go`)
  - Utilities (`util_test.go`)
- **Test Patterns:** Table-driven tests for multiple scenarios; mock external dependencies (APIs, Git commands)
- **CI Integration:** Tests run automatically in CNB CI/CD pipeline on pull requests
- **Future Improvements:** Consider adding integration tests for end-to-end migration scenarios

### Git Workflow
- **Branching Strategy:**
  - `main` branch: Protected, requires merge requests
  - Feature branches: Named descriptively (e.g., `vincent-patch-1`, `fix-gitlab-patch0`, `openspec`)
  - No direct commits to `main`
- **Commit Message Convention:** Chinese language predominant; follows pattern: `type: description`
  - **Types:** `feat:` (feature), `fix:` (bug fix), `docs:` (documentation), `merge:` (merge requests), `chore:` (maintenance)
  - **Example:** `feat: 优化手动流水线，设置CODING迁移按钮为默认选项`
- **Versioning:** Semantic versioning with Git tags (e.g., `v1.60.1`, `v1.61.0`)
- **CI/CD Triggers:**
  - **Push to main:** Docker build + push, auto-tagging, knowledge base build
  - **Pull requests:** AI-powered code review
  - **Tag push:** Changelog generation, release upload
- **Code Review:** All changes require merge request approval before merging

## Domain Context
**Code Repository Migration:**
- **Source Platforms:** Organizations/projects/repositories on various VCS platforms
- **Target Platform:** CNB (Chinese code hosting platform similar to GitLab/GitHub)
- **Migration Scope:**
  - Code (Git repositories with full history)
  - Releases (tags, release notes, assets) - optional
  - Large files (automatically converted to Git LFS if >256 MiB)
- **Organization Hierarchy:** CNB uses nested organizations (`root_org/sub_org/repo`); tool automatically creates sub-organizations based on source repository paths
- **CODING Platform Specifics:**
  - Maps CODING project display names to CNB sub-organization aliases
  - Maps CODING project descriptions to CNB sub-organization descriptions
  - Handles CODING-specific API quirks (team ID requirements, depot structure)
- **Migration State Tracking:** `successful.log` file records migrated repositories to enable incremental migrations and skip re-processing
- **Concurrency Control:** Limits concurrent operations to 10 to avoid API rate limiting and resource exhaustion
- **Credential Management:** Uses platform-specific authentication (tokens, OAuth2, AK/SK); credentials masked in logs

## Important Constraints
- **API Rate Limits:** Must respect platform-specific rate limits (e.g., GitHub: 5000 req/hour authenticated)
- **Large File Handling:** Files >256 MiB require Git LFS; standard Git push will fail
- **CNB Root Organization:** Must pre-exist before migration; tool cannot create root organization
- **Concurrent Migration Limit:** Maximum 10 concurrent repository migrations to avoid resource exhaustion
- **Docker Execution Timeout:** CNB web trigger has 10-hour timeout for long-running migrations
- **Network Connectivity:** Requires access to source platform APIs and CNB API; firewalls may block some platforms
- **Disk Space:** Clones entire repositories to `source_git_dir/` before pushing to CNB; ensure sufficient disk space
- **Authentication Tokens:** Require appropriate scopes:
  - Source platform: Read repositories, releases
  - CNB: Create organizations, repositories, push code
- **Platform Availability:** Tool behavior depends on source platform API availability and CNB service health
- **Git Version:** Requires Git 2.x+ and Git LFS installed in execution environment

## External Dependencies
**Required Services:**
1. **CNB Platform** (https://cnb.cool or custom instance)
   - API: Repository creation, organization management, code push
   - Authentication: Personal access tokens with `api`, `write_repository` scopes
2. **Source VCS Platforms:**
   - **CODING** (https://coding.net): OAuth2 or personal access token
   - **GitHub** (https://github.com): Personal access token with `repo`, `read:org` scopes
   - **GitLab** (https://gitlab.com or self-hosted): Personal access token with `api`, `read_repository` scopes
   - **Gitee** (https://gitee.com): Personal access token
   - **Gitea** (self-hosted): Personal access token
   - **Alibaba Cloud Codeup** (https://codeup.aliyun.com): AK/SK credentials
   - **Huawei Cloud CodeArts Repo** (https://devcloud.huaweicloud.com): AK/SK credentials
   - **Tencent Git (Gongfeng)**: Platform-specific authentication
   - **Generic platforms**: HTTP Basic Auth or token-based auth
3. **Git + Git LFS:** Local Git installation for repository operations
4. **Docker Registry** (for containerized deployment):
   - Images: `cnbcool/code-import:latest`, `cnbcool/code-import:v{VERSION}`

**Configuration Files:**
- `config.yaml`: Primary configuration (source/target platforms, migration options)
- `successful.log`: Tracks migrated repositories (auto-generated)
- `migrate.log`: Migration execution logs (auto-generated)
- `repo-path.txt`: Optional whitelist for selective repository migration

**Environment Variables (Docker mode):**
- All config keys can be set via `PLUGIN_*` prefixed variables (e.g., `PLUGIN_SOURCE_PLATFORM`, `PLUGIN_CNB_TOKEN`)

**CI/CD Integration:**
- `.cnb.yml`: CNB CI/CD pipeline configuration (Docker build, tagging, changelog generation)
- Uses CNB-provided Docker service and container registry
