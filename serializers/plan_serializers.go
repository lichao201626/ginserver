package serializers

import (
	"ginserver/models"
)

// PlansSubsetJSON ..
type PlansSubsetJSON struct {
	ListMeta
	Values []models.Plan `json:"data"`
}

// NewPlansSubsetJSON ...
func NewPlansSubsetJSON(plans []models.Plan, count int, skip int, total int) PlansSubsetJSON {
	json := PlansSubsetJSON{
		Values: []models.Plan{},
		ListMeta: ListMeta{
			Count: count,
			Skip:  skip,
			Total: total,
		},
	}
	for _, plan := range plans {
		json.Values = append(json.Values, plan)
	}
	return json
}

// SerializePlans ...
func SerializePlans(plans []models.Plan, params ...interface{}) interface{} {
	count := params[0].(int64)
	skip := params[1].(int64)
	total := params[2].(int)
	planSubsetJSON := NewPlansSubsetJSON(plans, int(count), int(skip), total)
	return NewResponse(0, planSubsetJSON, "Success")
}

// SerializePlan ..
func SerializePlan(plan models.Plan) interface{} {
	return NewResponse(0, plan, "Success")
}
