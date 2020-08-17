package go_api_testing_test

import (
	"github.com/cs-shuai/go-api-testing/common"
	"github.com/cs-shuai/go-api-testing/model"
	"testing"
)

func Test(t *testing.T) {
	common.AutoTestRun(t, new(model.JsonTesting), new(model.JsonGroupTesting))
}
