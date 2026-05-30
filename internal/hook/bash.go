package hook

const BashHook = `ocp() {
  if [[ "$1" == "use" ]]; then
    command ocp "$@"
    local profile_vars
    profile_vars=$(command ocp status --export 2>/dev/null)
    if [[ -n "$profile_vars" ]]; then
      eval "$profile_vars"
    fi
  else
    command ocp "$@"
  fi
}

ocp_prompt() {
  local active_profile
  active_profile=$(command ocp status --name-only 2>/dev/null)
  if [[ -n "$active_profile" ]]; then
    echo "[ocp:$active_profile]"
  fi
}
`

func GetBashHook() string {
	return BashHook
}
