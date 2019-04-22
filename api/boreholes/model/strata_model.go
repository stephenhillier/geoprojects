package model

// Strata is a soil layer/stratum and contains information such as description and depth of the layer
type Strata struct {
	ID          int64   `json:"id"`
	Borehole    int64   `json:"borehole"`
	Start       float64 `json:"start" db:"start_depth"`
	End         float64 `json:"end" db:"end_depth"`
	Description string  `json:"description"`
	Soils       string  `json:"soils"`
	Moisture    string  `json:"moisture"`
	Consistency string  `json:"consistency"`
}

// StrataRequest is a struct containing fields required to create a new strata layer
type StrataRequest struct {
	Borehole    int64   `json:"borehole,string"`
	Start       float64 `json:"start,string"`
	End         float64 `json:"end,string"`
	Description string  `json:"description"`
}
