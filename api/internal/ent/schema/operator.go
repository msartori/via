package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Operator holds the schema definition for the Operator entity.
type Operator struct {
	ent.Schema
}

func (Operator) Fields() []ent.Field {
	return []ent.Field{
		field.String("account").Unique().NotEmpty(),
		field.Bool("enabled"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Operator) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("guide_processes", GuideProcess.Type),
		edge.To("guide_process_histories", GuideProcessHistory.Type),
	}
}
