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

package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dockersdk "github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	// OSLinux is the name of Docker's Linux platform
	OSLinux = "linux"
)

const (
	// ArchAMD64 is the name of Docker's architecture platform
	ArchAMD64 = "amd64"
)

var (
	// ErrCreateTool represents error to create tool
	ErrCreateTool = errors.New("unable to instantiate tool")
	// ErrProtocCreateContainer represents error to create protoc container
	ErrProtocCreateContainer = errors.New("unable to create protoc container")
	// ErrProtoStartContainer represents error to start the ontainer
	ErrProtoStartContainer = errors.New("unable to start protoc container")
	// ErrRemovingContainer represents error to remove protoc container
	ErrRemovingContainer = errors.New("unable to remove container")
)

type Tool interface {
	// CompileProtosGo trigger protoc container to compile protofile
	CompileProtosGo(ctx context.Context, image string, containerName string, protoPath string, outPath string, proto string) (string, error)
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

func (t tool) CompileProtosGo(ctx context.Context, image string, containerName string, protoPath string, outPath string, proto string) (string, error) {
	return compileProtosGo(ctx, t.cli, image, containerName, t.osPlatform, t.archPlatform, protoPath, outPath, proto)
}

func compileProtosGo(ctx context.Context, client *dockersdk.Client, image string, name string, platformOS string, arch string, protoPath string, outPath string, proto string) (string, error) {

	platform := &v1.Platform{
		OS:           platformOS,
		Architecture: arch,
	}

	protofile := filepath.Join(protoPath, proto)

	localProtosDir := "/opt/protos"
	localProtos := fmt.Sprintf("%s/%s", localProtosDir, proto)
	localOutput := fmt.Sprintf("/opt/out")

	cmd := []string{"--proto_path=/usr/local/include", fmt.Sprintf("--proto_path=%s", localProtosDir), fmt.Sprintf("--go_out=%s", localOutput), localProtos}
	containConfig := &container.Config{
		Image: image,
		Cmd:   cmd,
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: protofile,
				Target: fmt.Sprintf("%s/%s", localProtosDir, proto),
			},
			{
				Type:   mount.TypeBind,
				Source: outPath,
				Target: localOutput,
			},
		},
	}

	resp, err := client.ContainerCreate(ctx, containConfig, hostConfig, nil, platform, name)
	if err != nil {
		return "", fmt.Errorf("%w-%v", ErrProtocCreateContainer, err)
	}

	if err := client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("%w-%v", ErrProtoStartContainer, err)
	}

	out, err := client.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: true})
	if err != nil {
		return "", err
	}
	defer out.Close()

	io.Copy(os.Stdout, out)

	return resp.ID, nil
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

func NewDefaultTool() (Tool, error) {
	cli, err := dockersdk.NewClientWithOpts(dockersdk.FromEnv, dockersdk.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("%w-%v", ErrCreateTool, err)
	}
	return &tool{
		cli:          cli,
		osPlatform:   OSLinux,
		archPlatform: ArchAMD64,
	}, nil
}
