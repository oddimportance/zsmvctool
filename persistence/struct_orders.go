package persistence

type Orders struct {
	//Id           string
	Resturant_id string
	User_id      string
	Mobile_phone string
	Status       string
	Address_id   string
}

type OrderItems struct {
	Order_id     string
	Menu_item_id string
	Price        string
	Parent_item  string
}

type OrderPayments struct {
	Id                   string
	Order_id             string
	Canceled_at          string
	Capture_method       string
	Amount               string
	Client_secret        string
	Currency             string
	Created              string
	Confirmation_method  string
	Payment_method_types string
	Cancellation_reason  string
	Status               string
}
