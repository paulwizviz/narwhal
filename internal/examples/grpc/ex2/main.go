// Copyright 2025 The Contributors to narwhal
// This file is part of the narwhal project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific
// language governing permissions and limitations under the License.
//
// For a list of contributors, refer to the CONTRIBUTORS file or the
// repository's commit history.

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/paulwizviz/narwhal/grpc"
)

func main() {

	// STEP 1: Specify location of protos file and compiled artefacts
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	protoPath := filepath.Join(pwd, "testdata", "protos")
	protoFile := filepath.Join(protoPath, "person.proto")
	outPath := filepath.Join(pwd, "tmp", "go", "ex2")
	if _, err := os.Stat(outPath); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(outPath, 0755); err != nil {
			log.Fatal(err)
		}
	}

	// STEP 3: Instantiate Etherreum tool
	// NOTE: In this case we have a custom build image
	tool, err := grpc.NewProtocWithLocalImageLinuxAMD64("narwhal/protoc:current")
	if err != nil {
		log.Fatal(err)
	}

	// STEP 4: Exxecute function to compile solidity
	containerID, err := tool.CompileProtosGRPC(context.Background(), "grpc_container", []string{protoPath}, outPath, protoFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(containerID)

	// STEP 5: Remove solidity compiler container
	if err := tool.RemoveContainerForce(context.TODO(), containerID); err != nil {
		log.Fatal(err)
	}

}
