# `eth` package

This package contains data types and functions to trigger Docker containers embedded with tools to support the development of Ethereum-based application development.

The artefacts in this package include:

* [solc compiler](../eth/solc.go)
* [ABI generator](../eth/abi.go)

## `solc` compiler

The following code snippets shows the core functions to trigger Solidity compiler

```go
solc, err := eth.NewDefaultSolc("0.8.28")
if err != nil {
    log.Fatal(err)
}

containerID, err := solc.CompileSolWithOverride(context.Background(), "solc_container", solPath, solFile, outPath, eth.EVMVerParis)
if err != nil {
    log.Fatal(err)
}
```
or 

```go
containerID, err := solc.CompileSol(context.Background(), "solc_container", solPath, solFile, outPath, eth.EVMVerParis)
if err != nil {
    log.Fatal(err)
}
```

Refer to [Example 1](../internal/examples/eth/ex1/main.go) for a working version incorporated as part of an application.

## ABI Gen -- Go binding generator

Use this package to build application to generate Go binding.

The following code snippet shows the core functions to generate Go binding

```go

abigen, err := eth.NewDefaultProtoc("alltools-stable")
if err != nil {
    log.Fatal(err)
}

containerID, err := abigen.GenGoBinding(context.Background(), "go-gen", abiPath, outPath, packageName, localType)
if err != nil {
    log.Fatal(err)
}
```

Refer to [Example 2](../internal/examples/eth/ex2/main.go) for a working version incorporated as part of an application.