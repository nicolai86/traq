#!/usr/bin/env bash

# convert date strings into unix timestamps
function date2timestamp() {
  date -j -f "%a %b %d %T %z %Y" "$1" "+%s"
}

# extract the date string from a traq line
function extract_date() {
  echo "$@" | cut -d';' -f1
}

# extract the tag string from a traq line
function extract_tag() {
  echo "$@" | cut -d';' -f2
}