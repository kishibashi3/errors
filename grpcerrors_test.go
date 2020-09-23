package errors

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestGrpcError(t *testing.T){

	t.Run("test create error", func(t *testing.T) {
		e1 := New("test 1")
		e2 := NewErrorFrom(e1, codes.Aborted, "test %d", 2)
		e3 := NewErrorFrom(e2, codes.Aborted, "test %d", 3)
		e4 := NewErrorFrom(e3, codes.Aborted, "test %d", 4)
		fmt.Printf("%+v", e4)
	})
}
