#!/usr/bin/env bash

# converts a date into the corresponding week number
function week_number() {
  if [ "$(uname)" == "Darwin" ]
  then
    date -j -f "%Y-%m-%d" "$1" "+%V"
  else
    date -d "$1" "+%V"
  fi
}

function year_number() {
  if [ "$(uname)" == "Darwin" ]
  then
    date -j -f "%Y-%m-%d" "$1" "+%Y"
  else
    date -d "$1" "+%Y"
  fi
}

function current_week_number() {
  date "+%V"
}

function current_year() {
  date "+%Y"
}

# output the current date like this: 'Sun Sep 23 11:28:02 +0400 2012'
function traq_timestamp() {
  printf "$(date +"%a %b %d %H:%M:%S %z %Y")"
}

# generate the correct tag for a given tag
function traq_tag() {
  if [ "$@" == "stop" ]
  then
    printf "stop\n"
  else
    printf "#$@\n"
  fi
}

# todays date, formatted for traq
function traq_date() {
  date "+%Y-%m-%d"
}

# return kw-<week>/timestamps-<current date>
function today_file() {
  local DATE=$1
  local WEEK_NUMBER=$2

  # no date given. use today
  if [ "$DATE" = "" ]; then
    DATE=$(traq_date)
    WEEK_NUMBER=$(week_number $DATE)
  fi

  printf "kw-$WEEK_NUMBER/timestamps-$DATE"
}

# creates a traq entry
function traq_entry() {
  local TAG=$1
  local COMMENT=$2

  printf "$(traq_timestamp);$(traq_tag $TAG);$COMMENT"
}

# echos the content of a file with delimiter
function print_traq_file() {
  local FILE=$1

  if [ -f $FILE ]; then
    cat $FILE
    printf "%%%%\n"
  fi
}

# traq logic.
function traq() {
  local TAG=$1
  local DATE=$2
  local WEEK=$3
  local PROJECT=$4
  local COMMENT=$5
  local YEAR=$(current_year)

  # no arguments given. use todays date to output
  if [ "$DATE" = "" -a "$TAG" = "" -a "$WEEK" = "" ]; then
    DATE=$(traq_date)
    WEEK="$(current_week_number)"
  fi
  if [ "$DATE" != "" ]; then
    WEEK="$(week_number $DATE)"
    YEAR="$(year_number $DATE)"
  fi
  # week given. use glob
  if [ "$DATE" = "" -a "$TAG" = "" -a "$WEEK" != "" ]; then
    DATE='*'
  fi

  local TRAQFILE="$TRAQ_DATA_DIR/$PROJECT/$YEAR/$(today_file "$DATE" "$WEEK")"

  if [ "$TAG" = "" ]; then
    # no tag was given; output the content, but only if the file exists
    for FILE in $TRAQFILE
    do
      print_traq_file "$FILE"
    done
  else
    # tag was given; handle tag
    if [ "$WEEK" = "" -a "$DATE" = "" ]
    then
      mkdir -p $(dirname $TRAQFILE)
      printf "$(traq_entry "$TAG" "$COMMENT")\n" >> $TRAQFILE
    else
      printf "can not combine -d or -w with a tag\n" 1>&2
    fi
  fi
}