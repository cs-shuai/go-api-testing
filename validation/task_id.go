package validation

import (
	"fmt"
	"github.com/gavv/httpexpect"
)

func init() {
	Register(new(TaskId))
}

/**
 * 禅道号
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type TaskId struct {
	Key   string
	Value interface{}
}

func (t *TaskId) GetJsonKey() string {
	t.Key = "task_id"
	return t.Key
}

func (t *TaskId) GetRunFunc() string {
	return TearDownTest
}

func (t *TaskId) SetJsonValue(value interface{}) {
	t.Value = value
}

func (t *TaskId) GetJsonValue() interface{} {
	return t.Value
}

func (t *TaskId) Run(res *httpexpect.Response, params *map[string]interface{}) {
	fmt.Println("执行禅道号", t.Value)
}
