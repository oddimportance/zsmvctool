package persistence

// User details
// @ param First_name 			string
// @ param Last_name  			string
// @ param Email      			string
// @ param Salutation 			string
// @ param Title      			string
// @ param Avatar				string
// @ param Language				string
// @ param Role		string
type UserDetails struct {
	Id           string
	Title        string
	Salutation   string
	First_name   string
	Last_name    string
	Email        string
	Avatar       string
	Language     string
	Telephone    string
	Mobile_phone string
	RestaurantId string
	CountryCode  string
}

var UserschemaRegister map[int]JsonForm = map[int]JsonForm{
	1: JsonForm{
		Key:         "first_name",
		Value:       "",
		Max:         20,
		Min:         3,
		IsMandatory: true,
		//Type:        "alphabetOnly",
	},
	2: JsonForm{
		Key:         "last_name",
		Value:       "",
		Max:         20,
		Min:         3,
		IsMandatory: true,
		//Type:        "alphabetOnly",
	},
	3: JsonForm{
		Key:         "mobile_phone",
		Value:       "",
		Max:         20,
		Min:         10,
		IsMandatory: true,
	},
}

var UserschemaLogin map[int]JsonForm = map[int]JsonForm{

	1: JsonForm{
		Key:         "mobile_phone",
		Value:       "",
		Max:         255,
		Min:         10,
		IsMandatory: true,
	},

	2: JsonForm{
		Key:         "password",
		Value:       "",
		Max:         20,
		Min:         8,
		IsMandatory: true,
	},
}

var UserschemaVerifyMobilenumber map[int]JsonForm = map[int]JsonForm{

	1: JsonForm{
		Key:         "mobile_phone",
		Value:       "",
		Max:         255,
		Min:         10,
		IsMandatory: true,
	},
	2: JsonForm{
		Key:         "verificationCode",
		Value:       "",
		Max:         6,
		Min:         6,
		IsMandatory: true,
	},
}

var UserschemaAddress map[int]JsonForm = map[int]JsonForm{
	1: JsonForm{
		Key:         "user_id",
		Value:       "",
		Max:         20,
		Min:         0,
		IsMandatory: true,
	},
	2: JsonForm{
		Key:         "address",
		Value:       "",
		Max:         20,
		Min:         0,
		IsMandatory: true,
	},
	3: JsonForm{
		Key:         "address_additional",
		Value:       "",
		Max:         100,
		Min:         0,
		IsMandatory: true,
	},

	4: JsonForm{
		Key:         "city",
		Value:       "",
		Max:         100,
		Min:         3,
		IsMandatory: true,
	},
	5: JsonForm{
		Key:         "postal_code",
		Value:       "",
		Max:         15,
		Min:         0,
		IsMandatory: true,
	},
	6: JsonForm{
		Key:         "country",
		Value:       "",
		Max:         20,
		Min:         0,
		IsMandatory: true,
	},
}
