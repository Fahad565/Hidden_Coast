package models

type SME struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"` // "food", "crafts", "experience"
	Location    Location  `json:"location"`
	Description string    `json:"description"`
	Story       string    `json:"story"` // the "3 generations" narrative
	Languages   []string  `json:"languages"`
	Opening     string    `json:"opening"`  // "10:00"
	Closing     string    `json:"closing"`  // "18:00"
	Capacity    int       `json:"capacity"` // max people at once
	Products    []Product `json:"products"`
}

type Location struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Area     string  `json:"area"`      // "Old Town", "Marikiti", etc.
	WalkFrom string  `json:"walk_from"` // e.g. "Mombasa Port"
	WalkMins int     `json:"walk_mins"`
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`        // KSh
	PrepMinutes int    `json:"prep_minutes"` // 0 for retail, 35 for food
	Stock       int    `json:"stock"`        // remaining today
	MaxStock    int    `json:"max_stock"`    // capacity today
	Custom      bool   `json:"custom"`       // made-to-order?
}

func (s *SME) FindProduct(productID string) *Product {
	for i := range s.Products {
		if s.Products[i].ID == productID {
			return &s.Products[i]
		}
	}
	return nil
}
