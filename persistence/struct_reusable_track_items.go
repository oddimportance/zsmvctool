package persistence

type ReusableTrackItems struct {
	UserID       string         `json:"user_id"`
	DeleveryAgentID   string         `json:"delevery_agent_id"`
	ReusableCodes []string `json:"reusable_codes"`
}
