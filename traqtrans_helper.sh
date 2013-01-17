#!/usr/bin/env bash

. $TRAQ_PATH/os_helper.sh

# convert date strings into unix timestamps
function date2timestamp() {
  if [[ $(is_osx) -eq 0 ]]
  then
    date -j -f "%a %b %d %T %z %Y" "$1" "+%s"
  else
    date -d "$1" "+%s"
  fi
}

# convert unix timestamp into date
function timestamp2date() {
  if [[ $(is_osx) -eq 0 ]]
  then
    date -j -f "%s" "$1" "+%a %b %d %T %z %Y"
  else
    date -d "@$1" "+%a %b %d %T %z %Y"
  fi
}

# convert date into %Y-%m-%d
function format_date() {
  if [[ $(is_osx) -eq 0 ]]
  then
    date -j -f "%a %b %d %T %z %Y" "$1" "+%Y-%m-%d"
  else
    date -d "$1" "+%Y-%m-%d"
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