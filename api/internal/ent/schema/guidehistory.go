package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// GuideHistory holds the schema definition for the GuideHistory entity.
type GuideHistory struct {
	ent.Schema
}

// Fields of the GuideHistory.
func (GuideHistory) Fields() []ent.Field {
	return []ent.Field{
		field.Int("guide_id").
			Immutable(),
		field.String("status").
			Immutable().
			NotEmpty().
			MaxLen(30),
		field.Int("operator_id").
			Immutable().
			Positive(),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the GuideHistory.
func (GuideHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guide", Guide.Type).
			Ref("history").
			Required().
			Unique().
			Field("guide_id").
			Immutable(),

		edge.From("operator", Operator.Type).
			Ref("guide_history").
			Required().
			Unique().
			Field("operator_id").
			Immutable(),
	}
}
