package persistence

type OrderDetails struct {
	OrderDetails struct {
		Orderinfo struct {
			Amount       string `json:"Amount"`
			Address_id   string `json:"Address_id"`
			Resturant_id string `json:"Resturant_id"`
			User_id      string `json:"User_id"`
			Status       string `json:"Status"`
			Mobile_phone string `json:"Mobile_phone"`
		} `json:"orderinfo"`
		Orderpayment struct {
			ID                   string `json:"Id"`
			Order_id             string `json:"Order_id"`
			Canceled_at          string `json:"Canceled_at"`
			Capture_method       string `json:"Capture_method"`
			Amount               string `json:"Amount"`
			Client_secret        string `json:"Client_secret"`
			Currency             string `json:"Currency"`
			Created              string `json:"Created"`
			Confirmation_method  string `json:"Confirmation_method"`
			Payment_method_types string `json:"Payment_method_types"`
			Cancellation_reason  string `json:"Cancellation_reason"`
			Status               string `json:"Status"`
		} `json:"orderpayment"`
		Items []struct {
			Menu_item_id string `json:"Menu_item_id"`
			Order_id     string `json:"Order_id"`
			Price        string `json:"Price"`
			Extras       []struct {
				Menu_item_id string `json:"Menu_item_id"`
				Price        string `json:"Price"`
			} `json:"extras"`
		} `json:"items"`
	} `json:"orderDetails"`
}
