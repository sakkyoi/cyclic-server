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
		field.String("name"),
		field.String("description").
			Optional().
			Annotations(entsql.Annotation{
				Size: 1000,
			}),
		field.Float("price").
			SchemaType(map[string]string{
				dialect.MySQL:    "decimal(9,2)",
				dialect.Postgres: "numeric",
			}),
		field.Time("start_from"),
		field.Enum("duration_type").
			Values("days", "months", "years"),
		field.Int16("duration"),
		field.String("status"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Enum("auto_notify").
			Values("automatic", "manual").
			Default("automatic"),
	}
}

// Edges of the Plan.
func (Plan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("host", User.Type).
			Ref("plans").
			Unique().
			Required(),
	}
}
