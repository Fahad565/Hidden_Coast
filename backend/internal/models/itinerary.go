package models

type Itinerary struct {
	ID        string `json:"id"`
	TouristID string `json:"tourist_id"`
	Stops     []Stop `json:"stops"`
	TotalCost int    `json:"total_cost"`
	Budget    int    `json:"budget"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type Stop struct {
	Time      string `json:"time"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	Cost      int    `json:"cost"`
	Duration  int    `json:"duration_mins"`
	SMEID     string `json:"sme_id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Note      string `json:"note"`
}
