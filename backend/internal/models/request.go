package models

type TouristRequest struct {
	ID        string `json:"id"`
	SMEID     string `json:"sme_id"`
	TouristID string `json:"tourist_id"`
	Items     []Item `json:"items"`
	PickupAt  string `json:"pickup_at"` // "13:15"
	Total     int    `json:"total"`
	Status    string `json:"status"` // "pending", "accepted", "declined"
}

type Item struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Qty       int    `json:"qty"`
	Price     int    `json:"price"`
}
