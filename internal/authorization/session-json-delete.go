package authorization

import (
	"encoding/json"
	"io"
)

type SessionJsonDelete struct {
	Id int `json:"id"`
}

func UnmarshalJsonDelete(r io.Reader) (*SessionJsonDelete, error) {
	jsonByte, err := io.ReadAll(r)

	if err != nil {
		return &SessionJsonDelete{}, err
	}

	var obj SessionJsonDelete
	err = json.Unmarshal(jsonByte, &obj)

	if err != nil {
		return &SessionJsonDelete{}, err
	}

	return &obj, nil
}
