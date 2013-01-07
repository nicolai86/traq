#!/usr/bin/env bash

. $TRAQ_PATH/os_helper.sh

# converts a date into the corresponding week number
function week_number() {
  is_osx
  OSX=$?
  if [ $OSX -eq 0 ]
  then
    date -j -f "%Y-%m-%d" "$1" "+%V"
  else
    date -d "$1" "+%V"
  fi
}

function year_number() {
  is_osx
  OSX=$?
  if [ $OSX -eq 0 ]
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

# traq logic.
function traq() {
  local TAG=$1
  local DATE=$2
  local WEEK=$3
  local PROJECT=$4
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

  local TRAQFILE="$HOME/.traq/$PROJECT/$YEAR/$(today_file "$DATE" "$WEEK")"

  if [ "$TAG" = "" ]; then # no tag was given; output the content
    for FILE in $TRAQFILE
    do
      cat $FILE
      printf "%%%%\n"
    done
  else # tag was given; handle tag
    if [ "$WEEK" = "" -a "$DATE" = "" ]
    then
      mkdir -p $(dirname $TRAQFILE)
      printf "$(traq_timestamp);$(traq_tag $TAG)\n" >> $TRAQFILE
    else
      printf "can not combine -d or -w with a tag\n" 1>&2
    fi
  fi
}