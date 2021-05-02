// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schemast

import (
	"go/ast"

	"entgo.io/ent"
)

// Mutator changes a Context.
type Mutator interface {
	Mutate(ctx *Context) error
}

// Mutate applies a sequence of mutations to a Context
func Mutate(ctx *Context, mutations ...Mutator) error {
	for _, mut := range mutations {
		if err := mut.Mutate(ctx); err != nil {
			return err
		}
	}
	return nil
}

// UpsertSchema implements Mutator. UpsertSchema will add to the Context the type named Name if not present and rewrite
// the type's Fields and Edges methods to return the desired fields and edges.
type UpsertSchema struct {
	Name   string
	Fields []ent.Field
	Edges  []ent.Edge
}

// Mutate applies the UpsertSchema mutation to the Context.
func (u *UpsertSchema) Mutate(ctx *Context) error {
	if !ctx.HasType(u.Name) {
		if err := ctx.AddType(u.Name); err != nil {
			return err
		}
	}
	fieldsReturn, err := ctx.fieldsReturnStmt(u.Name)
	if err != nil {
		return err
	}
	fieldsReturn.Results = []ast.Expr{ast.NewIdent("nil")} // Reset fields.
	edgesReturn, err := ctx.edgesReturnStmt(u.Name)
	if err != nil {
		return err
	}
	edgesReturn.Results = []ast.Expr{ast.NewIdent("nil")} // Reset edges.
	for _, fld := range u.Fields {
		if err := ctx.AppendField(u.Name, fld.Descriptor()); err != nil {
			return err
		}
	}
	for _, edg := range u.Edges {
		if err := ctx.AppendEdge(u.Name, edg.Descriptor()); err != nil {
			return err
		}
	}
	return nil
}
