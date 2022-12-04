package authorization

import (
	"encoding/json"
	"io"
)

type SessionJsonGet struct {
	SessionName string `json:"sname"`
	Password    string `json:"password"`
}

func UnmarshalJsonGet(r io.Reader) (*SessionJsonGet, error) {
	jsonByte, err := io.ReadAll(r)

	if err != nil {
		return &SessionJsonGet{}, err
	}

	var obj SessionJsonGet
	err = json.Unmarshal(jsonByte, &obj)

	if err != nil {
		return &SessionJsonGet{}, err
	}

	return &obj, nil
}
