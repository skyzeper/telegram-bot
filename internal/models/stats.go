package models

// Stats represents statistics data
type Stats struct {
	TotalOrders        int     `json:"total_orders"`
	WasteRemovalOrders int     `json:"waste_removal_orders"`
	DemolitionOrders   int     `json:"demolition_orders"`
	ConstructionOrders int     `json:"construction_orders"`
	TotalAmount        float64 `json:"total_amount"`
	DriverDebts        float64 `json:"driver_debts"`
}