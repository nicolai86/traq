#!/usr/bin/env bats

@test "traqtrans converts dates to timestamps" {
  run traqtrans "$BATS_TEST_DIRNAME/fixtures/timestamps-2013-03-23"
  [ ${lines[0]} = "1364047702;#development" ]
  [ ${lines[1]} = "1364052885;stop" ]
  [ ${lines[2]} = "%%" ]
}
