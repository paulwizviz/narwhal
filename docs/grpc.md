# `grpc` package

This package contains data types and functions to trigger Docker containers embedded with tools to support the development of GRPC-based applications

The artefacts in this package include [protoc compiler](../grpc/grpc.go).

## `protoc` Compiler

The following code snippet demonstrates the core functions to use `protoc` to generate protobuf Go binding:

```go
protoc, err := grpc.NewProtocWithLocalImageLinuxAMD64("narwhal/protoc:current")
if err != nil {
    log.Fatal(err)
}

containerID, err := protoc.CompileProtosGo(context.Background(), "grpc_container", []string{protoPath}, outPath, protoFile)
if err != nil {
    log.Fatal(err)
}
```

Refer to [Example 1](../internal/examples/grpc/ex1/main.go) for a working version that incorporate these functions in an application.

The following code snippet demonstrates the core functions using `protoc` to generate GRPC services in Go binding.

```go
protoc, err := grpc.NewProtocWithLocalImageLinuxAMD64("narwhal/protoc:current")
if err != nil {
    log.Fatal(err)
}

containerID, err := protoc.CompileProtosGRPC(context.Background(), "grpc_container", []string{protoPath}, outPath, protoFile)
if err != nil {
    log.Fatal(err)
}
```

Refer to [Example 2](../internal/examples/grpc/ex2/main.go) for a working version that incorporate these functions in an application.