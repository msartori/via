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

// Fields of the Operator.
func (Operator) Fields() []ent.Field {
	return []ent.Field{
		field.String("account").
			Immutable().
			NotEmpty().
			MaxLen(200).
			Unique(),
		field.Bool("enabled").
			Default(false),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// In operator.go > Edges()
func (Operator) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("guides", Guide.Type),
		edge.To("guide_history", GuideHistory.Type),
	}
}
