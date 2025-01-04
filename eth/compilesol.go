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

package eth

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dockersdk "github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func compileSol(ctx context.Context, client *dockersdk.Client, image string, name string, solPath string, solFile string, outPath string, evmVer string, override bool) (string, error) {

	if !isEVMVerCorrect(evmVer) {
		return "", ErrCompileSolEVMVersion
	}

	platform := &v1.Platform{
		OS:           "linux",
		Architecture: "amd64",
	}

	var cmd []string
	if override {
		cmd = []string{"--abi", "--bin", fmt.Sprintf("/opt/solidity/%s", solFile), "-o", "/opt/abi", "--evm-version", evmVer, "--overwrite"}
	} else {
		cmd = []string{"--abi", "--bin", fmt.Sprintf("/opt/solidity/%s", solFile), "-o", "/opt/abi", "--evm-version", evmVer}
	}

	containConfig := &container.Config{
		Image: image,
		Cmd:   cmd,
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: filepath.Join(solPath, solFile),
				Target: fmt.Sprintf("/opt/solidity/%s", solFile),
			},
			{
				Type:   mount.TypeBind,
				Source: outPath,
				Target: "/opt/abi",
			},
		},
	}

	resp, err := client.ContainerCreate(ctx, containConfig, hostConfig, nil, platform, name)
	if err != nil {
		return "", fmt.Errorf("%w-%v", ErrCompileSolEVMCreatContainer, err)
	}

	if err := client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("%w-%v", ErrCompileSolEVMStartContainer, err)
	}

	out, err := client.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: true})
	if err != nil {
		return "", err
	}
	defer out.Close()

	io.Copy(os.Stdout, out)

	return resp.ID, nil
}

func isEVMVerCorrect(version string) bool {
	switch version {
	case EVMVerFrontier,
		EVMVerHomstead,
		EVMVerByzantium,
		EVMVerConstantinople,
		EVMVerIstanbul,
		EVMVerBerlin,
		EVMVerLondon,
		EVMVerShanghai,
		EVMVerCancun,
		EVMVerParis:
		return true
	default:
		return false
	}
}
