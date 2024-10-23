package bot

import (
	"context"
	"github.com/Enthreeka/tg-posting-bot/internal/config"
	"github.com/Enthreeka/tg-posting-bot/internal/handler/callback"
	"github.com/Enthreeka/tg-posting-bot/internal/handler/middleware"
	"github.com/Enthreeka/tg-posting-bot/internal/handler/tgbot"
	"github.com/Enthreeka/tg-posting-bot/internal/handler/view"
	"github.com/Enthreeka/tg-posting-bot/internal/repo"
	"github.com/Enthreeka/tg-posting-bot/internal/scheduled"
	"github.com/Enthreeka/tg-posting-bot/internal/usecase"
	store "github.com/Enthreeka/tg-posting-bot/pkg/local_storage"
	"github.com/Enthreeka/tg-posting-bot/pkg/logger"
	"github.com/Enthreeka/tg-posting-bot/pkg/postgres"
	customMsg "github.com/Enthreeka/tg-posting-bot/pkg/tg_bot_api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	PostgresMaxAttempts = 5
)

type Bot struct {
	bot              *tgbotapi.BotAPI
	psql             *postgres.Postgres
	store            *store.Store
	cfg              *config.Config
	log              *logger.Logger
	tgMsg            *customMsg.TelegramMsg
	callbackStore    *store.CallbackStorage
	publicationArray *store.PublicationArray

	userService        service.UserService
	channelService     service.ChannelService
	publicationService service.PublicationService

	publicationSchedule scheduled.Schedule

	userRepo        repo.UserRepo
	channelRepo     repo.ChannelRepo
	publicationRepo repo.PublicationRepo

	callbackUser        callback.CallbackUser
	callbackChannel     callback.CallbackChannel
	callbackPublication callback.PublicationChannel

	viewGeneral *view.ViewGeneral
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) initHandler() {
	b.viewGeneral = view.NewViewGeneral(b.log, b.tgMsg)

	callbackUser, err := callback.NewCallbackUser(b.userService, b.log, b.store, b.tgMsg)
	if err != nil {
		b.log.Fatal("NewCallbackUser: ", err)
	}
	b.callbackUser = callbackUser

	callbackChannel, err := callback.NewCallbackChannel(b.channelService, b.publicationService, b.log, b.tgMsg, b.store)
	if err != nil {
		b.log.Fatal("NewCallbackChannel: ", err)
	}
	b.callbackChannel = callbackChannel

	callbackPublication, err := callback.NewCallbackPublication(b.publicationService, b.log, b.tgMsg, b.store, b.channelService, b.publicationArray)
	if err != nil {
		b.log.Fatal("callbackPublication: ", err)
	}
	b.callbackPublication = callbackPublication

	b.log.Info("Initializing handler")
}

func (b *Bot) initUsecase() {
	userService, err := service.NewUserService(b.userRepo, b.log)
	if err != nil {
		b.log.Fatal("NewUserService: ", err)
	}
	b.userService = userService

	channelService, err := service.NewChannelService(b.channelRepo, b.log)
	if err != nil {
		b.log.Fatal("NewChannelService:", err)
	}
	b.channelService = channelService

	publicationService, err := service.NewPublicationService(b.publicationRepo, b.log)
	if err != nil {
		b.log.Fatal("NewPublicationService:", err)
	}
	b.publicationService = publicationService

	b.log.Info("Initializing usecase")
}

func (b *Bot) initRepo() {
	userRepo, err := repo.NewUserRepo(b.psql)
	if err != nil {
		b.log.Fatal("NewUserRepo: ", err)
	}

	b.userRepo = userRepo

	channelRepo, err := repo.NewChannelRepo(b.psql)
	if err != nil {
		b.log.Fatal("NewChannelRepo: ", err)
	}
	b.channelRepo = channelRepo

	publicationRepo, err := repo.NewPublicationRepo(b.psql)
	if err != nil {
		b.log.Fatal("NewPublicationRepo: ", err)
	}
	b.publicationRepo = publicationRepo

	b.log.Info("Initializing repo")
}

func (b *Bot) initMessage() {
	b.tgMsg = customMsg.NewMessageSetting(b.bot, b.log)

	b.log.Info("Initializing message")
}

func (b *Bot) initPostgres(ctx context.Context) {
	psql, err := postgres.New(ctx, PostgresMaxAttempts, b.cfg.Postgres.URL)
	if err != nil {
		b.log.Fatal("failed to connect PostgreSQL: %v", err)
	}
	b.psql = psql

	b.log.Info("Initializing postgres")
}

func (b *Bot) initConfig() {
	cfg, err := config.New()
	if err != nil {
		b.log.Fatal("failed load config: %v", err)
	}
	b.cfg = cfg

	b.log.Info("Initializing config")
}

func (b *Bot) initLogger() {
	b.log = logger.New()

	b.log.Info("Initializing logger")
}

func (b *Bot) initStore() {
	b.store = store.NewStore()

	b.log.Info("Initializing store")
}

func (b *Bot) initCallbackStorage() {
	b.callbackStore = store.NewCallbackStorage()

	b.log.Info("Initializing callback storage")
}

