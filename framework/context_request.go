package framework

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
)

func (ctx *JolContext) QueryAll() map[string][]string {
	return ctx.request.URL.Query()
}

func (ctx *JolContext) QueryStringWithDefault(key string, defaultValue string) string {
	v, ok := ctx.QueryString(key)

	if !ok {
		return defaultValue
	}
	return v
}

func (ctx *JolContext) QueryString(key string) (string, bool) {
	queryAll := ctx.QueryAll()
	result := queryAll[key]

	if result == nil {
		return "", false
	}

	if len(result) == 0 {
		return "", false
	}

	return result[len(result)-1], true
}

func (ctx *JolContext) QueryIntWithDefault(key string, defaultValue int) int {
	v, ok := ctx.QueryInt("user_id")
	if !ok {
		return defaultValue
	}
	return v
}

func (ctx *JolContext) QueryInt(key string) (int, bool) {
	queryAll := ctx.QueryAll()
	result := queryAll[key]

	if result == nil {
		return 0, false
	}

	if len(result) == 0 {
		return 0, false
	}

	res := result[len(result)-1]

	v, e := strconv.Atoi(res)
	if e != nil {
		return 0, false
	}
	return v, true
}

func (ctx *JolContext) ParamString(param string) (string, bool) {
	result := ctx.paramsDicts[param]
	if result == "" {
		return "", false
	}

	return result, true
}

func (ctx *JolContext) ParamStringWithDefaultValue(param string, defaultValue string) string {
	v, ok := ctx.ParamString(param)
	if !ok {
		return defaultValue
	}
	return v
}

func (ctx *JolContext) BindJson(target interface{}) error {
	data, err := io.ReadAll(ctx.request.Body)

	if err != nil {
		return err
	}

	// recovery body after read it
	ctx.request.Body = io.NopCloser(bytes.NewBuffer(data))

	return json.Unmarshal(data, target)
}
