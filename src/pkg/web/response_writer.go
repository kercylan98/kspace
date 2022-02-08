package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ResponseWriter struct {
	gin.ResponseWriter
	*Response
}

func (slf *ResponseWriter) Write(data []byte) (int, error) {
	if slf.Status() == 200 {
		var dataset map[string]any
		if err := json.Unmarshal(data, &dataset); err != nil {
			return 0, errors.New(fmt.Sprintf("%s: %s", err, string(data)))
		}
		slf.Response.Pass(dataset)
		return 0, nil
	} else {
		return 0, errors.New(string(data))
	}
}
