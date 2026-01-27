package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	//"github.com/Gleipnir-Technology/bob"
	//"github.com/Gleipnir-Technology/bob/dialect/psql"
)

type QueryWriter interface {
	WriteQuery(ctx context.Context, w io.Writer, start int) ([]any, error)
}

func queryToString(query QueryWriter) string {
	buf := new(bytes.Buffer)
	_, err := query.WriteQuery(context.TODO(), buf, 0)
	if err != nil {
		return fmt.Sprintf("Failed to write query to buffer: %v", err)
	}
	return buf.String()
}

/*
func insertQueryToString(query bob.BaseQuery[*dialect.InsertQuery]) string {
	buf := new(bytes.Buffer)
	_, err := query.WriteQuery(context.TODO(), buf, 0)
	if err != nil {
		return fmt.Sprintf("Failed to write query: %v", err)
	}
	return buf.String()
}
*/
