//go:build !solution

package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	getReq := kvapi.GetRequest{Key: key}
	resp, err := c.Get(&getReq)

	var apiErr *kvapi.APIError
	if errors.As(err, &apiErr) {
		err2 := err
		err = apiErr.Unwrap()
		var authErr *kvapi.AuthError
		if errors.As(err, &authErr) {
			return err2
		} else if err != kvapi.ErrKeyNotFound {
			return UpdateValue(c, key, updateFn)
		}
	}
	var newValue string
	oldVersion := uuid.UUID{}
	oldValue := ""

	if resp != nil {
		newValue, err = updateFn(&resp.Value)
		if err != nil {
			return err
		}
		oldVersion = resp.Version
		oldValue = resp.Value
	} else {
		newValue, err = updateFn(nil)
		if err != nil {
			return err
		}
	}

	setReq := kvapi.SetRequest{Key: key, Value: newValue, OldVersion: oldVersion, NewVersion: uuid.Must(uuid.NewV4())}
	return setReqFun(c, &setReq, oldValue, updateFn)
}

func setReqFun(c kvapi.Client, setReq *kvapi.SetRequest, oldValue string, updateFn func(oldValue *string) (newValue string, err error)) error {
	_, err := c.Set(setReq)
	if err != nil {
		var apiErr *kvapi.APIError
		if errors.As(err, &apiErr) {
			err2 := err
			err = apiErr.Unwrap()
			var authErr *kvapi.AuthError
			var confErr *kvapi.ConflictError
			if errors.As(err, &authErr) {
				return err2
			} else if errors.As(err, &confErr) {
				if confErr.ExpectedVersion == setReq.NewVersion {
					return nil
				}
				getReq := kvapi.GetRequest{Key: setReq.Key}
				resp, _ := c.Get(&getReq)
				newValue, _ := updateFn(&resp.Value)
				if oldValue != newValue {
					setReq.Value = newValue
				}
				setReq.OldVersion = confErr.ExpectedVersion
				return setReqFun(c, setReq, oldValue, updateFn)
			} else if err == kvapi.ErrKeyNotFound {
				newValue, err := updateFn(nil)
				if err != nil {
					return err
				}
				setReq.Value = newValue
				setReq.NewVersion = uuid.UUID{}
				setReq.OldVersion = uuid.UUID{}
				return setReqFun(c, setReq, oldValue, updateFn)
			} else {
				return setReqFun(c, setReq, oldValue, updateFn)
			}
		}
	}
	return nil
}
