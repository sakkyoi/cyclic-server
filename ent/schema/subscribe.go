package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
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
			Default(time.Now),
		field.Time("left_at").
			Optional(),
	}
}

// Edges of the Subscribe.
func (Subscribe) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}
