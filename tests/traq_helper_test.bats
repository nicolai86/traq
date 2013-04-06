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

@test "traq_entry - returns timestamp with tag if comment is empty" {
  result="$(traq_entry "example")"
  [ "$result" = "$(traq_timestamp);#example;" ]
}

@test "traq_entry - returns timestamp, tag and comment" {
  result="$(traq_entry "example" "this is a comment")"
  [ "$result" = "$(traq_timestamp);#example;this is a comment" ]
}

@test "traq_entry - works with stop as well" {
  result="$(traq_entry "stop" "this is a comment")"
  [ "$result" = "$(traq_timestamp);stop;this is a comment" ]
}

@test "print_traq_file - without existing file prints nothing" {
  run print_traq_file "./barf"
  [ "$output" = "" ]
}

@test "print_traq_file - prints content and delimiter" {
  run print_traq_file "$BATS_TEST_DIRNAME/fixtures/timestamps-2013-03-23"
  [ ${lines[0]} = "Sat Mar 23 15:08:22 +0100 2013;#development;comment" ]
  [ ${lines[1]} = "Sat Mar 23 16:34:45 +0100 2013;stop;" ]
  [ ${lines[2]} = "%%" ]
}
