package persistence

type Reusabletracking struct {
	UserID       string         `json:"User_id"`
	DeleveryID   string         `json:"Delevery_id"`
	ReusableCode []ReusableCode `json:"Reusable_code"`
}

type ReusableCode struct {
	ResturentID  string `json:"Resturent_id"`
	ReusableCode string `json:"Reusable_code"`
}
