# JWT Authentication in Golang (Gorilla Mux) 🔐
A secure JWT-based authentication system built in Golang using Gorilla Mux and GORM (PostgreSQL).

✅ User Registration & Login with hashed passwords (bcrypt).  
✅ JWT Token Generation & Verification.  
✅ Middleware for Protected Routes using JWT.  
✅ PostgreSQL Database Integration with GORM.  
✅ Automated Quality Gate (CI) on Pull Requests.  
✅ Production Deployment Separation (CD) on Approved Merge to `main`.  
✅ Branch Protection on `main` to prevent unreviewed or failing code.  

---

# 📌 Installation & Setup

### 1️⃣ Clone the Repository
```bash
git clone https://github.com/sangeeth200494/JWT-AUTH_Golang.git
cd JWT-AUTH_Golang
```

### 2️⃣ Install Dependencies
```bash
go mod tidy
```

### 3️⃣ Configure Environment Variables
Copy `.env.example` to `.env` and update the values:
```bash
cp .env.example .env
```
Example configuration:
```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=jwt_auth_db
DB_PORT=5432
JWT_SECRET=your_secret_key_here
```

### 4️⃣ Run Tests & Build Locally
```bash
# Run static analysis
go vet ./...

# Run unit tests
go test -v ./...

# Build binary
go build main.go
```

### 5️⃣ Run the Application
```bash
go run main.go
```

---

# 🛠 API Endpoints

| Method | Endpoint     | Description                  | Authorization     |
|--------|--------------|------------------------------|-------------------|
| `POST` | `/users`     | Register a new user          | ❌ No Token        |
| `POST` | `/login`     | Authenticate & retrieve JWT  | ❌ No Token        |
| `GET`  | `/Home`      | API health / welcome check   | ❌ No Token        |
| `GET`  | `/protected` | Access protected area        | ✅ Requires Bearer |

---

# 🛡️ CI/CD & Branch Protection Architecture

```
        Developer Feature Branch
                   │
                   ▼ (Open Pull Request to main)
     ┌───────────────────────────────┐
     │   CI Workflow: ci.yml         │
     │   - gofmt format check        │
     │   - go vet static analysis    │
     │   - go test -race & coverage  │
     │   - go build check            │
     └──────────────┬────────────────┘
                    │
                    ▼ [Status: Quality Gate & Validation]
     ┌───────────────────────────────┐
     │   GitHub Branch Protection    │
     │   - 1+ Approving Review       │
     │   - Passing CI Quality Gate   │
     │   - Up-to-date branch check   │
     │   - Enforced on Administrators│
     │   - No Direct Push to main    │
     └──────────────┬────────────────┘
                    │ (Approved & Merged)
                    ▼
          `main` Production Branch
                    │
                    ▼ (Push to main triggers CD)
     ┌───────────────────────────────┐
     │   CD Workflow: deploy.yml     │
     │   - Re-validate on main       │
     │   - Build release binary      │
     │   - Deploy to Production env  │
     └───────────────────────────────┘
```

### 1. Pull Request Validation (CI Only)
- Workflow file: `.github/workflows/ci.yml`
- Triggers on: Pull Requests to `main` and `dev`, and pushes to `dev`.
- Pull Requests **never** deploy to production. They only execute the quality gates:
  1. `gofmt -s -l .`: Ensures code formatting complies with Go standards.
  2. `go vet ./...`: Detects subtle bugs and suspicious constructs.
  3. `go test -v -race ./...`: Runs unit tests with race detection.
  4. `go build -v ./...`: Verifies compilation.

### 2. Production Deployment (CD Only)
- Workflow file: `.github/workflows/deploy.yml`
- Triggers strictly on: `push` to `main` (i.e. after an approved PR is merged) and manual `workflow_dispatch`.
- Does **not** run on Pull Requests.
- Uses GitHub Environment `production` with deployment secrets and artifact releases.

### 3. Branch Protection on `main`
Direct pushes to `main` are blocked. Developers must open a Pull Request, pass all CI checks, and obtain review approvals.

#### How to apply Branch Protection via Script
Set your GitHub Personal Access Token (PAT) with `repo` admin permissions, then run:

**PowerShell (Windows):**
```powershell
$env:GITHUB_TOKEN = "your_github_token"
.\scripts\setup-branch-protection.ps1
```

**Bash (Linux / macOS):**
```bash
export GITHUB_TOKEN="your_github_token"
chmod +x scripts/setup-branch-protection.sh
./scripts/setup-branch-protection.sh
```

#### How to apply manually in GitHub Web UI
1. Go to repository **Settings** → **Branches**.
2. Click **Add branch protection rule** (or edit rule for `main`).
3. Set **Branch name pattern** to `main`.
4. Enable **Require a pull request before merging**:
   - Check **Require approvals** (minimum 1).
   - Check **Dismiss stale pull request approvals when new commits are pushed**.
   - Check **Require approval of the most recent reviewable push**.
5. Enable **Require status checks to pass before merging**:
   - Check **Require branches to be up to date before merging**.
   - Search for and select: `Quality Gate & Validation` (from `ci.yml`).
6. Enable **Do not allow bypassing the above settings** (Enforce for administrators).
7. Enable **Require conversation resolution before merging**.
8. Click **Save changes**.

---

# 📜 License
This project is open-source and licensed under the MIT License.

# 📞 Contact
- **Author:** Sangeeth Jayaraj
- **Email:** sangeethvv554@gmail.com
- **GitHub:** [github.com/sangeeth200494](https://github.com/sangeeth200494)
