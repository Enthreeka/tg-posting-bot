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
	PublicationButtonUpdate     TypeCommand = "update_publication_button"
	PublicationSentDateUpdate   TypeCommand = "update_publication_sent_date"
	PublicationDeleteDateUpdate TypeCommand = "update_publication_delete_date"
)

var MapTypes = map[TypeCommand]OperationType{
	AdminCreate:       Admin,
	AdminDelete:       Admin,
	PublicationCreate: Publication,
	PublicationDelete: Publication,
}
