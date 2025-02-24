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

package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ObCluster holds the schema definition for the ObCluster entity.
type ObCluster struct {
	ent.Schema
}

// Fields of the ObCluster.
func (ObCluster) Fields() []ent.Field {
	return []ent.Field{
		field.Time("create_time").Default(time.Now),
		field.Time("update_time").Default(time.Now).UpdateDefault(time.Now),
		field.String("name"),
		field.Int64("ob_cluster_id").Positive(),
		field.String("type"),
		field.String("rootservice_json").
			Annotations(entsql.Annotation{
				Size: 65536,
			}),
	}
}

func (ObCluster) Edges() []ent.Edge {
	return nil
}

func (ObCluster) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("update_time"),
		index.Fields("name", "ob_cluster_id").Unique(),
	}
}
