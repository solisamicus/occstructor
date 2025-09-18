package model_test

import (
	"fmt"
	"github.com/solisamicus/occstructor/internal/model"
	"testing"
)

func TestGetParentSeq(t *testing.T) {
	fmt.Println(model.GetParentSeq("1"))          // ""
	fmt.Println(model.GetParentSeq("1-01"))       // "1"
	fmt.Println(model.GetParentSeq("1-01-00"))    // "1-01"
	fmt.Println(model.GetParentSeq("1-01-00-01")) // "1-01-00"
}
