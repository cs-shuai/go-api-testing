package tests

import (
	"fmt"
	"gopkg.in/check.v1"
	"jccAPITest/common"
	_ "jccAPITest/model"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
	fmt.Println(common.TestList)
	for _, test := range common.TestList {
		var _ = check.Suite(test)
	}
}
