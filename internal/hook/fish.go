package hook

const FishHook = `function ocp
  if test "$argv[1]" = "use"
    command ocp $argv
    set -l profile_vars (command ocp status --export 2>/dev/null)
    if test -n "$profile_vars"
      eval $profile_vars
    end
  else
    command ocp $argv
  end
end

function ocp_prompt
  set -l active_profile (command ocp status --name-only 2>/dev/null)
  if test -n "$active_profile"
    echo "[ocp:$active_profile]"
  end
end
`

func GetFishHook() string {
	return FishHook
}
