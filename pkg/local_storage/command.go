package store

type OperationType string

type TypeCommand string

const (
	Admin       OperationType = "admin"
	Publication OperationType = "Publication"
)

const (
	AdminCreate TypeCommand = "create_role"
	AdminDelete TypeCommand = "delete_role"

	PublicationCreate           TypeCommand = "create_publication"
	PublicationDelete           TypeCommand = "delete_publication"
	PublicationTextUpdate       TypeCommand = "update_publication_text"
	PublicationImageUpdate      TypeCommand = "update_publication_image"
	PublicationButtonTextUpdate TypeCommand = "update_publication_button_text"
	PublicationSentDateUpdate   TypeCommand = "update_publication_sent_date"
	PublicationDeleteDateUpdate TypeCommand = "update_publication_delete_date"
	PublicationButtonLinkUpdate TypeCommand = "update_publication_button_link"
)

var MapTypes = map[TypeCommand]OperationType{
	AdminCreate:       Admin,
	AdminDelete:       Admin,
	PublicationCreate: Publication,
	PublicationDelete: Publication,
}
