#!/usr/bin/env bats

. $TRAQ_PATH/traqtrans_helper.sh

@test "extract_date returns correct timestamp" {
  result="$(extract_date "Thu Sep 27 07:05:05 +0400 2012;#foo")"
  [ "$result" = "Thu Sep 27 07:05:05 +0400 2012" ]
}

@test "extract_tag returns correct tag" {
  result="$(extract_tag "Thu Sep 27 07:05:05 +0400 2012;#foo")"
  [ "$result" = "#foo" ]
}

@test "performs correct date->timestamp conversion" {
  result="$(date2timestamp "Thu Sep 27 07:05:05 +0400 2012")"
  [ "$result" = "1348715105" ]
}
