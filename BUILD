load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_library", "go_binary")

go_prefix("github.com/nicolai86")

go_library(
    name = "traq",
    srcs = ["traq.go"],
)

go_binary(
    name = "cli",
    srcs = ["cmd/traq/main.go"],
    deps = [":traq"],
)
