package go_api_test

import (
	"github.com/cs-shuai/go-api-test/common"
	"github.com/cs-shuai/go-api-test/model"
	"testing"
)

func Test(t *testing.T) {
	common.AutoTestRun(t, new(model.JsonTesting), new(model.JsonGroupTesting))
}
