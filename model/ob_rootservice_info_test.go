/**
 * Copyright 2025 OceanBase
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Service struct {
	Address string `json:"address"`
}

func TestFillObRegion(t *testing.T) {
	info := &ObRootServiceInfo{
		ObCluster:   "helloworld",
		ObClusterId: 1,
	}
	info.Fill()
	require.Equal(t, int64(1), info.ObRegionId)
	require.Equal(t, "helloworld", info.ObRegion)
}

func TestFillObCluster(t *testing.T) {
	info := &ObRootServiceInfo{
		ObRegion:   "helloworld",
		ObRegionId: 1,
	}
	info.Fill()
	require.Equal(t, int64(1), info.ObClusterId)
	require.Equal(t, "helloworld", info.ObCluster)
}
