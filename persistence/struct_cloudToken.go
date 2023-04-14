package persistence

type CloudToken struct {
	PortalKey    string // sapp = super admins, app = restaurant manager, api = mobile users
	Destination  string
	FileName     string
	FileSize     string
	FileMimeType string
	CreatedBy    string
	Token        string
	CloudUrl     string
}

type CloudAction string

const (
	upload CloudAction = "upload"
	delete CloudAction = "delete"
)

func (c CloudAction) Upload() CloudAction {
	return upload
}

func (c CloudAction) Delete() CloudAction {
	return delete
}
