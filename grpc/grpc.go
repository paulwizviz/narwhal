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
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dockersdk "github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/paulwizviz/narwhal/shared"
)

type Protoc interface {
	// CompileProtosGo trigger protoc container to compile protofile
	CompileProtosGo(ctx context.Context, containerName string, protoPath []string, outPath string, proto string) (string, error)
	// CompileProtosGRPC trigger protoc container to compile protofile for grpc output
	CompileProtosGRPC(ctx context.Context, containerName string, protoPath []string, outPath string, proto string) (string, error)
	// RemoveContainer remove container for a given ID
	RemoveContainer(ctx context.Context, containerID string) error
	// RemoveContainerForce remove container for ID with no exception
	RemoveContainerForce(ctx context.Context, containerID string) error
}

type protoc struct {
	cli          *dockersdk.Client
	osPlatform   string
	archPlatform string
	image        string
}

func (p protoc) CompileProtosGo(ctx context.Context, containerName string, protoPaths []string, outPath string, proto string) (string, error) {
	return compileProtosGo(ctx, p.cli, p.image, containerName, p.osPlatform, p.archPlatform, protoPaths, outPath, proto)
}

func compileProtosGo(ctx context.Context, client *dockersdk.Client, image string, name string, platformOS string, arch string, protoPaths []string, outPath string, protoFile string) (string, error) {

	platform := &v1.Platform{
		OS:           platformOS,
		Architecture: arch,
	}

	var localProtoPaths string
	for _, pp := range protoPaths {
		localProtoPaths = localProtoPaths + fmt.Sprintf("--proto_path=%s,", pp)
	}

	localProtoPaths = localProtoPaths[:len(localProtoPaths)-1]
	localOutput := fmt.Sprintf("/opt/out")

	cmd := []string{"--proto_path=/usr/local/include", localProtoPaths,
		fmt.Sprintf("--go_out=%s", localOutput),
		"--go_opt=paths=source_relative",
		protoFile}

	containConfig := &container.Config{
		Image: image,
		Cmd:   cmd,
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: protoFile,
				Target: protoFile,
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
		return "", shared.CreateContainerErr(err, "grpc", "compileProtosGo")
	}

	if err := client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", shared.StartContainerErr(err, "grpc", "compileProtosGo")
	}

	out, err := client.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: true})
	if err != nil {
		return "", shared.ContainerLogErr(err, "grpc", "compileProtosGo")
	}
	defer out.Close()

	io.Copy(os.Stdout, out)

	return resp.ID, nil
}

func (p protoc) CompileProtosGRPC(ctx context.Context, containerName string, protoPaths []string, outPath string, proto string) (string, error) {
	return compileProtosGRPC(ctx, p.cli, p.image, containerName, p.osPlatform, p.archPlatform, protoPaths, outPath, proto)
}

func compileProtosGRPC(ctx context.Context, client *dockersdk.Client, image string, name string, platformOS string, arch string, protoPaths []string, outPath string, protoFile string) (string, error) {

	platform := &v1.Platform{
		OS:           platformOS,
		Architecture: arch,
	}

	var localProtoPaths string
	for _, pp := range protoPaths {
		localProtoPaths = localProtoPaths + fmt.Sprintf("--proto_path=%s,", pp)
	}

	localProtoPaths = localProtoPaths[:len(localProtoPaths)-1]
	localOutput := fmt.Sprintf("/opt/out")

	cmd := []string{"--proto_path=/usr/local/include", localProtoPaths,
		fmt.Sprintf("--go_out=%s", localOutput),
		"--go_opt=paths=source_relative",
		fmt.Sprintf("--go-grpc_out=%s", localOutput),
		"--go-grpc_opt=paths=source_relative",
		protoFile}
	containConfig := &container.Config{
		Image: image,
		Cmd:   cmd,
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: protoFile,
				Target: protoFile,
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
		return "", shared.CreateContainerErr(err, "grpc", "compileProtosGRPC")
	}

	if err := client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", shared.StartContainerErr(err, "grpc", "compileProtosGRPC")
	}

	out, err := client.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: true})
	if err != nil {
		return "", shared.ContainerLogErr(err, "grpc", "compileProtosGRPC")
	}
	defer out.Close()

	io.Copy(os.Stdout, out)

	return resp.ID, nil
}

func (p protoc) RemoveContainer(ctx context.Context, containerID string) error {
	return shared.RemoveContainer(ctx, p.cli, containerID)
}

func (p protoc) RemoveContainerForce(ctx context.Context, containerID string) error {
	return shared.RemoveContainerForce(ctx, p.cli, containerID)
}

// NewProtocWithLocalImageLinuxAMD64 instantiate a user specified image base on Linux and AMD64 platform
func NewProtocWithLocalImageLinuxAMD64(img string) (Protoc, error) {
	cli, err := dockersdk.NewClientWithOpts(dockersdk.FromEnv, dockersdk.WithAPIVersionNegotiation())
	if err != nil {
		return nil, shared.InstantiateClientErr(err, "grpc", "NewProtocWithLocalImageLinuxAMD64")
	}
	p := shared.PlatformLinuxAMD64()
	return &protoc{
		cli:          cli,
		osPlatform:   p.OS,
		archPlatform: p.Arch,
		image:        img,
	}, nil
}
