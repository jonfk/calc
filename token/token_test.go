package token

import (
	"testing"
)

func TestString(T *testing.T) {
	if ILLEGAL.String() != "ILLEGAL" {
		T.Error("String() not working properly on ILLEGAL token")
	}
	if EOF.String() != "EOF" {
		T.Error("String() not working properly on EOF token")
	}
	if COMMENT.String() != "COMMENT" {
		T.Error("String() not working properly on COMMENT token")
	}
	if INT.String() != "INT" {
		T.Error("String() not working properly on INT token")
	}
	if FLOAT.String() != "FLOAT" {
		T.Error("String() not working properly on FLOAT token")
	}
	if ADD.String() != "+" {
		T.Error("String() not working properly on ADD token")
	}
	if RPAREN.String() != ")" {
		T.Error("String() not working properly on RPAREN token")
	}
}

func TestIsLiteral(T *testing.T) {
	if COMMENT.IsLiteral() {
		T.Error("COMMENT.IsLiteral() not working properly")
	}
	if !INT.IsLiteral() {
		T.Error("INT.IsLiteral() not working properly")
	}
}
