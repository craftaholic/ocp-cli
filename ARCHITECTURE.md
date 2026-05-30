# OCP Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          OCP - OpenCode Profile Switcher                     │
└─────────────────────────────────────────────────────────────────────────────┘

                                USER INTERACTION
                                       │
                    ┌──────────────────┼──────────────────┐
                    │                  │                  │
                    ▼                  ▼                  ▼
            ┌──────────────┐   ┌──────────────┐  ┌──────────────┐
            │  ocp use     │   │  ocp run     │  │  ocp list    │
            │  <profile>   │   │  <profile>   │  │  ocp status  │
            └──────┬───────┘   └──────┬───────┘  └──────┬───────┘
                   │                  │                  │
                   └──────────────────┼──────────────────┘
                                      │
                         ┌────────────▼────────────┐
                         │   OCP CLI (Cobra)       │
                         │   ─────────────         │
                         │   • Root Command        │
                         │   • Subcommands         │
                         │   • Flag Parsing        │
                         └────────────┬────────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    │                 │                 │
                    ▼                 ▼                 ▼
         ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
         │  Config Module  │ │   Env Module    │ │   Hook Module   │
         │  ─────────────  │ │   ──────────    │ │   ───────────   │
         │  • LoadConfig   │ │  • InjectVars   │ │  • ZshHook      │
         │  • SaveConfig   │ │  • ExpandPath   │ │  • BashHook     │
         │  • LoadProfile  │ │  • MaskValue    │ │  • FishHook     │
         │  • SaveProfile  │ │  • IsSensitive  │ │                 │
         │  • ListProfiles │ │                 │ │                 │
         │  • DeleteProfile│ │                 │ │                 │
         └────────┬────────┘ └────────┬────────┘ └────────┬────────┘
                  │                   │                    │
                  └───────────────────┼────────────────────┘
                                      │
                                      ▼
                    ┌──────────────────────────────────┐
                    │    FILE SYSTEM STORAGE           │
                    │    ~/.config/ocp/                │
                    │                                  │
                    │  ┌────────────────────────────┐  │
                    │  │ config.json                │  │
                    │  │ {                          │  │
                    │  │   "active": "work"         │  │
                    │  │ }                          │  │
                    │  └────────────────────────────┘  │
                    │                                  │
                    │  ┌────────────────────────────┐  │
                    │  │ profiles/                  │  │
                    │  │  ├─ work.json              │  │
                    │  │  │  {                      │  │
                    │  │  │    "name": "work",      │  │
                    │  │  │    "vars": {            │  │
                    │  │  │      "KEY": "val"       │  │
                    │  │  │    }                    │  │
                    │  │  │  }                      │  │
                    │  │  │                         │  │
                    │  │  └─ personal.json          │  │
                    │  └────────────────────────────┘  │
                    └──────────────────────────────────┘


                            USAGE WORKFLOWS
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  WORKFLOW 1: Shell Integration (Persistent Environment)                     │
│  ──────────────────────────────────────────────────────                     │
│                                                                              │
│   ~/.zshrc: eval "$(ocp init hook zsh)"                                     │
│                                                                              │
│   User Shell                    OCP Process              File System        │
│   ──────────                    ───────────              ───────────        │
│       │                              │                        │             │
│       │  $ ocp use work              │                        │             │
│       ├─────────────────────────────>│                        │             │
│       │                              │                        │             │
│       │                              │  Write active="work"   │             │
│       │                              ├───────────────────────>│             │
│       │                              │                        │             │
│       │  (hook intercepts)           │                        │             │
│       │  $ ocp status --export       │                        │             │
│       ├─────────────────────────────>│                        │             │
│       │                              │                        │             │
│       │                              │  Read work.json        │             │
│       │                              │<───────────────────────┤             │
│       │                              │                        │             │
│       │  export ANTHROPIC_API_KEY=.. │                        │             │
│       │<─────────────────────────────┤                        │             │
│       │                              │                        │             │
│   eval export ...                    │                        │             │
│   (vars now in shell)                │                        │             │
│       │                              │                        │             │
│       │  $ opencode                  │                        │             │
│       │  (uses work env vars)        │                        │             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘


┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  WORKFLOW 2: One-off Command (No Shell Modification)                        │
│  ───────────────────────────────────────────────────                        │
│                                                                              │
│   User Shell                    OCP Process              Target Command     │
│   ──────────                    ───────────              ──────────────     │
│       │                              │                        │             │
│       │  $ ocp run work -- opencode  │                        │             │
│       ├─────────────────────────────>│                        │             │
│       │                              │                        │             │
│       │                      Load work.json                   │             │
│       │                      Inject vars into env             │             │
│       │                      Find 'opencode' in PATH          │             │
│       │                              │                        │             │
│       │                      syscall.Exec()                   │             │
│       │                      (replaces process)               │             │
│       │                              ├───────────────────────>│             │
│       │                              │                        │             │
│       │                              X                    opencode          │
│       │                         (process replaced)        runs with         │
│       │                                                   work env vars     │
│       │                                                        │             │
│       │<───────────────────────────────────────────────────────┤             │
│       │                   (output from opencode)                            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘


                            DATA FLOW DIAGRAM
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│                           ┌─────────────────┐                               │
│                           │  User Command   │                               │
│                           └────────┬────────┘                               │
│                                    │                                        │
│                    ┌───────────────┼───────────────┐                        │
│                    │               │               │                        │
│                    ▼               ▼               ▼                        │
│              ┌─────────┐     ┌─────────┐    ┌─────────┐                    │
│              │   use   │     │   run   │    │  status │                    │
│              └────┬────┘     └────┬────┘    └────┬────┘                    │
│                   │               │              │                          │
│                   │               │              │                          │
│         ┌─────────▼────────┐      │              │                          │
│         │ 1. Check profile │      │              │                          │
│         │    exists        │      │              │                          │
│         └─────────┬────────┘      │              │                          │
│                   │               │              │                          │
│         ┌─────────▼────────┐      │              │                          │
│         │ 2. Write to      │      │              │                          │
│         │    config.json   │      │              │                          │
│         │    (atomic)      │      │              │                          │
│         └─────────┬────────┘      │              │                          │
│                   │               │              │                          │
│                   │               │              │                          │
│                   │      ┌────────▼────────┐     │                          │
│                   │      │ 1. Load profile │     │                          │
│                   │      │                 │     │                          │
│                   │      ├─────────────────┤     │                          │
│                   │      │ 2. Inject vars  │     │                          │
│                   │      │    + expand ~   │     │                          │
│                   │      │                 │     │                          │
│                   │      ├─────────────────┤     │                          │
│                   │      │ 3. Find command │     │                          │
│                   │      │    in PATH      │     │                          │
│                   │      │                 │     │                          │
│                   │      ├─────────────────┤     │                          │
│                   │      │ 4. syscall.Exec │     │                          │
│                   │      │    (no fork!)   │     │                          │
│                   │      └─────────────────┘     │                          │
│                   │                              │                          │
│                   │                    ┌─────────▼────────┐                 │
│                   │                    │ 1. Load config   │                 │
│                   │                    │                  │                 │
│                   │                    ├──────────────────┤                 │
│                   │                    │ 2. Load profile  │                 │
│                   │                    │                  │                 │
│                   │                    ├──────────────────┤                 │
│                   │                    │ 3. Mask secrets  │                 │
│                   │                    │    (if needed)   │                 │
│                   │                    │                  │                 │
│                   │                    ├──────────────────┤                 │
│                   │                    │ 4. Format output │                 │
│                   │                    │    (--export?)   │                 │
│                   │                    └──────────────────┘                 │
│                   │                                                          │
│                   └─────────────────┬────────────────────────┘              │
│                                     │                                        │
│                           ┌─────────▼────────┐                              │
│                           │   Output to      │                              │
│                           │   User/Shell     │                              │
│                           └──────────────────┘                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘


                          PROJECT STRUCTURE
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  ocp/                                                                        │
│  │                                                                           │
│  ├── main.go ───────────────────────> Entry point, calls cmd.Execute()     │
│  │                                                                           │
│  ├── cmd/ ──────────────────────────> CLI commands (Cobra)                 │
│  │   ├── root.go                       • Root command definition            │
│  │   ├── use.go                        • ocp use <profile>                  │
│  │   ├── run.go                        • ocp run <profile> -- <cmd>         │
│  │   ├── list.go                       • ocp list                           │
│  │   ├── status.go                     • ocp status                         │
│  │   ├── add.go                        • ocp add <profile>                  │
│  │   ├── edit.go                       • ocp edit <profile>                 │
│  │   ├── delete.go                     • ocp delete <profile>               │
│  │   └── hook.go                       • ocp init hook <shell>              │
│  │                                                                           │
│  ├── internal/ ─────────────────────> Core logic                           │
│  │   │                                                                       │
│  │   ├── config/ ──────────────────> Config & profile management           │
│  │   │   ├── config.go                 • LoadConfig, SaveConfig             │
│  │   │   │                             • GetConfigDir                       │
│  │   │   │                             • EnsureConfigDir                    │
│  │   │   │                                                                  │
│  │   │   └── profile.go                • LoadProfile, SaveProfile           │
│  │   │                                 • ListProfiles, DeleteProfile        │
│  │   │                                 • ProfileExists                      │
│  │   │                                                                       │
│  │   ├── env/ ─────────────────────> Environment handling                  │
│  │   │   └── inject.go                 • InjectProfileVars                  │
│  │   │                                 • ExpandPath (~/...)                 │
│  │   │                                 • IsSensitive                        │
│  │   │                                 • MaskValue                          │
│  │   │                                                                       │
│  │   └── hook/ ────────────────────> Shell integration                     │
│  │       ├── zsh.go                    • GetZshHook()                       │
│  │       ├── bash.go                   • GetBashHook()                      │
│  │       └── fish.go                   • GetFishHook()                      │
│  │                                                                           │
│  ├── go.mod ────────────────────────> Dependencies (just cobra)            │
│  ├── go.sum ────────────────────────> Checksums                            │
│  ├── Makefile ──────────────────────> Build targets                        │
│  ├── README.md ─────────────────────> Documentation                        │
│  └── LICENSE ───────────────────────> MIT License                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘


                     SECURITY & SAFETY FEATURES
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────┐      │
│  │  ATOMIC WRITES                                                    │      │
│  │  ──────────────                                                   │      │
│  │                                                                    │      │
│  │  Write Operation:                                                 │      │
│  │  1. Write to: config.json.tmp                                     │      │
│  │  2. Rename to: config.json  (atomic on Unix)                      │      │
│  │  3. No partial writes or corruption                               │      │
│  └──────────────────────────────────────────────────────────────────┘      │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────┐      │
│  │  SENSITIVE VALUE MASKING                                          │      │
│  │  ────────────────────────                                         │      │
│  │                                                                    │      │
│  │  Keywords: "key", "secret", "token", "password"                   │      │
│  │                                                                    │      │
│  │  Input:  ANTHROPIC_API_KEY=sk-ant-work-abc123xyz                 │      │
│  │  Output: ANTHROPIC_API_KEY=sk-ant-w...                            │      │
│  │          (first 8 chars + ...)                                    │      │
│  │                                                                    │      │
│  │  Note: --export flag bypasses masking (for shell hooks)           │      │
│  └──────────────────────────────────────────────────────────────────┘      │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────┐      │
│  │  ERROR HANDLING                                                   │      │
│  │  ───────────────                                                  │      │
│  │                                                                    │      │
│  │  • No panics - all errors returned and handled                    │      │
│  │  • Clear error messages for users                                 │      │
│  │  • Proper exit codes: 0=success, 1=user error                     │      │
│  │  • Graceful handling of missing config (first run)                │      │
│  └──────────────────────────────────────────────────────────────────┘      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```
