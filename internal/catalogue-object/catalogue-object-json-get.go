package catalogueobject

import (
	"encoding/json"
	"io"
)

type CatalogueObjectJsonGet struct {
	Name string `json:"name"`
}

func UnmarshalJsonGet(r io.Reader) (*CatalogueObjectJsonGet, error) {
	jsonByte, err := io.ReadAll(r)

	if err != nil {
		return &CatalogueObjectJsonGet{}, err
	}

	var obj CatalogueObjectJsonGet
	err = json.Unmarshal(jsonByte, &obj)

	if err != nil {
		return &CatalogueObjectJsonGet{}, err
	}

	return &obj, nil
}
