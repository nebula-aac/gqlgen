package graphql

import (
	"context"
	"io"
)

var (
	nullLit      = []byte(`null`)
	trueLit      = []byte(`true`)
	falseLit     = []byte(`false`)
	openBrace    = []byte(`{`)
	closeBrace   = []byte(`}`)
	openBracket  = []byte(`[`)
	closeBracket = []byte(`]`)
	colon        = []byte(`:`)
	comma        = []byte(`,`)
)

var (
	Null  = &lit{nullLit}
	True  = &lit{trueLit}
	False = &lit{falseLit}
)

type Marshaler interface {
	MarshalGQL(w io.Writer)
}

type Unmarshaler interface {
	UnmarshalGQL(v any) error
}

type ContextMarshaler interface {
	MarshalGQLContext(ctx context.Context, w io.Writer) error
}

type ContextUnmarshaler interface {
	UnmarshalGQLContext(ctx context.Context, v any) error
}

type contextMarshalerAdapter struct {
	Context context.Context
	ContextMarshaler
}

func WrapContextMarshaler(ctx context.Context, m ContextMarshaler) Marshaler {
	return contextMarshalerAdapter{Context: ctx, ContextMarshaler: m}
}

func (a contextMarshalerAdapter) MarshalGQL(w io.Writer) {
	err := a.MarshalGQLContext(a.Context, w)
	if err != nil {
		AddError(a.Context, err)
		Null.MarshalGQL(w)
	}
}

type WriterFunc func(writer io.Writer)

func (f WriterFunc) MarshalGQL(w io.Writer) {
	f(w)
}

type ContextWriterFunc func(ctx context.Context, writer io.Writer) error

func (f ContextWriterFunc) MarshalGQLContext(ctx context.Context, w io.Writer) error {
	return f(ctx, w)
}

type Array []Marshaler

func (a Array) MarshalGQL(writer io.Writer) {
	writer.Write(openBracket)
	var notEmpty bool
	for _, val := range a {
		if _, ok := val.(Empty); ok {
			continue
		}
		if notEmpty {
			writer.Write(comma)
		}
		val.MarshalGQL(writer)
		notEmpty = true
	}
	writer.Write(closeBracket)
}

type lit struct{ b []byte }

func (l lit) MarshalGQL(w io.Writer) {
	w.Write(l.b)
}

func (l lit) MarshalGQLContext(_ context.Context, w io.Writer) error {
	w.Write(l.b)
	return nil
}

type Empty struct{}

func (e Empty) MarshalGQL(_ io.Writer) {}

func (e Empty) MarshalGQLContext(_ context.Context, _ io.Writer) error {
	return nil
}
