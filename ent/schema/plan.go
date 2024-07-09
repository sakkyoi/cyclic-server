package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// Plan holds the schema definition for the Plan entity.
type Plan struct {
	ent.Schema
}

// Fields of the Plan.
func (Plan) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		// name of the plan
		field.String("name"),
		// description of the plan
		field.String("description").
			Optional().
			Annotations(entsql.Annotation{
				Size: 1000,
			}),
		// price of the plan
		field.Float32("price").
			SchemaType(map[string]string{
				// both of the following are equivalent
				dialect.MySQL:    "decimal(9,2)",
				dialect.Postgres: "numeric(11,2)",
			}),
		// currency of the plan
		field.String("currency"),
		// is it can be exchanged to another currency (automatically to the user's currency)
		field.Bool("exchangeable").
			Default(false),
		// the date when the plan starts
		field.Time("start_from"),
		// duration type of the plan
		field.Enum("duration_type").
			Values("days", "months", "years"),
		// duration value of the plan
		field.Int16("duration"),
		// auto notify the user when the plan is about to expire (payment is required)
		field.Bool("auto_notify").
			Default(true),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("deleted_at").
			Optional(),
	}
}

// Edges of the Plan.
func (Plan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("host", User.Type).
			Ref("plans").
			Unique().
			Required(),
		edge.To("subscriptions", Subscription.Type),
	}
}
