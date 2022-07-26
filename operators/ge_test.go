// Copyright 2022 Juan Pablo Tosso
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package operators

import (
	_ "fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGe(t *testing.T) {
	geo := &ge{}

	err := geo.Init("2500")
	require.NoError(t, err, "Cannot init geo")
	require.True(t, geo.Evaluate(nil, "2800"), "Invalid result for @ge operator")
	require.True(t, geo.Evaluate(nil, "2500"), "Invalid result for @ge operator")
	require.False(t, geo.Evaluate(nil, "2400"), "Invalid result for @ge operator")
}
