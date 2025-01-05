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

	"github.com/paulwizviz/narwhal/pkg/eth"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	packageName := "hello"
	localType := "HelloWorld"

	abiPath := filepath.Join(pwd, "tmp", packageName)
	outPath := filepath.Join(pwd, "tmp", "internal", packageName)
	if _, err := os.Stat(outPath); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(outPath, 0755); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(abiPath)
	fmt.Println(outPath)

	tool, err := eth.NewDefaultTool()
	if err != nil {
		log.Fatal(err)
	}

	containerID, err := tool.GenGoBinding(context.Background(), "alltools-stable", "go-gen", abiPath, outPath, packageName, localType)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(containerID)

	if err := tool.RemoveContainerForce(context.TODO(), containerID); err != nil {
		log.Fatal(err)
	}
}