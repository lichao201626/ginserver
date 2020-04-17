package serializers

// ListMeta ...
type ListMeta struct {
	Count int `json:"count"`
	Skip  int `json:"skip"`
	Total int `json:"total"`
}
