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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEVMVerCorrect(t *testing.T) {
	testcases := []struct {
		input string
		want  bool
	}{
		{
			input: "frontier",
			want:  true,
		},
		{
			input: "homestead",
			want:  true,
		},
		{
			input: "byzantium",
			want:  true,
		},
		{
			input: "constantinople",
			want:  true,
		},
		{
			input: "istanbul",
			want:  true,
		},
		{
			input: "berlin",
			want:  true,
		},
		{
			input: "london",
			want:  true,
		},
		{
			input: "shanghai",
			want:  true,
		},
		{
			input: "cancun",
			want:  true,
		},
		{
			input: "paris",
			want:  true,
		},
		{
			input: "Cancun",
			want:  false,
		},
		{
			input: "hello",
			want:  false,
		},
	}
	for i, tc := range testcases {
		got := isEVMVerCorrect(tc.input)
		assert.Equal(t, tc.want, got, fmt.Sprintf("Case: %d Want: %v Got: %v", i, tc.want, got))
	}
}
