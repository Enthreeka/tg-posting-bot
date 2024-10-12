package repo

import (
	"context"
	"errors"
	"github.com/Enthreeka/tg-posting-bot/internal/entity"
	"github.com/Enthreeka/tg-posting-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"time"
)

type PublicationRepo interface {
	CreatePublication(ctx context.Context, publication *entity.Publication) (int, error)

	DeletePublication(ctx context.Context, publicationID int) error

	GetPublicationByPublicationID(ctx context.Context, publicationID int) (*entity.Publication, error)
	GetAllPublicationByChannelID(ctx context.Context, channelID int) ([]entity.Publication, error)
	GetAwaitingPublication(ctx context.Context) ([]*entity.Publication, error)
	GetPublicationAndChannel(ctx context.Context, publicationID int) (*entity.Publication, error)
	GetAllPublicationByID(ctx context.Context, publicationID int) ([]entity.Publication, error)
	GetOnePublicationByID(ctx context.Context, publicationID int) (*entity.Publication, error)
	GetSentAndWaitingToDeletePublication(ctx context.Context) ([]*entity.Publication, error)

	UpdatePublicationButton(ctx context.Context, publicationID int, buttonUrl, buttonText *string) error
	UpdatePublicationText(ctx context.Context, publicationID int, text string) error
	UpdatePublicationStatus(ctx context.Context, publicationID int, status entity.PublicationStatus) error
	UpdatePublicationImage(ctx context.Context, publicationID int, image *string) error
	UpdatePublicationDate(ctx context.Context, publicationID int, date time.Time) error
	UpdateDeleteDate(ctx context.Context, publicationID int, date time.Time) error
	UpdateMessageID(ctx context.Context, publicationID int, messageID int64) error

	IsExistPublication(ctx context.Context, publicationID int) (bool, error)
}

type publicationRepo struct {
	*postgres.Postgres
}

func NewPublicationRepo(pg *postgres.Postgres) (PublicationRepo, error) {
	if pg == nil {
		return nil, errors.New("postgres connection is nil")
	}

	return &publicationRepo{
		pg,
	}, nil
}

func (p *publicationRepo) collectRow(row pgx.Row) (*entity.Publication, error) {
	var publication entity.Publication
	err := row.Scan(&publication.ID,
		&publication.PublicationStatus,
		&publication.PublicationDate,
		&publication.Image,
		&publication.Text,
		&publication.DeleteDate,
		&publication.ChannelID,
		&publication.ButtonUrl,
		&publication.ButtonText)
	if checkErr := ErrorHandler(err); checkErr != nil {
		return nil, checkErr
	}
	return &publication, err
}

func (p *publicationRepo) collectRows(rows pgx.Rows) ([]entity.Publication, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Publication, error) {
		publication, err := p.collectRow(row)
		return *publication, err
	})
}

func (p *publicationRepo) GetPublicationByPublicationID(ctx context.Context, publicationID int) (*entity.Publication, error) {
	query := `select id,publication_status,publication_date,image,text,delete_date,channel_id,button_url,button_text 
				from publication where id = $1`

	row := p.Pool.QueryRow(ctx, query, publicationID)
	return p.collectRow(row)
}

func (p *publicationRepo) CreatePublication(ctx context.Context, publication *entity.Publication) (int, error) {
	query := `insert into publication (channel_id,text,image,publication_date,delete_date,button_url,button_text)
			values ($1,$2,$3,$4,$5,$6,$7) returning id`
	var id int

	err := p.Pool.QueryRow(ctx, query,
		publication.ChannelID,
		publication.Text,
		publication.Image,
		publication.PublicationDate,
		publication.DeleteDate,
		publication.ButtonUrl,
		publication.ButtonText).Scan(&id)
	return id, err
}

