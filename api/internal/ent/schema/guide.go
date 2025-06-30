package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Guide holds the schema definition for the Guide entity.
type Guide struct {
	ent.Schema
}

// Fields of the Guide.
func (Guide) Fields() []ent.Field {
	return []ent.Field{
		field.String("via_guide_id").
			NotEmpty().
			Immutable().
			MaxLen(12).
			MinLen(12),
		field.String("recipient").
			NotEmpty().
			Immutable().
			MaxLen(100),
		field.String("status").
			NotEmpty().
			MaxLen(30),
		field.String("payment").
			NotEmpty().
			Immutable().
			MaxLen(1).
			MinLen(1),
		field.Int("operator_id"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Guide.
func (Guide) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("operator", Operator.Type).
			Ref("guides").
			Unique().
			Required().
			Field("operator_id"),

		edge.To("history", GuideHistory.Type),
	}
}
