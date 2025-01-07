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
	// EthereumGethToolImage is the name of Geth tool Docker image
	EthereumGethToolImage = "ethereum/client-go"
)

var (
	// ErrCreateABIGenClient represents error creating protoc Docker client
	ErrCreateABIGenClient = errors.New("unable to create docker abigen client")
	// ErrCreateABIGenContainer represents error creating ABIGen container
	ErrCreateABIGenContainer = errors.New("unable to create abigen container")
	// ErrRemoveABIGenContainer represents error removing ABIGen container
	ErrRemoveABIGenContainer = errors.New("unable to remove abigen container")
	// ErrCreateABIGenContainer represents error staring ABIGen container
	ErrStartABIGenContainer = errors.New("unable to start abigen container")
)

// ABIGen is an abstraction of Ethereum ABIGen docker client
type ABIGen interface {
	// GenGoBinding generates Go binding
	GenGoBinding(ctx context.Context, name string, abiPath string, outPath string, pkgName string, localType string) (string, error)
	// RemoveContainer remove container for a given ID
	RemoveContainer(ctx context.Context, containerID string) error
	// RemoveContainerForce remove container for ID with no exception
	RemoveContainerForce(ctx context.Context, containerID string) error
}

type abigen struct {
	cli          *dockersdk.Client
	osPlatform   string
	archPlatform string
	image        string
}

func (a abigen) GenGoBinding(ctx context.Context, name string, abiPath string, outPath string, pkgName string, localType string) (string, error) {
	return generateGoBinding(ctx, a.cli, a.image, name, a.osPlatform, a.archPlatform, abiPath, outPath, pkgName, localType)
}

func generateGoBinding(ctx context.Context, client *dockersdk.Client, image string, name string, platformOS string, arch string, abiPath string, outPath string, pkgName string, localType string) (string, error) {

	platform := &v1.Platform{
		OS:           platformOS,
		Architecture: arch,
	}

	localABIFolder := "/opt/abi"
	localBindingFolder := "/opt/binding"

	containConfig := &container.Config{
		Image: image,
		Cmd:   []string{"abigen", "--abi", fmt.Sprintf("%s/%s.abi", localABIFolder, localType), "--bin", fmt.Sprintf("/opt/abi/%s.bin", localType), "--pkg", pkgName, "--type", localType, "--out", fmt.Sprintf("/opt/binding/%s/%s.go", pkgName, localType)},
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: filepath.Join(abiPath, fmt.Sprintf("%s.abi", localType)),
				Target: fmt.Sprintf("%s/%s", localABIFolder, fmt.Sprintf("%s.abi", localType)),
			},
			{
				Type:   mount.TypeBind,
				Source: filepath.Join(abiPath, fmt.Sprintf("%s.bin", localType)),
				Target: fmt.Sprintf("%s/%s", localABIFolder, fmt.Sprintf("%s.bin", localType)),
			},
			{
				Type:   mount.TypeBind,
				Source: outPath,
				Target: fmt.Sprintf("%s/%s", localBindingFolder, pkgName),
			},
		},
	}

	resp, err := client.ContainerCreate(ctx, containConfig, hostConfig, nil, platform, name)
	if err != nil {
		return "", fmt.Errorf("%w-%v", ErrCreateABIGenContainer, err)
	}

	if err := client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("%w-%v", ErrStartABIGenContainer, err)
	}

	out, err := client.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: true})
	if err != nil {
		return "", err
	}
	defer out.Close()

	io.Copy(os.Stdout, out)

	return resp.ID, nil
}

func (a abigen) RemoveContainer(ctx context.Context, containerID string) error {
	if err := a.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: false}); err != nil {
		return fmt.Errorf("%w-%v", ErrRemoveABIGenContainer, err)
	}
	return nil
}

func (a abigen) RemoveContainerForce(ctx context.Context, containerID string) error {
	if err := a.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("%w-%v", ErrRemoveABIGenContainer, err)
	}
	return nil
}

// NewDefaultProtoc instantiate an ethereum/client-go client for Linux/amd64 platform
//
// Arguments:
//
// - imgTag is the tag associated with ethereum/client-go
func NewDefaultProtoc(imgTag string) (ABIGen, error) {

	cli, err := dockersdk.NewClientWithOpts(dockersdk.FromEnv, dockersdk.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("%w-%v", ErrCreateABIGenClient, err)
	}
	gethToolImage := fmt.Sprintf("%s:%s", EthereumGethToolImage, imgTag)

	p := shared.PlatformLinuxAMD64()
	reader, err := cli.ImagePull(context.Background(), gethToolImage, image.PullOptions{
		Platform: fmt.Sprintf("%s/%s", p.OS, p.Arch),
	})
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	return &abigen{
		cli:          cli,
		osPlatform:   p.OS,
		archPlatform: p.Arch,
		image:        gethToolImage,
	}, nil
}
