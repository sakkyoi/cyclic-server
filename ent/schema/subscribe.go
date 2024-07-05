package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Subscribe holds the schema definition for the Subscribe entity.
type Subscribe struct {
	ent.Schema
}

// Fields of the Subscribe.
func (Subscribe) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Time("subscribed_at").
			Optional(),
		field.Time("left_at").
			Optional(),
	}
}

// Edges of the Subscribe.
func (Subscribe) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("subscriptions").
			Unique().
			Required(),
		edge.From("plan", Plan.Type).
			Ref("subscriptions").
			Unique().
			Required(),
	}
}

// Indexes of the Subscribe.
func (Subscribe) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("user", "plan").
			Unique(),
	}
}
