package scheduled

import (
	"context"
	"errors"
	"github.com/Enthreeka/tg-posting-bot/internal/entity"
	service "github.com/Enthreeka/tg-posting-bot/internal/usecase"
	store "github.com/Enthreeka/tg-posting-bot/pkg/local_storage"
	"github.com/Enthreeka/tg-posting-bot/pkg/logger"
	customMsg "github.com/Enthreeka/tg-posting-bot/pkg/tg_bot_api"
	"time"
)

type Schedule interface {
	LoadDatabaseInPubDelStore(ctx context.Context) error
	Start(ctx context.Context) error
}

type schedule struct {
	publicationService service.PublicationService
	tgMsg              customMsg.Message
	pubStore           *store.PublicationArray
	log                *logger.Logger
}

func NewSchedule(publicationService service.PublicationService,
	tgMsg customMsg.Message,
	pubStore *store.PublicationArray,
	log *logger.Logger) (Schedule, error) {
	if tgMsg == nil {
		return nil, errors.New("tgMsg cannot be nil")
	}
	if pubStore == nil {
		return nil, errors.New("pubStore cannot be nil")
	}
	if publicationService == nil {
		return nil, errors.New("publicationService cannot be nil")
	}
	if log == nil {
		return nil, errors.New("log cannot be nil")
	}

	return &schedule{
		tgMsg:              tgMsg,
		pubStore:           pubStore,
		log:                log,
		publicationService: publicationService,
	}, nil
}

func (s *schedule) LoadDatabaseInPubDelStore(ctx context.Context) error {
	publications, err := s.publicationService.GetAwaitingPublication(ctx)
	if err != nil {
		s.log.Error("Failed to get pub publications", "err", err)
		return err
	}

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		s.log.Error("Failed to load time location", "err", err)
		return err
	}

	var countPubs int
	for _, value := range publications {
		if value.PublicationDate.Before(time.Now().In(loc)) {
			go func() {
				if err := s.publicationService.UpdatePublicationStatus(ctx, value.ID, entity.StatusErrorOnSending); err != nil {
					s.log.Error("Failed to update publication status err: ", err)
				}
			}()
			continue
		}
		s.pubStore.AppendPub(&store.PubData{
			PublicationID: value.ID,
			PubDate:       value.PublicationDate,
		})
		countPubs++
	}
	s.log.Info("Successfully loaded publications count - %d", countPubs)

	delPublications, err := s.publicationService.GetSentAndWaitingToDeletePublication(ctx)
	if err != nil {
		s.log.Error("Failed to get del publications", "err", err)
		return err
	}

	var countDelPubs int
	for _, value := range delPublications {
		if value.DeleteDate.Before(time.Now().In(loc)) {
			go func() {
				if err := s.publicationService.UpdatePublicationStatus(ctx, value.ID, entity.StatusErrorOnDeleting); err != nil {
					s.log.Error("Failed to update publication status err: ", err)
				}
			}()
			continue
		}
		s.pubStore.AppendDel(&store.PubData{
			PublicationID: value.ID,
			DelDate:       *value.DeleteDate,
			SentMsgID:     int(value.MessageID),
		})
		countDelPubs++
	}

	s.log.Info("Successfully loaded del publications count - %d", countDelPubs)

	return nil
}

func (s *schedule) Start(ctx context.Context) error {
	timeTicker := time.NewTicker(time.Minute)
	defer func() {
		timeTicker.Stop()
		s.log.Info("Scheduler stopped")
	}()

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}

	for {
		select {
		case <-timeTicker.C:
			s.pubStore.SortPub()
			for _, value := range s.pubStore.GetPub() {
				if value == nil {
					continue
				}
				now := time.Now().In(loc).Round(time.Minute)
				pubDate := value.PubDate.Round(time.Minute)
				if now.Equal(pubDate) {
					s.pubStore.RemovePub(value)

					go func(pubID int) {
						publication, err := s.publicationService.GetPublicationAndChannel(ctx, pubID)
						if err != nil {
							s.log.Error("Failed to get publication by publicationID: publicationID - %d, err - %v", pubID, err)
						}

						if publication != nil && publication.PublicationStatus != entity.StatusSent {
							var status entity.PublicationStatus
							msgID, err := s.tgMsg.SendMessageToUser(publication.TelegramChannelID, publication)
							if err != nil {
								s.log.Error("Failed to send message to channel - %d, err - %v", publication.ChannelID, err)
								status = entity.StatusErrorOnSending
							}
							if err == nil {
								status = entity.StatusSent
								if publication.DeleteDate != nil {
									s.pubStore.AppendDel(&store.PubData{
										PublicationID: publication.ID,
										DelDate:       *publication.DeleteDate,
										SentMsgID:     msgID,
										ChannelID:     publication.TelegramChannelID,
									})

									if err := s.publicationService.UpdateMessageID(ctx, publication.ID, int64(msgID)); err != nil {
										s.log.Error("Failed to update message id - %d, err - %v", publication.ID, err)
									}
								}
							}

							if err = s.publicationService.UpdatePublicationStatus(ctx, publication.ID, status); err != nil {
								s.log.Error("Failed to update publication, publicationID - %d, status - %v, err - %v", pubID, status, err)
							}

							s.log.Info("Sent publication for publicationID: %d, channel_id: %d, msg_id: %d",
								pubID, publication.TelegramChannelID, msgID)
						}
					}(value.PublicationID)
				}
			}

			for _, value := range s.pubStore.GetDel() {
				if value == nil {
					continue
				}
				if time.Now().In(loc).Round(time.Minute).Equal(value.DelDate.Round(time.Minute)) {
					s.pubStore.ReplaceDel(value)

					go func(channelID int64, sentMsgID int, pubID int) {
						if channelID == 0 {
							publication, err := s.publicationService.GetPublicationAndChannel(ctx, pubID)
							if err != nil {
								s.log.Error("Failed to get publication by publicationID: publicationID - %d, err - %v", pubID, err)
							}
							channelID = publication.TelegramChannelID
						}

						var status entity.PublicationStatus
						if err = s.tgMsg.DeleteMessage(channelID, sentMsgID); err != nil {
							s.log.Error("Failed to delete message from channel - %d, sentMsgID -%d, err - %v", channelID, sentMsgID, err)
							status = entity.StatusErrorOnDeleting
						}
						if err == nil {
							status = entity.StatusDeletedByBot
						}

						if err = s.publicationService.UpdatePublicationStatus(ctx, pubID, status); err != nil {
							s.log.Error("Failed to update publication, publicationID - %d, status - %v, err - %v", pubID, status, err)
						}

						s.log.Info("Deleted publication for publicationID %d", pubID)
					}(value.ChannelID, value.SentMsgID, value.PublicationID)
				}
			}

		case <-ctx.Done():
			s.log.Error("context canceled")
			return ctx.Err()
		}
	}
}
