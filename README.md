# Overview

`Narwhal` is an open-source Go module designed to help developers seamlessly integrate container technologies into their applications.

## Project Scope

The library provides high-level abstractions for container operations to perform various tasks. For instance, using container-based artifacts to compile C source code and many other potential use cases.

## Audience

`Narwhal` is intended for developers with a working knowledge of Go programming and Docker operations including SDK. Developers wanting further knowledge of developing in Go using SDK, please refer to Docker's [official SDK documentation](https://docs.docker.com/reference/api/engine/sdk/) or [Learn Go Docker programming](https://github.com/paulwizviz/learn-go-docker) 

## Exported Packages

* [eth](./eth/doc.go)
    * [Example 1 - Compile solidity](./internal/examples/eth/ex1/main.go)
    * [Example 2 - Generate Go Binding](./internal/examples/eth/ex2/main.go)

## Disclaimer

`Narwhal` is provided as-is, without any guarantees or support. While every effort has been made to ensure the library is reliable and functional, users shall assume all risks associated with its use.  

The project is subject to change without notice, with modifications initiated and implemented at the sole discretion of its contributors.  

The contributors to this project shall not be held liable for any damages, losses, or issues arising from its use. Developers are encouraged to thoroughly test and evaluate the library in their own environments before incorporating it into their applications. Where necessary, developers are advised to fork or clone the project and modify it to meet their specific needs.

## Copyright

Unless otherwise specified, this project is copyrighted as follows:

Copyright 2024 The Contributors to `Narwhal`

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

Please refer to the CONTRIBUTORS file for a list of contributors or the repository's commit history.