func (b *Bot) initTelegramBot() {
	bot, err := tgbotapi.NewBotAPI(b.cfg.Telegram.Token)
	if err != nil {
		b.log.Fatal("failed to load token %v", err)
	}
	bot.Debug = false
	b.bot = bot

	b.log.Info("Initializing telegram bot")
	b.log.Info("Authorized on account %s", bot.Self.UserName)
}

func (b *Bot) initScheduled(ctx context.Context) {
	publicationSchedule, err := scheduled.NewSchedule(b.publicationService, b.tgMsg, b.publicationArray, b.log)
	if err != nil {
		b.log.Fatal("NewSchedule: %v", err)
	}
	b.publicationSchedule = publicationSchedule

	if err = publicationSchedule.LoadDatabaseInPubDelStore(ctx); err != nil {
		b.log.Fatal("LoadDatabaseInPubDelStore: %v", err)
	}
	go publicationSchedule.StartPub(ctx)
	go publicationSchedule.StartDel(ctx)

	b.log.Info("Initializing scheduled")
}

func (b *Bot) initStoreScheduled() {
	b.publicationArray = store.NewSortPublication(200)
	b.log.Info("Initializing store scheduled")
}

func (b *Bot) initialize(ctx context.Context) {
	b.initLogger()
	b.initConfig()
	b.initTelegramBot()
	b.initStore()
	b.initStoreScheduled()
	b.initCallbackStorage()
	b.initPostgres(ctx)
	b.initMessage()
	b.initRepo()
	b.initUsecase()
	b.initHandler()
	b.initScheduled(ctx)
}

func (b *Bot) Run(ctx context.Context) {
	startBot := time.Now()
	b.initialize(ctx)
	newBot, err := tgbot.NewBot(b.bot, b.log, b.store, b.tgMsg, b.userService, b.channelService, b.publicationService, b.callbackStore, b.publicationArray)
	if err != nil {
		b.log.Fatal("failed go create new bot: ", err)
	}
	defer b.psql.Close()

	newBot.RegisterCommandView("secret", middleware.AdminMiddleware(b.userService, b.viewGeneral.CallbackStartAdminPanel()))

	// user domain
	newBot.RegisterCommandCallback("main_menu", middleware.AdminMiddleware(b.userService, b.callbackUser.MainMenu()))
	newBot.RegisterCommandCallback("user_setting", middleware.AdminMiddleware(b.userService, b.callbackUser.AdminRoleSetting()))
	newBot.RegisterCommandCallback("admin_look_up", middleware.AdminMiddleware(b.userService, b.callbackUser.AdminLookUp()))
	newBot.RegisterCommandCallback("admin_delete_role", middleware.AdminMiddleware(b.userService, b.callbackUser.AdminDeleteRole()))
	newBot.RegisterCommandCallback("admin_set_role", middleware.AdminMiddleware(b.userService, b.callbackUser.AdminSetRole()))

	// channel domain
	newBot.RegisterCommandCallback("show_channels", middleware.AdminMiddleware(b.userService, b.callbackChannel.CallbackShowAllChannels()))
	newBot.RegisterCommandCallback("channel_get", middleware.AdminMiddleware(b.userService, b.callbackChannel.CallbackGetChannel()))
	newBot.RegisterCommandCallback("back_setting", middleware.AdminMiddleware(b.userService, b.callbackChannel.CallbackGetChannel()))
	newBot.RegisterCommandCallback("cancel_create", middleware.AdminMiddleware(b.userService, b.callbackChannel.CallbackCancelCreate()))

	// publication domain
	newBot.RegisterCommandCallback("publication_create", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackCreatePublication()))
	newBot.RegisterCommandCallback("publication_update", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackUpdatePublicationSettings()))
	newBot.RegisterCommandCallback("publication_get", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackGetPublicationGet()))
	newBot.RegisterCommandCallback("text_update", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackUpdatePublicationText()))
	newBot.RegisterCommandCallback("image_update", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackUpdatePublicationImage()))
	newBot.RegisterCommandCallback("buttontext_update", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackUpdatePublicationButtonText()))
	newBot.RegisterCommandCallback("buttonlink_update", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackUpdatePublicationButtonLink()))
	newBot.RegisterCommandCallback("sent-date_update", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackUpdatePublicationSentDate()))
	newBot.RegisterCommandCallback("delete-date_update", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackUpdatePublicationDeleteDate()))
	newBot.RegisterCommandCallback("check_publication", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackCheckPublication()))
	newBot.RegisterCommandCallback("publication_cancel", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackGetListForCancelPublication()))
	newBot.RegisterCommandCallback("publication_delete", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackDeletePublication()))
	newBot.RegisterCommandCallback("cancel_update", middleware.AdminMiddleware(b.userService, b.callbackPublication.CallbackCancelUpdate()))

	b.log.Info("Initialize bot took [%f] seconds", time.Since(startBot).Seconds())
	if err := newBot.Run(ctx); err != nil {
		b.log.Fatal("failed to run Telegram Bot: %v", err)
	}
}
