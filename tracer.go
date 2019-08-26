package dbUtil

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/opentracing/opentracing-go"
	ottag "github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/volatiletech/sqlboiler/boil"
)

func WithTracer(executor boil.ContextExecutor) boil.ContextExecutor {
	return &TracerExecutor{
		ContextExecutor: executor,
	}
}

type TracerExecutor struct {
	boil.ContextExecutor
}

func (tr *TracerExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "DB Exec")
	span.LogFields(otlog.String("sql.query", fmt.Sprint(query, ",", args)))
	defer span.Finish()
	result, err := tr.ContextExecutor.ExecContext(ctx, query, args...)
	if err != nil {
		logErrorToSpan(span, err)
	}
	return result, err
}

func (tr *TracerExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "DB Query")
	span.LogFields(otlog.String("sql.query", fmt.Sprint(query, ",", args)))
	defer span.Finish()

	result, err := tr.ContextExecutor.QueryContext(ctx, query, args...)
	if err != nil {
		logErrorToSpan(span, err)
	}
	return result, err
}

func logErrorToSpan(span opentracing.Span, err error) {
	ottag.Error.Set(span, true)
	span.LogFields(otlog.Error(err))
}
