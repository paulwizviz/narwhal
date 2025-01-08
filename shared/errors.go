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

package shared

import (
	"errors"
	"fmt"
)

var (
	// ErrInstantiateClient represents error instantiating a Docker client
	ErrInstantiateClient = errors.New("unable to instantiate a docker client")
	// ErrCreateContainer represents error creating container
	ErrCreateContainer = errors.New("unable to create a container")
	// ErrStartContainer represents error starting container
	ErrStartContainer = errors.New("unable to start a container")
	// ErrRemoveContainer represents error removing container
	ErrRemoveContainer = errors.New("unable to remove a container")
	// ErrContainerLog represents error instantiate a container log
	ErrContainerLog = errors.New("unable to instantiate container log")
	// ErrPullImage represents error pulling an image
	ErrPullImage = errors.New("unable to pull image")
)

// InstantiateClientErr returns an error handler instatiating a client
func InstantiateClientErr(err error, pkg string, fname string) error {
	return fmt.Errorf("%w-%s-%s-%v", ErrInstantiateClient, pkg, fname, err)
}

// CreateContainerErr returns an error handler creating a container
func CreateContainerErr(err error, pkg string, fname string) error {
	return fmt.Errorf("%w-%s-%s-%v", ErrCreateContainer, pkg, fname, err)
}

// StartContainerErr returns an error handler starting a container
func StartContainerErr(err error, pkg string, fname string) error {
	return fmt.Errorf("%w-%s-%s-%v", ErrStartContainer, pkg, fname, err)
}

// RemoveContainerErr returns an error handler removing a container
func RemoveContainerErr(err error, pkg string, fname string) error {
	return fmt.Errorf("%w-%s-%s-%v", ErrRemoveContainer, pkg, fname, err)
}

// ContainerLogErr returns an error handler removing a container
func ContainerLogErr(err error, pkg string, fname string) error {
	return fmt.Errorf("%w-%s-%s-%v", ErrRemoveContainer, pkg, fname, err)
}

// PullImageError returns an error handler pulling image
func PullImageError(err error, pkg string, fname string) error {
	return fmt.Errorf("%w-%s-%s-%v", ErrPullImage, pkg, fname, err)
}
