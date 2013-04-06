#!/usr/bin/env bats

@test "traqeval calculates time between timestamps" {
  run traqeval < <(traqtrans "$BATS_TEST_DIRNAME/fixtures/timestamps-2013-03-23")
  [ ${lines[0]} = "2013-03-23" ]
  [ ${lines[1]} = "#development:1.4397" ]
  [ ${lines[2]} = "%%" ]
}
