package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("username").
			MinLen(4).
			MaxLen(15).
			Optional().
			Unique(),
		field.Bytes("password").
			Optional().
			Sensitive(),
		field.String("email").
			Optional().
			Unique(),
		field.String("name").
			Default("unknown"),
		field.String("role").
			Default("user"),
		field.Bool("active").
			Default(false),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("links", Link.Type),
		edge.To("plans", Plan.Type),
		edge.From("subscriptions", Subscribe.Type).
			Ref("users"),
	}
}
