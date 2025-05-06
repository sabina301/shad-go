//go:build !solution

package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func truePlaceForCtx(method reflect.Method, i int) bool {
	return method.Type.In(i) == reflect.TypeOf((*context.Context)(nil)).Elem()
}

func trueSecondPlace(method reflect.Method) bool {
	return method.Type.In(2).Kind() == reflect.Ptr && method.Type.In(2).Elem().Kind() == reflect.Struct
}

func trueNumOfInAndOut(method reflect.Method, i, j int) bool {
	return method.Type.NumIn() == i && method.Type.NumOut() == j
}

func trueZeroOut(method reflect.Method) bool {
	return method.Type.Out(0).Kind() == reflect.Ptr && method.Type.Out(0).Elem().Kind() == reflect.Struct
}

func trueFirstOut(method reflect.Method) bool {
	return method.Type.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem())
}

func successValidateMethod(method reflect.Method) bool {
	if !(truePlaceForCtx(method, 1) && trueNumOfInAndOut(method, 3, 2) &&
		trueSecondPlace(method) && trueZeroOut(method) && trueFirstOut(method)) {
		return false
	}
	return true
}

func MakeHandler(service interface{}) http.Handler {
	mux := http.NewServeMux()
	tService := reflect.TypeOf(service)

	for i := 0; i < tService.NumMethod(); i++ {
		mService := tService.Method(i)
		if !successValidateMethod(mService) {
			continue
		}
		nMethod := mService.Name
		handler := func(w http.ResponseWriter, r *http.Request) {
			reqValue := reflect.New(mService.Type.In(2).Elem()).Interface()
			err := json.NewDecoder(r.Body).Decode(reqValue)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			mValue := reflect.ValueOf(service).MethodByName(nMethod)
			args := []reflect.Value{reflect.ValueOf(r.Context()), reflect.ValueOf(reqValue)}
			res := mValue.Call(args)
			rVal := res[0].Interface()
			errRes := res[1].Interface()

			if errRes != nil {
				http.Error(w, errRes.(error).Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(rVal); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		mux.HandleFunc("/"+nMethod, handler)
	}
	return mux
}

func Call(ctx context.Context, endpoint string, method string, req, rsp interface{}) error {
	url := endpoint + "/" + method
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	cl := &http.Client{}
	httpRes, err := cl.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpRes.Body)
		return fmt.Errorf("%s", string(body))
	}

	if err = json.NewDecoder(httpRes.Body).Decode(rsp); err != nil {
		return err
	}
	return nil
}
