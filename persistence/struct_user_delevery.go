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
