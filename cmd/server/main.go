package main

import (
	"context"

	"os"

	"github.com/heroticket/internal/cache/redis"
	"github.com/heroticket/internal/db"
	"github.com/heroticket/internal/db/mongo"
	"github.com/heroticket/internal/did"
	didrepo "github.com/heroticket/internal/did/repository/mongo"
	"github.com/heroticket/internal/http"
	"github.com/heroticket/internal/http/rest"
	"github.com/heroticket/internal/jwt"
	"github.com/heroticket/internal/notice"
	noticerepo "github.com/heroticket/internal/notice/repository/mongo"
	"github.com/heroticket/internal/user"
	userrepo "github.com/heroticket/internal/user/repository/mongo"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
)

//	@title			Hero Ticket API
//	@version		1.0
//	@description	This is Hero Ticket API server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	logger, _ := zap.NewProduction(zap.Fields(zap.String("service", "hero-ticket")))
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	zap.L().Info("starting server")

	defer func() {
		if r := recover(); r != nil {
			zap.L().Info("recovered from panic", zap.Any("r", r))
		}
	}()

	client, err := mongo.New(context.Background(), os.Getenv("MONGO_URL"))
	if err != nil {
		panic(err)
	}

	zap.L().Info("connected to mongo")

	cache, err := redis.New(context.Background(), os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	zap.L().Info("connected to redis")

	didRepo := didrepo.NewRepository(client, "hero-ticket", "dids")
	noticeRepo := noticerepo.New(client, "hero-ticket", "notices")
	userRepo := userrepo.NewMongoRepository(client, "hero-ticket", "users")

	// didSvc := did.New(config.DIDServiceConfig{
	// 	RPCUrl:    os.Getenv("RPC_URL"),
	// 	IssuerUrl: "https://issuer.heroticket.xyz",
	// 	Username:  "user-issuer",
	// 	Password:  "password-issuer",

	// 	Repo:         didRepo,
	// 	RequestCache: redis.NewCache(cache),
	// 	Client:       ohttp.DefaultClient,
	// })

	jwtSvc := jwt.New("secret1", "secret2")
	noticeSvc := notice.New(noticeRepo)
	userSvc := user.New(userRepo)
	tx := mongo.NewTx(client)

	_ = newServer(
		initUserController(didSvc, jwtSvc, userSvc, tx),
		initNoticeController(jwtSvc, noticeSvc, userSvc),
		initTicketController(),
	)

	// go func() {
	// 	if err := server.Run(); err != nil && err != http.ErrServerClosed {
	// 		panic(err)
	// 	}
	// }()

	// stop := shutdown.GracefulShutdown(func() {
	// 	if err := server.Shutdown(context.Background()); err != nil {
	// 		panic(err)
	// 	}

	// 	zap.L().Info("server stopped")

	// 	if err := client.Disconnect(context.Background()); err != nil {
	// 		panic(err)
	// 	}

	// 	zap.L().Info("disconnected from mongo")
	// }, syscall.SIGINT, syscall.SIGTERM)

	// <-stop
}

func initUserController(didSvc did.Service, jwtSvc jwt.Service, userSvc user.Service, tx db.Tx) *rest.UserCtrl {
	return rest.NewUserCtrl(didSvc, jwtSvc, userSvc, tx)
}

func initNoticeController(jwtSvc jwt.Service, noticeSvc notice.Service, userSvc user.Service) *rest.NoticeCtrl {
	return rest.NewNoticeCtrl(jwtSvc, noticeSvc, userSvc)
}

func initTicketController() *rest.TicketCtrl {
	return rest.NewTicketCtrl()
}

func newServer(ctrls ...http.Controller) *http.Server {
	return http.NewServer(http.DefaultConfig(), ctrls...)
}
