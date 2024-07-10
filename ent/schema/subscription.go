package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

// Subscription holds the schema definition for the Subscription entity.
type Subscription struct {
	ent.Schema
}

// Fields of the Subscription.
func (Subscription) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("start_from"),
		field.Time("subscribed_at").
			Optional(),
		field.Time("left_at").
			Optional(),
	}
}

// Edges of the Subscription.
func (Subscription) Edges() []ent.Edge {
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

// Indexes of the Subscription.
func (Subscription) Indexes() []ent.Index {
	return []ent.Index{
		// add left_at to the index to let the plan is re-subscribable
		index.Edges("user", "plan").
			Unique(),
	}
}
