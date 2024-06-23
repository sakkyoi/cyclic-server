// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"cyclic/ent/link"
	"cyclic/ent/user"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// LinkCreate is the builder for creating a Link entity.
type LinkCreate struct {
	config
	mutation *LinkMutation
	hooks    []Hook
}

// SetType sets the "type" field.
func (lc *LinkCreate) SetType(l link.Type) *LinkCreate {
	lc.mutation.SetType(l)
	return lc
}

// SetAddress sets the "address" field.
func (lc *LinkCreate) SetAddress(s string) *LinkCreate {
	lc.mutation.SetAddress(s)
	return lc
}

// SetID sets the "id" field.
func (lc *LinkCreate) SetID(u uuid.UUID) *LinkCreate {
	lc.mutation.SetID(u)
	return lc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (lc *LinkCreate) SetNillableID(u *uuid.UUID) *LinkCreate {
	if u != nil {
		lc.SetID(*u)
	}
	return lc
}

// SetOwnerID sets the "owner" edge to the User entity by ID.
func (lc *LinkCreate) SetOwnerID(id uuid.UUID) *LinkCreate {
	lc.mutation.SetOwnerID(id)
	return lc
}

// SetOwner sets the "owner" edge to the User entity.
func (lc *LinkCreate) SetOwner(u *User) *LinkCreate {
	return lc.SetOwnerID(u.ID)
}

// Mutation returns the LinkMutation object of the builder.
func (lc *LinkCreate) Mutation() *LinkMutation {
	return lc.mutation
}

// Save creates the Link in the database.
func (lc *LinkCreate) Save(ctx context.Context) (*Link, error) {
	lc.defaults()
	return withHooks(ctx, lc.sqlSave, lc.mutation, lc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (lc *LinkCreate) SaveX(ctx context.Context) *Link {
	v, err := lc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (lc *LinkCreate) Exec(ctx context.Context) error {
	_, err := lc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lc *LinkCreate) ExecX(ctx context.Context) {
	if err := lc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (lc *LinkCreate) defaults() {
	if _, ok := lc.mutation.ID(); !ok {
		v := link.DefaultID()
		lc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (lc *LinkCreate) check() error {
	if _, ok := lc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "Link.type"`)}
	}
	if v, ok := lc.mutation.GetType(); ok {
		if err := link.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Link.type": %w`, err)}
		}
	}
	if _, ok := lc.mutation.Address(); !ok {
		return &ValidationError{Name: "address", err: errors.New(`ent: missing required field "Link.address"`)}
	}
	if _, ok := lc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner", err: errors.New(`ent: missing required edge "Link.owner"`)}
	}
	return nil
}

func (lc *LinkCreate) sqlSave(ctx context.Context) (*Link, error) {
	if err := lc.check(); err != nil {
		return nil, err
	}
	_node, _spec := lc.createSpec()
	if err := sqlgraph.CreateNode(ctx, lc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	lc.mutation.id = &_node.ID
	lc.mutation.done = true
	return _node, nil
}

func (lc *LinkCreate) createSpec() (*Link, *sqlgraph.CreateSpec) {
	var (
		_node = &Link{config: lc.config}
		_spec = sqlgraph.NewCreateSpec(link.Table, sqlgraph.NewFieldSpec(link.FieldID, field.TypeUUID))
	)
	if id, ok := lc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := lc.mutation.GetType(); ok {
		_spec.SetField(link.FieldType, field.TypeEnum, value)
		_node.Type = value
	}
	if value, ok := lc.mutation.Address(); ok {
		_spec.SetField(link.FieldAddress, field.TypeString, value)
		_node.Address = value
	}
	if nodes := lc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   link.OwnerTable,
			Columns: []string{link.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.user_links = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// LinkCreateBulk is the builder for creating many Link entities in bulk.
type LinkCreateBulk struct {
	config
	err      error
	builders []*LinkCreate
}

// Save creates the Link entities in the database.
func (lcb *LinkCreateBulk) Save(ctx context.Context) ([]*Link, error) {
	if lcb.err != nil {
		return nil, lcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(lcb.builders))
	nodes := make([]*Link, len(lcb.builders))
	mutators := make([]Mutator, len(lcb.builders))
	for i := range lcb.builders {
		func(i int, root context.Context) {
			builder := lcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*LinkMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, lcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, lcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, lcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (lcb *LinkCreateBulk) SaveX(ctx context.Context) []*Link {
	v, err := lcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (lcb *LinkCreateBulk) Exec(ctx context.Context) error {
	_, err := lcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lcb *LinkCreateBulk) ExecX(ctx context.Context) {
	if err := lcb.Exec(ctx); err != nil {
		panic(err)
	}
}
