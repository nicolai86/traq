#!/usr/bin/env bats

. $TRAQ_PATH/traq_helper.sh

@test "traq_tag stop - stop is a special tag which does not get a prefix" {
  result="$(traq_tag stop)"
  [ "$result" = "stop" ]
}

@test "traq_tag tag - regular tags are prefixed with a #" {
  result="$(traq_tag test)"
  [ "$result" = "#test" ]
}

@test "week_number - return the week number of a given date" {
  result="$(week_number '2012-12-13')"
  [ "$result" = "50" ]
}

@test "today_file - generate kw-<week number>/timestamp-<date>" {
  result="$(today_file "2012-12-13" "50")"
  [ "$result" = "kw-50/timestamps-2012-12-13" ]
}

@test "today_file - generate glob if date is *" {
  result="$(today_file '*' "50")"
  [ "$result" = "kw-50/timestamps-*" ]
}
