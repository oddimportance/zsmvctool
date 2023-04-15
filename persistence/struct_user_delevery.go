package persistence

type UserDeleveryagentMobilePhone struct {
	Mobile_phone string
}

type UserDeleveryagentCodeVerification struct {
	Verification_code string
	Resturant_id      string
	Mobile_phone      string
}
type UserQrCode struct {
	User_qr_code string
}

type ReusableQrCode struct {
	Reusable_qr_code string
}

type ReusableMarkReturned struct {
	Reusable_qr_code string
	Delevery_id      string
}

var ReusableMarkReturnedchema map[int]JsonForm = map[int]JsonForm{
	1: JsonForm{
		Key:         "reusable_qr_code",
		Value:       "",
		Max:         100,
		Min:         0,
		IsMandatory: true,
	},
	2: JsonForm{
		Key:         "delevery_id",
		Value:       "",
		Max:         20,
		Min:         0,
		IsMandatory: true,
	},
}

var ReusableQRCodechema map[int]JsonForm = map[int]JsonForm{
	1: JsonForm{
		Key:         "reusable_qr_code",
		Value:       "",
		Max:         100,
		Min:         0,
		IsMandatory: true,
	},
}

var UserQRCodechema map[int]JsonForm = map[int]JsonForm{
	1: JsonForm{
		Key:         "user_qr_code",
		Value:       "",
		Max:         100,
		Min:         0,
		IsMandatory: true,
	},
}
