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

type ObClusterIdcRegionInfo struct {
	Cluster        string           `json:"ObRegion"`
	ClusterId      int64            `json:"ObRegionId"`
	IdcList        []*IdcRegionInfo `json:"IDCList"`
	ReadonlyRsList string           `json:"ReadonlyRsList"`
}

type IdcRegionInfo struct {
	Idc    string `json:"idc"`
	Region string `json:"region"`
}
