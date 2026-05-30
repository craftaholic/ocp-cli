# What Each OCP Command Does

## Quick Reference

| Command | What It Does | Files Accessed |
|---------|-------------|----------------|
| `ocp list` | Shows all profiles, marks active with * | Read: profiles/*.json, config.json |
| `ocp use <profile>` | Sets active profile | Write: config.json (atomic) |
| `ocp status` | Shows current profile vars (masked) | Read: config.json, profiles/<active>.json |
| `ocp status --export` | Outputs shell export commands | Read: config.json, profiles/<active>.json |
| `ocp run <profile> -- <cmd>` | Runs command with profile env | Read: profiles/<profile>.json, then exec |
| `ocp add <profile>` | Creates new profile, opens $EDITOR | Write: profiles/<profile>.json |
| `ocp edit <profile>` | Opens profile in $EDITOR | Read/Write: profiles/<profile>.json |
| `ocp delete <profile>` | Removes profile | Delete: profiles/<profile>.json |
| `ocp init hook <shell>` | Prints shell integration code | No files (just outputs code) |

---

## Detailed Breakdown

### 1. `ocp list`

**Purpose:** See all available profiles

**What happens internally:**
1. Opens directory `~/.config/ocp/profiles/`
2. Finds all `.json` files
3. Reads `~/.config/ocp/config.json` to see active profile
4. Sorts alphabetically
5. Prints list with `*` marking active profile

**Example:**
```bash
$ ocp list
  personal
* work
  staging
```

---

### 2. `ocp use <profile>`

**Purpose:** Switch to a different profile

**What happens internally:**
1. Check if `~/.config/ocp/profiles/<profile>.json` exists
2. If not found → error "profile does not exist"
3. Load current `config.json` (or create empty struct)
4. Set `active` field to `<profile>`
5. Marshal to JSON with indentation
6. Write to `~/.config/ocp/config.json.tmp`
7. Rename tmp → config.json (atomic operation)
8. Print success message

**File changes:**
```json
// ~/.config/ocp/config.json BEFORE
{
  "active": "personal"
}

// AFTER running: ocp use work
{
  "active": "work"
}
```

**With shell hook installed:**
- After writing config, shell hook runs `ocp status --export`
- Evaluates the output to inject vars into current shell
- Your shell now has the new profile's environment variables

---

### 3. `ocp status`

**Purpose:** View current profile and its variables

**What happens internally:**
1. Read `~/.config/ocp/config.json`
2. If no active profile → print "No active profile" and exit
3. Load `~/.config/ocp/profiles/<active>.json`
4. For each variable in profile:
   - Check if key contains: "key", "secret", "token", "password" (case-insensitive)
   - If sensitive: mask value (show first 8 chars + "...")
   - If normal: show full value
5. Sort variables alphabetically
6. Print formatted output

**Example output:**
```
Active profile: work

Environment variables:
  ANTHROPIC_API_KEY=sk-ant-w...      ← MASKED (contains "key")
  MY_SECRET_TOKEN=mytoken1...         ← MASKED (contains "token")
  NORMAL_VAR=full_value_here          ← NOT MASKED
  OPENCODE_CONFIG_DIR=~/.config/ocp   ← NOT MASKED
```

---

### 4. `ocp status --export`

**Purpose:** Output variables in shell-exportable format (used by shell hooks)

**What happens internally:**
1. Same as `ocp status` but:
   - **NO masking** (exports full values)
   - **Expands tildes** (`~/path` → `/home/user/path`)
   - Outputs in format: `export KEY="value"`

**Example output:**
```bash
export ANTHROPIC_API_KEY="sk-ant-work-key-123"
export OPENCODE_CONFIG_DIR="/home/user/.config/opencode-work"
export MY_VAR="some_value"
```

**This output is meant to be `eval`'d:**
```bash
eval "$(ocp status --export)"
# Now your shell has all those environment variables set
```

---

### 5. `ocp run <profile> -- <cmd>`

**Purpose:** Run a command with profile's environment variables (one-off, doesn't affect shell)

**What happens internally:**
1. Parse arguments: profile name and command
2. Load `~/.config/ocp/profiles/<profile>.json`
3. Get all current environment variables (`os.Environ()`)
4. Create a map from current env
5. Overlay profile vars on top (profile vars override existing)
6. Expand any `~` in paths
7. Find command in PATH using `exec.LookPath()`
8. Call `syscall.Exec(cmdPath, args, mergedEnv)`
   - **This REPLACES the ocp process!**
   - The command becomes PID of ocp (no subprocess)
   - Zero overhead

**Example:**
```bash
$ ocp run work -- opencode

# What happens:
# 1. Loads work.json vars
# 2. Current env: PATH=/usr/bin:..., HOME=/home/user, ...
# 3. Merges: ANTHROPIC_API_KEY=sk-ant-work-key, ...
# 4. Finds 'opencode' at /usr/local/bin/opencode
# 5. Calls syscall.Exec("/usr/local/bin/opencode", ["opencode"], mergedEnv)
# 6. ocp process BECOMES opencode process
```

**Default command:**
If you don't specify a command, it defaults to `opencode`:
```bash
ocp run work          # Same as: ocp run work -- opencode
```

---

### 6. `ocp add <profile>`

**Purpose:** Create a new profile

**What happens internally:**
1. Check if profile already exists → error if yes
2. Create empty profile struct:
   ```json
   {
     "name": "<profile>",
     "vars": {}
   }
   ```
3. Marshal to JSON with indentation
4. Write to `~/.config/ocp/profiles/<profile>.json.tmp`
5. Rename tmp → .json (atomic)
6. Print "Created profile '<profile>'"
7. Get `$EDITOR` environment variable (default: `vi`)
8. Run: `$EDITOR ~/.config/ocp/profiles/<profile>.json`
9. Wait for editor to exit

**User then edits the file:**
```json
{
  "name": "newprofile",
  "vars": {
    "ANTHROPIC_API_KEY": "sk-ant-...",
    "OPENCODE_CONFIG_DIR": "~/.config/opencode-new"
  }
}
```

---

### 7. `ocp edit <profile>`

**Purpose:** Modify an existing profile

**What happens internally:**
1. Check if profile exists → error if not
2. Get `$EDITOR` (default: `vi`)
3. Run: `$EDITOR ~/.config/ocp/profiles/<profile>.json`
4. Wait for editor to exit
5. No validation (user can break the JSON if they want)

---

### 8. `ocp delete <profile>`

**Purpose:** Remove a profile

**What happens internally:**
1. Load `config.json`
2. If deleting active profile:
   - Clear active field: `{"active": ""}`
   - Save config.json (atomic write)
3. Delete `~/.config/ocp/profiles/<profile>.json`
4. Print "Deleted profile '<profile>'"

---

### 9. `ocp init hook <shell>`

**Purpose:** Generate shell integration code

**What happens internally:**
1. Check shell argument (zsh, bash, or fish)
2. If invalid → error "unsupported shell"
3. Load pre-written hook template from internal/hook/<shell>.go
4. Print to stdout (no file operations)

**The hook code wraps the `ocp` command:**

For zsh/bash:
```bash
ocp() {
  if [[ "$1" == "use" ]]; then
    command ocp "$@"                           # Run real ocp
    local profile_vars
    profile_vars=$(command ocp status --export)  # Get export commands
    if [[ -n "$profile_vars" ]]; then
      eval "$profile_vars"                     # Inject into shell
    fi
  else
    command ocp "$@"                           # Pass through other commands
  fi
}
```

**User adds this to ~/.zshrc:**
```bash
eval "$(ocp init hook zsh)"
```

Now when user runs `ocp use work`, the hook:
1. Calls real ocp binary
2. Gets export commands
3. Evaluates them to inject vars into current shell

---

## Special Features

### Atomic Writes
All file writes use this pattern:
```go
1. Write to: file.tmp
2. Rename to: file (atomic on Unix)
3. If error, remove tmp
```

This prevents corruption if process is interrupted.

### Tilde Expansion
Paths like `~/config` are expanded to full paths:
```
~/.config/ocp → /home/user/.config/ocp
```

### Sensitive Value Masking
Variables with these keywords are masked in `ocp status`:
- key
- secret
- token
- password

**Logic:**
```go
if strings.Contains(strings.ToLower(key), "key") {
    value = value[:8] + "..."
}
```

### Process Replacement (syscall.Exec)
Traditional approach (fork + exec):
```
Shell (PID 100)
  └─ ocp (PID 101)
       └─ opencode (PID 102)  ← New process created
```

OCP's approach (exec only):
```
Shell (PID 100)
  └─ ocp (PID 101) → opencode (PID 101)  ← Same PID!
```

The ocp process is **replaced** by opencode, not creating a subprocess.

---

## Summary Table

| Command | Reads Files | Writes Files | Spawns Editor | Exec's Command |
|---------|-------------|--------------|---------------|----------------|
| list | ✓ | | | |
| use | ✓ | ✓ | | |
| status | ✓ | | | |
| run | ✓ | | | ✓ |
| add | | ✓ | ✓ | |
| edit | | | ✓ | |
| delete | ✓ | ✓ | | |
| init hook | | | | |
