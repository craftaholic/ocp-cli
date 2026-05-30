# ocp - OpenCode Profile Switcher

A CLI tool for managing multiple named profiles for [opencode](https://opencode.ai) and [claude code](https://claude.ai/code). Each profile contains a set of environment variables (API keys, config directories, model preferences). Switch between profiles seamlessly in your shell without manually exporting variables.

## Features

- **Auto-Migration**: Automatically detects and migrates existing opencode configuration on first run
- **Multiple Profiles**: Create and manage separate profiles for work, personal, or different projects
- **Directory-Based**: Each profile is a self-contained directory with all config files
- **Symlink Switching**: Uses symlinks for instant profile switching - no environment variables needed
- **Seamless Switching**: Change active profile with a single command
- **Shell Integration**: Automatic environment variable loading in your shell
- **Secure**: Masks sensitive values (keys, tokens, secrets) when displaying profiles
- **Zero Overhead**: `ocp run` uses syscall.Exec for direct process replacement
- **Cross-Shell Support**: Works with zsh, bash, and fish

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/craftaholic/ocp-cli.git
cd ocp-cli

# Install using Make
make install

# Or build manually
go build -o ocp .
sudo mv ocp /usr/local/bin/
```

### Using go install

```bash
go install github.com/craftaholic/ocp-cli@latest
```

## Quick Start

### First Run: Auto-Migration

If you already have an existing `~/.config/opencode` configuration, ocp will automatically migrate it on first run:

```bash
$ ocp list

🔍 Detected existing opencode configuration at ~/.config/opencode
📦 Migrating to ocp default profile...
✅ Migration complete!
   - Moved ~/.config/opencode to ~/.config/ocp/profiles/default/
   - Created symlink: ~/.config/opencode -> ~/.config/ocp/profiles/default/
   - Set 'default' as active profile

Your existing configuration is now in the 'default' profile.
You can create additional profiles with: ocp add <profile>

* default
```

**All your existing files are preserved:**
- Configuration (opencode.json)
- Conversation history
- Any other files

You can continue using opencode immediately - everything works as before!

See [MIGRATION.md](MIGRATION.md) for details.

### 1. Create Your First Profile

```bash
# Create a personal profile
ocp add personal
```

This opens your `$EDITOR` with a JSON template:

```json
{
  "name": "personal",
  "vars": {
    "ANTHROPIC_API_KEY": "sk-ant-api03-...",
    "OPENCODE_CONFIG_DIR": "~/.config/opencode",
    "CLAUDE_CONFIG_DIR": "~/.config/claude"
  }
}
```

### 2. Create Additional Profiles

```bash
# Create a work profile
ocp add work
```

Edit it with your work credentials:

```json
{
  "name": "work",
  "vars": {
    "ANTHROPIC_API_KEY": "sk-ant-api03-work-key-...",
    "OPENCODE_CONFIG_DIR": "~/.config/opencode-work",
    "CLAUDE_CONFIG_DIR": "~/.config/claude-work",
    "ANTHROPIC_MODEL": "claude-opus-4-20250514"
  }
}
```

### 3. Set Up Shell Integration

Add the following to your shell's RC file:

**For zsh (~/.zshrc):**
```bash
eval "$(ocp init hook zsh)"
```

**For bash (~/.bashrc):**
```bash
eval "$(ocp init hook bash)"
```

**For fish (~/.config/fish/config.fish):**
```fish
ocp init hook fish | source
```

Reload your shell:
```bash
exec $SHELL
```

### 4. Switch Profiles

```bash
# Switch to work profile
ocp use work

# Switch to personal profile
ocp use personal

# List all profiles (active marked with *)
ocp list
```

## Commands

### `ocp use <profile>`
Set the active profile. When shell integration is enabled, environment variables are automatically exported to your current shell.

```bash
ocp use work
# Switched to profile 'work'
```

### `ocp run <profile> [-- <cmd>]`
Run a command with profile environment variables injected. Defaults to running `opencode` if no command is specified.

```bash
# Run opencode with work profile
ocp run work

# Run specific command with personal profile
ocp run personal -- opencode

# Run claude with work profile
ocp run work -- claude --version
```

### `ocp list`
List all profiles. The active profile is marked with `*`.

```bash
ocp list
# * work
#   personal
```

### `ocp status`
Show the active profile and its environment variables. Sensitive values are masked.

```bash
ocp status
# Active profile: work
# 
# Environment variables:
#   ANTHROPIC_API_KEY=sk-ant-a...
#   CLAUDE_CONFIG_DIR=/home/user/.config/claude-work
#   OPENCODE_CONFIG_DIR=/home/user/.config/opencode-work
```

**Flags:**
- `--export`: Output in shell export format (for shell hooks)
- `--name-only`: Output only the active profile name

### `ocp add <profile>`
Create a new profile and open it in `$EDITOR`.

```bash
ocp add staging
```

### `ocp edit <profile>`
Edit an existing profile in `$EDITOR`.

```bash
ocp edit work
```

### `ocp delete <profile>`
Delete a profile. If it's the active profile, the active selection is cleared.

```bash
ocp delete old-profile
```

### `ocp init hook <shell>`
Print shell hook code for zsh, bash, or fish.

```bash
ocp init hook zsh
```

## Configuration

### Config Location

All configuration is stored in `~/.config/ocp/`:

```
~/.config/ocp/
├── config.json          # Active profile selection
└── profiles/
    ├── personal.json    # Personal profile
    └── work.json        # Work profile
```

### Profile Schema

Each profile is a JSON file with the following structure:

```json
{
  "name": "profile-name",
  "vars": {
    "ENV_VAR_NAME": "value",
    "ANOTHER_VAR": "~/path/with/tilde/expansion"
  }
}
```

### Path Expansion

Paths starting with `~/` are automatically expanded to your home directory when injected into the environment.

### Sensitive Value Masking

Variables containing these keywords in their names are masked in `ocp status` output:
- key
- secret
- token
- password

Masked values show only the first 8 characters followed by `...`

## Example Profiles

### OpenCode Personal Profile

```json
{
  "name": "personal",
  "vars": {
    "ANTHROPIC_API_KEY": "sk-ant-api03-personal-...",
    "OPENCODE_CONFIG_DIR": "~/.config/opencode",
    "ANTHROPIC_MODEL": "claude-sonnet-4-20250514"
  }
}
```

### Work Profile with Custom Settings

```json
{
  "name": "work",
  "vars": {
    "ANTHROPIC_API_KEY": "sk-ant-api03-work-...",
    "OPENCODE_CONFIG_DIR": "~/.config/opencode-work",
    "CLAUDE_CONFIG_DIR": "~/.config/claude-work",
    "ANTHROPIC_MODEL": "claude-opus-4-20250514",
    "ANTHROPIC_BASE_URL": "https://api.company-proxy.com"
  }
}
```

### Testing Profile

```json
{
  "name": "testing",
  "vars": {
    "ANTHROPIC_API_KEY": "sk-ant-api03-test-...",
    "OPENCODE_CONFIG_DIR": "~/.config/opencode-test",
    "ANTHROPIC_MODEL": "claude-sonnet-4-20250514",
    "OPENCODE_LOG_LEVEL": "debug"
  }
}
```

## Using Different opencode.json Configs Per Profile

Each profile can point to a different config directory, allowing you to have completely separate `opencode.json` configurations for each profile.

### Setup

**1. Create separate config directories:**

```bash
mkdir -p ~/.config/opencode-work
mkdir -p ~/.config/opencode-personal
```

**2. Create different opencode.json in each directory:**

`~/.config/opencode-work/opencode.json`:
```json
{
  "model": "claude-opus-4-20250514",
  "temperature": 0.2,
  "max_tokens": 8000,
  "system_prompt": "You are a professional coding assistant."
}
```

`~/.config/opencode-personal/opencode.json`:
```json
{
  "model": "claude-sonnet-4-20250514",
  "temperature": 0.7,
  "max_tokens": 4000,
  "system_prompt": "You are a helpful coding assistant."
}
```

**3. Create ocp profiles pointing to these directories:**

```bash
# Work profile
cat > ~/.config/ocp/profiles/work.json << EOF
{
  "name": "work",
  "vars": {
    "ANTHROPIC_API_KEY": "sk-ant-work-key",
    "OPENCODE_CONFIG_DIR": "~/.config/opencode-work"
  }
}
EOF

# Personal profile
cat > ~/.config/ocp/profiles/personal.json << EOF
{
  "name": "personal",
  "vars": {
    "ANTHROPIC_API_KEY": "sk-ant-personal-key",
    "OPENCODE_CONFIG_DIR": "~/.config/opencode-personal"
  }
}
EOF
```

**4. Switch between profiles:**

```bash
$ ocp use work
# opencode now uses ~/.config/opencode-work/opencode.json

$ ocp use personal
# opencode now uses ~/.config/opencode-personal/opencode.json
```

### Result

Each profile now has:
- Different API key
- Different config directory
- Different opencode.json settings (model, temperature, prompts, etc.)
- Separate conversation history
- Separate preferences

This gives you **complete isolation** between work and personal environments!

## Shell Integration Details

### How It Works

When you run `ocp use <profile>`:

1. The profile name is saved to `~/.config/ocp/config.json`
2. The shell hook intercepts the `ocp use` command
3. After setting the profile, it runs `ocp status --export`
4. The exported variables are evaluated into your current shell

### Prompt Integration

The shell hook provides an `ocp_prompt()` function you can use in your prompt:

**For zsh:**
```bash
# Add to your prompt in ~/.zshrc
PROMPT='$(ocp_prompt) %~ %# '
```

**For bash:**
```bash
# Add to PS1 in ~/.bashrc
PS1='$(ocp_prompt) \w \$ '
```

**For fish:**
```fish
# Add to your fish_prompt function
function fish_prompt
    echo (ocp_prompt) (prompt_pwd) '> '
end
```

This displays `[ocp:profile-name]` in your prompt when a profile is active.

## Workflows

### Quick Profile Switch for Different Projects

```bash
# Working on personal project
ocp use personal
opencode

# Switch to work project
ocp use work
opencode
```

### One-off Commands Without Switching

```bash
# Keep personal as active, but run one command with work profile
ocp run work -- opencode --version
```

### Testing New Configuration

```bash
# Create test profile
ocp add test

# Edit with test settings
ocp edit test

# Test without affecting your active profile
ocp run test -- opencode
```

## Troubleshooting

### Shell integration not working

Make sure you've added the hook to your RC file and reloaded your shell:

```bash
exec $SHELL
```

### Editor not opening for `ocp add` or `ocp edit`

Set your `EDITOR` environment variable:

```bash
export EDITOR=vim  # or nano, emacs, code, etc.
```

### Command not found when using `ocp run`

Ensure the command is in your `PATH`. The `ocp run` command searches your PATH for executables.

## Exit Codes

- `0`: Success
- `1`: User error (invalid arguments, profile not found, etc.)
- `2`: Internal error

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.
