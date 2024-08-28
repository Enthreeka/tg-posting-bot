package repo

import (
	"context"
	"errors"
	"github.com/Enthreeka/tg-posting-bot/internal/entity"
	"github.com/Enthreeka/tg-posting-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type ChannelRepo interface {
	Create(ctx context.Context, channel *entity.Channel) error
	GetByID(ctx context.Context, id int) (*entity.Channel, error)
	DeleteByID(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]entity.Channel, error)
	UpdateStatusByTgID(ctx context.Context, status entity.ChannelStatus, telegramID int64) error
	IsChannelExistByTgID(ctx context.Context, telegramID int64) (bool, error)
	GetAllAdminChannel(ctx context.Context) ([]entity.Channel, error)
	GetChannelIDByChannelName(ctx context.Context, channelName string) (int64, error)
	GetByChannelName(ctx context.Context, channelName string) (*entity.Channel, error)
	//GetChannelByUserID(ctx context.Context, userID int64) (string, error)
}

type channelRepo struct {
	*postgres.Postgres
}

func NewChannelRepo(pg *postgres.Postgres) (ChannelRepo, error) {
	if pg == nil {
		return nil, errors.New("postgres connection is nil")
	}

	return &channelRepo{
		pg,
	}, nil
}

func (u *channelRepo) collectRow(row pgx.Row) (*entity.Channel, error) {
	var channel entity.Channel
	err := row.Scan(&channel.ID, &channel.TgID, &channel.ChannelName, &channel.ChannelUrl, &channel.ChannelStatus)
	if checkErr := ErrorHandler(err); checkErr != nil {
		return nil, checkErr
	}
	return &channel, err
}

func (u *channelRepo) collectRows(rows pgx.Rows) ([]entity.Channel, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Channel, error) {
		channel, err := u.collectRow(row)
		return *channel, err
	})
}

func (u *channelRepo) Create(ctx context.Context, channel *entity.Channel) error {
	query := `insert into channel (tg_id,channel_name,channel_url,channel_status) values ($1,$2,$3,$4)`

	_, err := u.Pool.Exec(ctx, query, channel.TgID, channel.ChannelName, channel.ChannelUrl, channel.ChannelStatus)
	return err
}

func (u *channelRepo) GetByID(ctx context.Context, id int) (*entity.Channel, error) {
	query := `select * from channel where id = $1`

	row := u.Pool.QueryRow(ctx, query, id)
	return u.collectRow(row)
}

func (u *channelRepo) DeleteByID(ctx context.Context, id int) error {
	query := `delete from channel where id = $1`

	_, err := u.Pool.Exec(ctx, query, id)
	return err
}

func (u *channelRepo) GetAll(ctx context.Context) ([]entity.Channel, error) {
	query := `select * from channel`

	rows, err := u.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return u.collectRows(rows)
}

func (u *channelRepo) UpdateStatusByTgID(ctx context.Context, status entity.ChannelStatus, telegramID int64) error {
	query := `update channel set channel_status = $1 where tg_id = $2`

	_, err := u.Pool.Exec(ctx, query, status, telegramID)
	if checkErr := ErrorHandler(err); checkErr != nil {
		return checkErr
	}

	return err
}

func (u *channelRepo) IsChannelExistByTgID(ctx context.Context, telegramID int64) (bool, error) {
	query := `select exists (select id from channel where tg_id = $1)`
	var isExist bool

	err := u.Pool.QueryRow(ctx, query, telegramID).Scan(&isExist)
	if checkErr := ErrorHandler(err); checkErr != nil {
		return isExist, checkErr
	}

	return isExist, err
}

func (u *channelRepo) GetAllAdminChannel(ctx context.Context) ([]entity.Channel, error) {
	query := `select * from channel where channel_status = 'administrator'`

	rows, err := u.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return u.collectRows(rows)
}

func (u *channelRepo) GetChannelIDByChannelName(ctx context.Context, channelName string) (int64, error) {
	query := `select tg_id from channel where channel_name = $1`
	var ChannelTelegramID int64

	err := u.Pool.QueryRow(ctx, query, channelName).Scan(&ChannelTelegramID)
	if checkErr := ErrorHandler(err); checkErr != nil {
		return 0, checkErr
	}

	return ChannelTelegramID, err
}

func (u *channelRepo) GetByChannelName(ctx context.Context, channelName string) (*entity.Channel, error) {
	query := `select * from channel where channel_name = $1`
	channel := new(entity.Channel)

	err := u.Pool.QueryRow(ctx, query, channelName).Scan(
		&channel.ID,
		&channel.ChannelName,
		&channel.ChannelUrl,
		&channel.ChannelStatus,
		&channel.TgID,
	)
	if checkErr := ErrorHandler(err); checkErr != nil {
		return channel, checkErr
	}

	return channel, err
}

//func (u *channelRepo) GetChannelByUserID(ctx context.Context, userID int64) (string, error) {
//	query := `select channel_name from channel
//				join user_channel on  user_channel.channel_tg_id = channel.tg_id
//				where user_channel.user_id = $1 `
//
//	var channel string
//
//	err := u.Pool.QueryRow(ctx, query, userID).Scan(&channel)
//	if checkErr := ErrorHandler(err); checkErr != nil {
//		return "", checkErr
//	}
//
//	return channel, err
//}
