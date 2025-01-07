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
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	dockersdk "github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/paulwizviz/narwhal/shared"
)

const (
	// EVMVerFrontier version name for frontier EVM
	EVMVerFrontier = "frontier"
	// EVMVerHomstead version name for homestead EVM
	EVMVerHomstead = "homestead"
	// EVMVerByzantium version name for homestead EVM
	EVMVerByzantium = "byzantium"
	// EVMVerConstantinople version name for constantinople EVM
	EVMVerConstantinople = "constantinople"
	// EVMVerIstanbul version name for istanbul EVM
	EVMVerIstanbul = "istanbul"
	// EVMVerBerlin version name for berlin EVM
	EVMVerBerlin = "berlin"
	// EVMVerLondon version name for london EVM
	EVMVerLondon = "london"
	// EVMVerShanghai version name for shanghai EVM
	EVMVerShanghai = "shanghai"
	// EVMVerCancun version name for cancun EVM
	EVMVerCancun = "cancun"
	// EVMVerParis version name for paris EVM
	EVMVerParis = "paris"
)

const (
	// EthereumSolcImage is the name Ethereum Solidity Compiler Docker image
	EthereumSolcImage = "ethereum/solc"
)

var (
	// ErrSolcClient represents error creating a ethtool client
	ErrSolcClient = errors.New("unable to create solc client")
	// ErrSolcCreatContainer represent container creation error
	ErrSolcCreatContainer = errors.New("unable to create solc container")
	// ErrRemovingSolcContainer represents error removing container
	ErrRemovingSolcContainer = errors.New("unable to remove container")
	// ErrSolcStartContainer represent container starting error
	ErrSolcStartContainer = errors.New("unable to start solc container")
)

var (
	// ErrInvalidEVMVersion represent an invalid EVM version declared
	ErrInvalidEVMVersion = errors.New("invalid evm version")
)

// Solc represents docker clients that wrap solidity compiler
type Solc interface {

	// CompileSol is a function trigger a container to compile solidity
	//
	// Arguments:
	//
	//	- containerName   a unique name of a container
	//	- solPath         path to location of solidity contracts
	//	- solFile         solidity file name
	//	- outPath         path to where the compiled artefact should be
	//	- evmVer          version of EVM as per constant value
	CompileSol(ctx context.Context, containerName string, solPath string, solFile string, outPath string, evmVer string) (string, error)
	// CompileSolWithOverride is compile solidity and override existing compile versions
	CompileSolWithOverride(ctx context.Context, containerName string, solPath string, solFile string, outPath string, evmVer string) (string, error)
	// RemoveContainer remove container for a given ID
	RemoveContainer(ctx context.Context, containerID string) error
	// RemoveContainerForce remove container for ID with no exception
	RemoveContainerForce(ctx context.Context, containerID string) error
}

type solc struct {
	cli          *dockersdk.Client
	osPlatform   string
	archPlatform string
	image        string
}

func (s solc) CompileSol(ctx context.Context, name string, solPath string, solFile string, outPath string, evmVer string) (string, error) {
	return compileSol(ctx, s.cli, s.image, name, s.osPlatform, s.archPlatform, solPath, solFile, outPath, evmVer, false)
}

func (s solc) CompileSolWithOverride(ctx context.Context, name string, solPath string, solFile string, outPath string, evmVer string) (string, error) {
	return compileSol(ctx, s.cli, s.image, name, s.osPlatform, s.archPlatform, solPath, solFile, outPath, evmVer, true)
}

func compileSol(ctx context.Context, client *dockersdk.Client, image string, name string, platformOS string, arch string, solPath string, solFile string, outPath string, evmVer string, override bool) (string, error) {

	if !isEVMVerCorrect(evmVer) {
		return "", ErrInvalidEVMVersion
	}

	platform := &v1.Platform{
		OS:           platformOS,
		Architecture: arch,
	}

	localSolFolder := "/opt/solidity"
	localABIFolder := "/opt/abi"

	var cmd []string
	if override {
		cmd = []string{"--abi", "--bin", fmt.Sprintf("%s/%s", localSolFolder, solFile), "-o", localABIFolder, "--evm-version", evmVer, "--overwrite"}
	} else {
		cmd = []string{"--abi", "--bin", fmt.Sprintf("%s/%s", localSolFolder, solFile), "-o", localABIFolder, "--evm-version", evmVer}
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
				Target: fmt.Sprintf("%s/%s", localSolFolder, solFile),
			},
			{
				Type:   mount.TypeBind,
				Source: outPath,
				Target: localABIFolder,
			},
		},
	}

	resp, err := client.ContainerCreate(ctx, containConfig, hostConfig, nil, platform, name)
	if err != nil {
		return "", fmt.Errorf("%w-%v", ErrSolcCreatContainer, err)
	}

	if err := client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("%w-%v", ErrSolcStartContainer, err)
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

func (s solc) RemoveContainer(ctx context.Context, containerID string) error {
	if err := s.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: false}); err != nil {
		return fmt.Errorf("%w-%v", ErrRemovingSolcContainer, err)
	}
	return nil
}

func (t solc) RemoveContainerForce(ctx context.Context, containerID string) error {
	if err := t.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("%w-%v", ErrRemovingSolcContainer, err)
	}
	return nil
}

// NewDefaultSolc instantiate an ethereum/solc client for Linux/amd64 platform
//
// Arguments:
//
// - imgTag is the tag associated with ethereum/solc
func NewDefaultSolc(imageTag string) (Solc, error) {
	cli, err := dockersdk.NewClientWithOpts(dockersdk.FromEnv, dockersdk.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("%w-%v", ErrSolcClient, err)
	}
	solcImage := fmt.Sprintf("%s:%s", EthereumSolcImage, imageTag)

	p := shared.PlatformLinuxAMD64()
	reader, err := cli.ImagePull(context.Background(), solcImage, image.PullOptions{
		Platform: fmt.Sprintf("%s/%s", p.OS, p.Arch),
	})
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)
	return &solc{
		cli:          cli,
		osPlatform:   p.OS,
		archPlatform: p.Arch,
		image:        solcImage,
	}, nil
}
