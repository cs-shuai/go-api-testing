package go_api_test

import (
	"jccAPITest/common"
	"jccAPITest/model"
	"testing"
)

func Test(t *testing.T) {
	common.AutoTestRun(t, new(model.JsonTesting))
	common.AutoTestRun(t, new(model.JsonGroupTesting))
}
