package go_api_test

import (
	"jccAPITest/common"
	"jccAPITest/model"
	"testing"
)

func Test(t *testing.T) {
	common.AutoTestRun(t, new(model.JsonTesting), new(model.JsonGroupTesting))
}
