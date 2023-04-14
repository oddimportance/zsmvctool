package persistence

type UserAppDetails struct {
	Id           string
	First_name   string
	Last_name    string
	Email        string
	Blocked      string
	Crdate       string
	Telephone    string
	Password     string
	Mobile_phone string

	Role string
}

type UserAppAddress struct {
	User_id            string
	Address            string
	Address_additional string
	City               string
	Postal_code        string
	Country            string
}

var UserKeyFromJson UserAppDetails = UserAppDetails{
	First_name:   "first_name",
	Last_name:    "last_name",
	Email:        "email",
	Password:     "password",
	Mobile_phone: "mobile_phone",
	Telephone:    "telephone",
}

var AddressKeyFromJson UserAppAddress = UserAppAddress{
	User_id:            "user_id",
	Address:            "address",
	Address_additional: "address_additional",
	City:               "city",
	Postal_code:        "postal_code",
	Country:            "country",
}
