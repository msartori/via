package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// GuideProcessHistory holds the schema definition for the GuideProcessHistory entity.
type GuideProcessHistory struct {
	ent.Schema
}

func (GuideProcessHistory) Fields() []ent.Field {
	return []ent.Field{
		field.String("status").MaxLen(30).Optional(),
		field.Time("created_at").Default(time.Now),
	}
}

func (GuideProcessHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guide_process", GuideProcess.Type).
			Ref("history").
			Required().
			Unique(),
		edge.From("operator", Operator.Type).
			Ref("guide_process_histories").
			Unique(),
	}
}
