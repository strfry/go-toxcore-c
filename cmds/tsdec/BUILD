load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["tsdec.go"],
    importpath = "github.com/TokTok/go-toxcore-c/cmds/tsdec",
    visibility = ["//visibility:private"],
    deps = ["//go-toxcore-c:go_default_library"],
)

go_binary(
    name = "tsdec",
    embed = [":go_default_library"],
    importpath = "github.com/TokTok/go-toxcore-c/cmds/tsdec",
    visibility = ["//visibility:public"],
)
