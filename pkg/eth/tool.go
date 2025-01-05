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

	"github.com/docker/docker/api/types/container"
	dockersdk "github.com/docker/docker/client"
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
	EthereumSolcImage = "ethereum/solc"
)

// Shared operations error
var (
	// ErrCreateClient represents error creating a ethtool client
	ErrCreateClient = errors.New("unable to create client")
	// ErrRemovingContainer represents error removing container
	ErrRemovingContainer = errors.New("unable to remove container")
)

// Compile Solidity function errors
var (
	// ErrCompileSolEVMVersion represent an invalid EVM version declared
	ErrCompileSolEVMVersion = errors.New("invalid evm version")
	// ErrCompileSolEVMCreatContainer represent container creation error
	ErrCompileSolEVMCreatContainer = errors.New("unable to create container")
	// ErrCompileSolEVMStartContainer represent container starting error
	ErrCompileSolEVMStartContainer = errors.New("unable to start container")
)

const (
	OSLinux = "linux"
)

const (
	ArchAMD64 = "amd64"
)

// Tool represents a client of Ethereum tool
type Tool interface {
	// CompileSol is a function trigger a container to compile solidity
	//
	// Arguments:
	//   imageTag        corresponsing to a solidity compiler version
	//   containerName   a unique name of a container
	//   solPath         path to location of solidity contracts
	//   solFile         solidity file name
	//   outPath         path to where the compiled artefact should be
	//   evmVer          version of EVM as per constant value
	CompileSol(ctx context.Context, imageTag string, containerName string, solPath string, solFile string, outPath string, evmVer string) (string, error)

	// CompileSolWithOverride is compile solidity and override existing compile versions
	CompileSolWithOverride(ctx context.Context, imageTag string, containerName string, solPath string, solFile string, outPath string, evmVer string) (string, error)

	// RemoveContainer remove container for a given ID
	RemoveContainer(ctx context.Context, containerID string) error

	// RemoveContainerForce remove container for ID with no exception
	RemoveContainerForce(ctx context.Context, containerID string) error
}

type tool struct {
	cli          *dockersdk.Client
	osPlatform   string
	archPlatform string
}

func (t tool) CompileSol(ctx context.Context, imageTag string, name string, solPath string, solFile string, outPath string, evmVer string) (string, error) {
	image := fmt.Sprintf("%s:%s", EthereumSolcImage, imageTag)
	return compileSol(ctx, t.cli, image, name, t.osPlatform, t.archPlatform, solPath, solFile, outPath, evmVer, false)
}

func (t tool) CompileSolWithOverride(ctx context.Context, imageTag string, name string, solPath string, solFile string, outPath string, evmVer string) (string, error) {
	image := fmt.Sprintf("%s:%s", EthereumSolcImage, imageTag)
	return compileSol(ctx, t.cli, image, name, t.osPlatform, t.archPlatform, solPath, solFile, outPath, evmVer, true)
}

func (t tool) RemoveContainer(ctx context.Context, containerID string) error {
	if err := t.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: false}); err != nil {
		return fmt.Errorf("%w-%v", ErrRemovingContainer, err)
	}
	return nil
}

func (t tool) RemoveContainerForce(ctx context.Context, containerID string) error {
	if err := t.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("%w-%v", ErrRemovingContainer, err)
	}
	return nil
}

// NewDefaultTool is an operation to instantiate an ethtool client with default setting
func NewDefaultTool() (Tool, error) {
	cli, err := dockersdk.NewClientWithOpts(dockersdk.FromEnv, dockersdk.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("%w-%v", ErrCreateClient, err)
	}
	return &tool{
		cli:          cli,
		osPlatform:   OSLinux,
		archPlatform: ArchAMD64,
	}, nil
}
