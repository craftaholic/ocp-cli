# OCP Flow Diagrams

## Quick Visual: How Profile Switching Works

```
┌─────────────────────────────────────────────────────────────────┐
│                    BEFORE: No Profile Active                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Your Shell                                                      │
│  ──────────                                                      │
│  ANTHROPIC_API_KEY=<not set>                                    │
│  OPENCODE_CONFIG_DIR=<not set>                                  │
│                                                                  │
│  ~/.config/ocp/config.json:                                     │
│  { "active": "" }                                               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ $ ocp use work
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     AFTER: Work Profile Active                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Your Shell (with shell hook)                                   │
│  ────────────────────────────                                   │
│  ANTHROPIC_API_KEY=sk-ant-work-abc123xyz                        │
│  OPENCODE_CONFIG_DIR=/home/user/.config/opencode-work           │
│  ANTHROPIC_MODEL=claude-opus-4-20250514                         │
│                                                                  │
│  ~/.config/ocp/config.json:                                     │
│  { "active": "work" }                                           │
│                                                                  │
│  Now when you run 'opencode', it uses these vars!               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ $ ocp use personal
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   AFTER: Personal Profile Active                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Your Shell (vars updated instantly)                            │
│  ────────────────────────────────                               │
│  ANTHROPIC_API_KEY=sk-ant-personal-xyz789abc                    │
│  OPENCODE_CONFIG_DIR=/home/user/.config/opencode-personal       │
│  ANTHROPIC_MODEL=claude-sonnet-4-20250514                       │
│                                                                  │
│  ~/.config/ocp/config.json:                                     │
│  { "active": "personal" }                                       │
│                                                                  │
│  Now 'opencode' uses your personal settings!                    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Command Flow Charts

### 1. ocp use <profile>

```
     START
       │
       ▼
┌─────────────┐
│ User types: │
│ ocp use work│
└──────┬──────┘
       │
       ▼
┌──────────────────┐     NO      ┌──────────────────┐
│ Does 'work'      ├────────────>│ Show error:      │
│ profile exist?   │             │ Profile not found│──> EXIT(1)
└──────┬───────────┘             └──────────────────┘
       │ YES
       ▼
┌──────────────────┐
│ Load config.json │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Set active="work"│
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Write to         │
│ config.json.tmp  │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Rename tmp to    │
│ config.json      │
│ (atomic)         │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Print: "Switched │
│ to profile 'work'│
└──────┬───────────┘
       │
       ▼
    EXIT(0)
       
       
   [IF SHELL HOOK IS INSTALLED]
       │
       ▼
┌──────────────────┐
│ Hook intercepts  │
│ and runs:        │
│ ocp status       │
│ --export         │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Load work.json   │
│ Get all vars     │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Output:          │
│ export KEY=val   │
│ export KEY2=val2 │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Shell evaluates  │
│ export commands  │
│ Vars now in env! │
└──────────────────┘
```

### 2. ocp run <profile> -- <cmd>

```
       START
         │
         ▼
  ┌──────────────┐
  │ User types:  │
  │ ocp run work │
  │ -- opencode  │
  └──────┬───────┘
         │
         ▼
  ┌──────────────────┐
  │ Parse arguments: │
  │ profile='work'   │
  │ cmd='opencode'   │
  └──────┬───────────┘
         │
         ▼
  ┌──────────────────┐     NO      ┌──────────────┐
  │ Load work.json   ├────────────>│ Show error   │──> EXIT(1)
  │                  │             └──────────────┘
  └──────┬───────────┘
         │ YES
         ▼
  ┌──────────────────┐
  │ Get profile vars:│
  │ API_KEY=xxx      │
  │ CONFIG_DIR=yyy   │
  └──────┬───────────┘
         │
         ▼
  ┌──────────────────┐
  │ Merge with       │
  │ os.Environ()     │
  │ Expand ~ paths   │
  └──────┬───────────┘
         │
         ▼
  ┌──────────────────┐     NO      ┌──────────────┐
  │ Find 'opencode'  ├────────────>│ Command not  │──> EXIT(1)
  │ in PATH          │             │ found error  │
  └──────┬───────────┘             └──────────────┘
         │ YES
         │ (e.g., /usr/local/bin/opencode)
         ▼
  ┌──────────────────┐
  │ syscall.Exec(    │
  │   path,          │
  │   args,          │
  │   merged_env     │
  │ )                │
  └──────┬───────────┘
         │
         ▼
  ┌──────────────────┐
  │ OCP PROCESS      │
  │ REPLACED!        │
  │ Now running      │
  │ 'opencode' with  │
  │ injected env vars│
  └──────────────────┘
         │
         ▼
    (opencode runs...)
