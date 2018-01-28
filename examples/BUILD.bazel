load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["toxecho.go"],
    importpath = "github.com/TokTok/go-toxcore-c/examples",
    visibility = ["//visibility:private"],
    deps = ["//go-toxcore-c:go_default_library"],
)

go_binary(
    name = "examples",
    embed = [":go_default_library"],
    importpath = "github.com/TokTok/go-toxcore-c/examples",
    visibility = ["//visibility:public"],
)