func (p *publicationRepo) GetAllPublicationByChannelID(ctx context.Context, channelID int) ([]entity.Publication, error) {
	query := `select id,channel_id,text, publication_date, publication_status from publication where channel_id = $1`
	rows, err := p.Pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	publications := make([]entity.Publication, 0)
	for rows.Next() {
		publication := entity.Publication{}
		err := rows.Scan(&publication.ID,
			&publication.ChannelID,
			&publication.Text,
			&publication.PublicationDate,
			&publication.PublicationStatus)
		if err != nil {
			if checkErr := ErrorHandler(err); checkErr != nil {
				return nil, checkErr
			}
			return nil, err
		}

		publications = append(publications, publication)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return publications, nil
}

func (p *publicationRepo) UpdatePublicationButton(ctx context.Context, publicationID int, buttonUrl, buttonText *string) error {
	query := `update publication set button_url = $1, button_text = $2  where id = $3`
	_, err := p.Pool.Exec(ctx, query, buttonUrl, buttonText, publicationID)
	return err
}

func (p *publicationRepo) UpdatePublicationText(ctx context.Context, publicationID int, text string) error {
	query := `update publication set text = $1 where id = $2`
	_, err := p.Pool.Exec(ctx, query, text, publicationID)
	return err
}

func (p *publicationRepo) UpdatePublicationStatus(ctx context.Context, publicationID int, status entity.PublicationStatus) error {
	query := `update publication set publication_status = $1 where id = $2`
	_, err := p.Pool.Exec(ctx, query, status, publicationID)
	return err
}

func (p *publicationRepo) UpdatePublicationImage(ctx context.Context, publicationID int, image *string) error {
	query := `update publication set image = $1 where id = $2`
	_, err := p.Pool.Exec(ctx, query, image, publicationID)
	return err
}

func (p *publicationRepo) UpdatePublicationDate(ctx context.Context, publicationID int, date time.Time) error {
	query := `update publication set publication_date = $1 where id = $2`
	_, err := p.Pool.Exec(ctx, query, date, publicationID)
	return err
}

func (p *publicationRepo) UpdateDeleteDate(ctx context.Context, publicationID int, date time.Time) error {
	query := `update publication set delete_date = $1 where id = $2`
	_, err := p.Pool.Exec(ctx, query, date, publicationID)
	return err
}

func (p *publicationRepo) DeletePublication(ctx context.Context, publicationID int) error {
	query := `delete from publication where id = $1`
	_, err := p.Pool.Exec(ctx, query, publicationID)
	return err
}

func (p *publicationRepo) IsExistPublication(ctx context.Context, publicationID int) (bool, error) {
	query := `select exists (select id from publication where id = $1)`
	var isExist bool

	err := p.Pool.QueryRow(ctx, query, publicationID).Scan(&isExist)
	return isExist, err
}

func (p *publicationRepo) GetAwaitingPublication(ctx context.Context) ([]*entity.Publication, error) {
	query := `select id, publication_date from publication where publication_status  = 'awaits'`
	rows, err := p.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	publications := make([]*entity.Publication, 0)
	for rows.Next() {
		publication := new(entity.Publication)
		err := rows.Scan(&publication.ID,
			&publication.PublicationDate)
		if err != nil {
			if checkErr := ErrorHandler(err); checkErr != nil {
				return nil, checkErr
			}
			return nil, err
		}

		publications = append(publications, publication)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return publications, nil
}

func (p *publicationRepo) GetPublicationAndChannel(ctx context.Context, publicationID int) (*entity.Publication, error) {
	query := `select c.tg_id,
       				c.channel_name,
					   p.id,
					   p.publication_status,
					   p.publication_date,
					   p.image,
					   p.text,
					   p.delete_date,
					   p.channel_id,
					   p.button_url,
					   p.button_text
				from publication p
				join channel c on p.channel_id = c.id
				where p.id = $1`
	pub := new(entity.Publication)

	err := p.Pool.QueryRow(ctx, query, publicationID).Scan(
		&pub.TelegramChannelID,
		&pub.ChannelName,
		&pub.ID,
		&pub.PublicationStatus,
		&pub.PublicationDate,
		&pub.Image,
		&pub.Text,
		&pub.DeleteDate,
		&pub.ChannelID,
		&pub.ButtonUrl,
		&pub.ButtonText)
	return pub, err
}

func (p *publicationRepo) GetAllPublicationByID(ctx context.Context, publicationID int) ([]entity.Publication, error) {
	query := `select p.id, p.channel_id, p.text, p.publication_date, p.publication_status, c.channel_name
												from publication p
													join channel c on p.channel_id = c.id
												where p.channel_id = (select channel_id from channel с
                                                 join publication p on с.id = p.channel_id
												  where p.id = $1)`
	rows, err := p.Pool.Query(ctx, query, publicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	publications := make([]entity.Publication, 0)
	for rows.Next() {
		publication := entity.Publication{}
		err := rows.Scan(&publication.ID,
			&publication.ChannelID,
			&publication.Text,
			&publication.PublicationDate,
			&publication.PublicationStatus,
			&publication.ChannelName)
		if err != nil {
			if checkErr := ErrorHandler(err); checkErr != nil {
				return nil, checkErr
			}
			return nil, err
		}

		publications = append(publications, publication)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return publications, nil
}

func (p *publicationRepo) GetOnePublicationByID(ctx context.Context, publicationID int) (*entity.Publication, error) {
	query := `select p.id, p.channel_id, p.text, p.publication_date, p.publication_status, c.channel_name
												from publication p
													join channel c on p.channel_id = c.id
												where p.channel_id = (select channel_id from channel с
                                                 join publication p on с.id = p.channel_id
												  where p.id = $1) limit 1`

	publication := entity.Publication{}
	err := p.Pool.QueryRow(ctx, query, publicationID).Scan(
		&publication.ID,
		&publication.ChannelID,
		&publication.Text,
		&publication.PublicationDate,
		&publication.PublicationStatus,
		&publication.ChannelName)
	if err != nil {
		return nil, err
	}

	return &publication, nil
}

func (p *publicationRepo) UpdateMessageID(ctx context.Context, publicationID int, messageID int64) error {
	query := `update publication set message_id = $1 where id = $2`

	_, err := p.Pool.Exec(ctx, query, messageID, publicationID)
	return err
}

func (p *publicationRepo) GetSentAndWaitingToDeletePublication(ctx context.Context) ([]*entity.Publication, error) {
	query := `select id, message_id, delete_date from publication
				where publication_status = 'sent' and delete_date > CURRENT_TIMESTAMP and message_id is not null`
	rows, err := p.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	publications := make([]*entity.Publication, 0)
	for rows.Next() {
		publication := new(entity.Publication)
		err := rows.Scan(&publication.ID,
			&publication.MessageID,
			&publication.DeleteDate,
		)
		if err != nil {
			if checkErr := ErrorHandler(err); checkErr != nil {
				return nil, checkErr
			}
			return nil, err
		}

		publications = append(publications, publication)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return publications, nil

}
