package errors

import (
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestGrpcError(t *testing.T) {

	t.Run("test create error", func(t *testing.T) {
		e1 := New("test 1")
		e2 := NewErrorFrom(e1, codes.Aborted, "test %d", 2)
		e3 := NewErrorFrom(e2, codes.Aborted, "test %d", 3)
		e4 := NewErrorFrom(e3, codes.Aborted, "test %d", 4)
		fmt.Printf("%+v", e4)

		t.Error("error")
	})

	t.Run("test stack", func(t *testing.T) {

		hoge := func() error {
			e := New("error1")
			return NewErrorFrom(e, codes.Internal, "error2")
		}
		e := hoge()
		e2 := WithMessage(e, "error3")

		fmt.Printf("%+v", e2)

		t.Error("e")
	})

}
