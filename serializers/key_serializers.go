package serializers

import (
	"ginserver/models"
)

// KeysSubsetJSON ..
type KeysSubsetJSON struct {
	ListMeta
	Values []models.Key `json:"data"`
}

// NewKeysSubsetJSON ...
func NewKeysSubsetJSON(Keys []models.Key, count int, skip int, total int) KeysSubsetJSON {
	json := KeysSubsetJSON{
		Values: []models.Key{},
		ListMeta: ListMeta{
			Count: count,
			Skip:  skip,
			Total: total,
		},
	}
	for _, Key := range Keys {
		json.Values = append(json.Values, Key)
	}
	return json
}

// SerializeKeys ...
func SerializeKeys(keys []models.Key, params ...interface{}) interface{} {
	count := params[0].(int64)
	skip := params[1].(int64)
	total := params[2].(int)
	keySubsetJSON := NewKeysSubsetJSON(keys, int(count), int(skip), total)
	return NewResponse(0, keySubsetJSON, "Success")
}

// SerializeKey ..
func SerializeKey(key models.Key) interface{} {
	return NewResponse(0, key, "Success")
}
