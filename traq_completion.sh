# traq completion by Raphael Randschau <nicolai86@me.com>

function _traq() {
  COMPREPLY=()

  local cur="${COMP_WORDS[COMP_CWORD]}"
  local prev="${COMP_WORDS[COMP_CWORD-1]}"

  case "$prev" in
    # only traq was entered. suggest different traq arguments
    traq)
      COMPREPLY=( $(compgen -W "-d -e -p" -- $cur) )
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

    # month flag was entered. suggest the current month
    -m)
      local current_month=$(date "+%m")
      COMPREPLY=( $(compgen -W "$current_month" -- $cur) )
      return 0
      ;;

    # year flag was entered. suggest the current year
    -y)
      local current_year=$(date "+%Y")
      COMPREPLY=( $(compgen -W "$current_year" -- $cur) )
      return 0
      ;;

    *)
      ;;
  esac

  # secondary option. only available unless present.
  local options=""
  if [[ ! "$COMP_LINE" =~ "-d" ]] && [[ ! "$COMP_LINE" =~ "-e" ]]
  then
    options="$options -d"
  fi

  if [[ ! "$COMP_LINE" =~ "-e" ]] && [[ ! "$COMP_LINE" =~ "-d" ]]
  then
    options="$options -e"
  fi

  if [[ ! "$COMP_LINE" =~ "-p" ]]
  then
    options="$options -p"
  fi

  COMPREPLY=( $(compgen -W "$options" -- $cur) )
  return 0
}
complete -F _traq traq