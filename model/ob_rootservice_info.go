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

type ObRootServiceInfo struct {
	ObClusterId    int64           `json:"ObClusterId"`
	ObRegionId     int64           `json:"ObRegionId"`
	ObCluster      string          `json:"ObCluster"`
	ObRegion       string          `json:"ObRegion"`
	ReadonlyRsList []*ObServerInfo `json:"ReadonlyRsList"`
	RsList         []*ObServerInfo `json:"RsList"`
	Type           string          `json:"Type"`
	TimeStamp      int64           `json:"timestamp"`
}

type ObServerInfo struct {
	Address string `json:"address"`
	Role    string `json:"role"`
	SqlPort int    `json:"sql_port"`
}

func (r *ObRootServiceInfo) Fill() {
	// fill ob cluster and ob region with real
	if len(r.ObCluster) > 0 {
		r.ObRegion = r.ObCluster
	} else if len(r.ObRegion) > 0 {
		r.ObCluster = r.ObRegion
	}

	// fill ob cluster id and ob region id with real
	if r.ObClusterId > 0 {
		r.ObRegionId = r.ObClusterId
	} else if r.ObRegionId > 0 {
		r.ObClusterId = r.ObRegionId
	}
}
