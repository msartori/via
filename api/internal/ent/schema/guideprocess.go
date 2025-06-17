package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// GuideProcess holds the schema definition for the GuideProcess entity.
type GuideProcess struct {
	ent.Schema
}

func (GuideProcess) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").MaxLen(12).Unique().NotEmpty(),
		field.String("recipient").MaxLen(100).Optional(),
		field.String("status").MaxLen(30).Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (GuideProcess) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("operator", Operator.Type).
			Ref("guide_processes").
			Unique(),
		edge.To("history", GuideProcessHistory.Type),
	}
}