```

### 3. ocp status

```
      START
        │
        ▼
 ┌──────────────┐
 │ User types:  │
 │ ocp status   │
 └──────┬───────┘
        │
        ▼
 ┌──────────────────┐
 │ Load config.json │
 └──────┬───────────┘
        │
        ▼
 ┌──────────────────┐     NO      ┌──────────────┐
 │ Is there an      ├────────────>│ Print: No    │──> EXIT(0)
 │ active profile?  │             │ active profile│
 └──────┬───────────┘             └──────────────┘
        │ YES
        ▼
 ┌──────────────────┐
 │ Load active      │
 │ profile.json     │
 └──────┬───────────┘
        │
        ▼
 ┌──────────────────┐
 │ For each var:    │
 │ Check if         │
 │ sensitive?       │
 └──────┬───────────┘
        │
        ├─────────────────┐
        │                 │
        ▼                 ▼
 ┌──────────────┐  ┌──────────────┐
 │ Sensitive    │  │ Normal var   │
 │ (has "key",  │  │ Show full    │
 │  "secret",   │  │ value        │
 │  "token",    │  │              │
 │  "password") │  │ FOO=bar      │
 │              │  └──────────────┘
 │ Mask value:  │
 │ KEY=first8...│
 └──────────────┘
        │
        └─────────────────┘
        │
        ▼
 ┌──────────────────┐
 │ Print formatted  │
 │ output           │
 └──────┬───────────┘
        │
        ▼
     EXIT(0)
     
     
 [IF --export FLAG]
        │
        ▼
 ┌──────────────────┐
 │ NO masking!      │
 │ Output format:   │
 │ export KEY="val" │
 │ export KEY2="v2" │
 └──────────────────┘
```

## Real-World Example Flow

```
┌────────────────────────────────────────────────────────────────────────┐
│                  TYPICAL DAY WITH OCP                                  │
└────────────────────────────────────────────────────────────────────────┘

    Morning - Start Work
    ────────────────────
         User                    OCP                     opencode
          │                       │                          │
          │  ocp use work        │                          │
          ├─────────────────────>│                          │
          │                       │                          │
          │                   [loads work                    │
          │                    profile vars]                 │
          │                       │                          │
          │  opencode             │                          │
          ├──────────────────────────────────────────────────>
          │                       │                          │
          │                       │      [uses work API key] │
          │                       │      [work config dir]   │
          │<──────────────────────────────────────────────────
          │     (work session)                               │
          │                                                   │
          
          
    Lunch Break - Personal Project
    ───────────────────────────────
          │                       │                          │
          │  ocp use personal    │                          │
          ├─────────────────────>│                          │
          │                       │                          │
          │                   [loads personal                │
          │                    profile vars]                 │
          │                       │                          │
          │  opencode             │                          │
          ├──────────────────────────────────────────────────>
          │                       │                          │
          │                       │  [uses personal API key] │
          │                       │  [personal config dir]   │
          │<──────────────────────────────────────────────────
          │   (personal session)                             │
          │                                                   │
          
          
    Afternoon - Back to Work
    ────────────────────────
          │                       │                          │
          │  ocp use work        │                          │
          ├─────────────────────>│                          │
          │                       │                          │
          │  opencode             │                          │
          ├──────────────────────────────────────────────────>
          │                       │                          │
          │                       │      [work env again!]   │
          │<──────────────────────────────────────────────────
          │                                                   │
          
          
    Quick Test - One-off Command
    ─────────────────────────────
          │                       │                          │
          │  ocp run staging     │                          │
          │  -- opencode test    │                          │
          ├─────────────────────>│                          │
          │                       │                          │
          │                  [OCP process                    │
          │                   replaced by                    │
          │                   opencode with                  │
          │                   staging vars]                  │
          │                                                   │
          │<──────────────────────────────────────────────────
          │  (test output)                                   │
          │                                                   │
          │  Still using 'work' profile in shell!            │
          │                                                   │
```

## The Magic of syscall.Exec

```
┌─────────────────────────────────────────────────────────────────────┐
│                    TRADITIONAL APPROACH                             │
│                    (NOT what ocp does)                              │
└─────────────────────────────────────────────────────────────────────┘

    Your Shell (PID 100)
         │
         │ ocp run work -- opencode
         │
         ▼
    OCP Process (PID 101)
         │
         │ fork()
         ▼
    Child Process (PID 102)
         │
         │ exec opencode
         ▼
    opencode (PID 102)
    
    Problem: OCP process (101) stays alive, wasting memory
    

┌─────────────────────────────────────────────────────────────────────┐
│                      OCP'S APPROACH                                 │
│                  (using syscall.Exec)                               │
└─────────────────────────────────────────────────────────────────────┘

    Your Shell (PID 100)
         │
         │ ocp run work -- opencode
         │
         ▼
    OCP Process (PID 101)
         │
         │ syscall.Exec(opencode, args, env)
         │ ↓
         │ [Process 101 REPLACED IN-PLACE]
         ▼
    opencode (PID 101)  ← Same PID! No extra process!
    
    Benefit: Zero overhead, clean process tree
```
