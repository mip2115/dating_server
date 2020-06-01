package types

type Location struct {
	Name    *string   `json:"name"`
	Address *string   `json:"address"`
	City    *string   `json:"city"`
	Zipcode *string   `json:"zipcode"`
	Hours   []*string `json:"hours"` // todo – make these structs
	Dates   []*string `json:"dates"`
	Orders  []*string `json:"orders"`
}
