#!/usr/bin/env bash

. $TRAQ_PATH/os_helper.sh

# convert date strings into unix timestamps
function date2timestamp() {
  if [ $(is_osx) ]
  then
    date -j -f "%a %b %d %T %z %Y" "$1" "+%s"
  else
    date -d "$1" "+%s"
  fi
}

# extract the date string from a traq line
function extract_date() {
  echo "$@" | cut -d';' -f1
}

# extract the tag string from a traq line
function extract_tag() {
  echo "$@" | cut -d';' -f2
}