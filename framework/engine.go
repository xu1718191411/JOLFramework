package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Engine struct {
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	panicCh := make(chan any)
	successCh := make(chan any)

	jolContext := NewContext(w, r)

	ctx, cancel := context.WithTimeout(jolContext, time.Second*5)

	defer cancel()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicCh <- p
			}
		}()
		time.Sleep(time.Second)
		ControllerWithContext(jolContext)
		//ControllerWithOutContext(w, r)
		successCh <- 1
	}()

	select {
	case <-ctx.Done():
		jolContext.Lock()
		defer jolContext.UnLock()
		jolContext.Send("timeout")
		jolContext.setIsTimeout(true)
	case <-panicCh:
		jolContext.Lock()
		defer jolContext.UnLock()
		jolContext.Send("internal error")
	case <-successCh:
		fmt.Println("success")
	}

}

func ControllerWithContext(ctx *JolContext) {
	output := struct {
		Msg string `json:"msg"`
	}{
		Msg: "ok",
	}

	ctx.Json(output)
}

func ControllerWithOutContext(w http.ResponseWriter, r *http.Request) {
	output := struct {
		Msg string `json:"msg"`
	}{
		Msg: "ok",
	}

	data, err := json.Marshal(output)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}
