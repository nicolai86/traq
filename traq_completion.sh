# traq completion by Raphael Randschau <nicolai86@me.com>

function _traq() {
  COMPREPLY=()

  local cur="${COMP_WORDS[COMP_CWORD]}"
  local prev="${COMP_WORDS[COMP_CWORD-1]}"

  case "$prev" in
    # only traq was entered. suggest different traq arguments
    traq)
      COMPREPLY=( $(compgen -W "-d -w -p" -- $cur) )
      return 0
      ;;

    # project flag was entered. suggest known projects
    -p)
      local projects=$(ls $TRAQ_DATA_DIR | grep -v traq)
      COMPREPLY=( $(compgen -W "${projects}" -- $cur) )
      return 0
      ;;

    # date flag was entered. suggest the current date
    -d)
      local today=$(date "+%Y-%m-%d")
      COMPREPLY=( $(compgen -W "$today" -- $cur) )
      return 0
      ;;

    # week flag was entered. suggest the current week
    -w)
      local week_number=$(date "+%V")
      COMPREPLY=( $(compgen -W "$week_number" -- $cur) )
      return 0
      ;;

    *)
      ;;
  esac

  # secondary option. only available unless present.
  local options=""
  if [[ ! "$COMP_LINE" =~ "-d" ]] && [[ ! "$COMP_LINE" =~ "-w" ]]
  then
    options="$options -d"
  fi

  if [[ ! "$COMP_LINE" =~ "-w" ]] && [[ ! "$COMP_LINE" =~ "-d" ]]
  then
    options="$options -w"
  fi

  if [[ ! "$COMP_LINE" =~ "-p" ]]
  then
    options="$options -p"
  fi

  COMPREPLY=( $(compgen -W "$options" -- $cur) )
  return 0
}
complete -F _traq traq