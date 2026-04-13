package routes

import (
	"database/sql"
	"mini-exchange/app/libraries"
	"mini-exchange/config"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type routes struct {
	Api fiber.Router
	MysqlDb *sql.DB
	TransactionManager libraries.TransactionManager
}

type router struct {
	app  *fiber.App
	MysqlDb *sql.DB
	goquLibrary libraries.GoquLibrary
	transactionManager libraries.TransactionManager
	cfg                 *config.Config
	startOnce sync.Once
}

func NewRouter(
	app *fiber.App,
	db *sql.DB,
	cfg *config.Config,
) *router {

	return &router{
		app:app,
		goquLibrary: libraries.NewGoquLibrary(db),
		transactionManager: libraries.NewTransactionManager(db),
		MysqlDb: db,
		cfg: cfg,
	}
}

// Wrapper function to initialize all routes
func (r *router) Init() routes {
	return routes{
		Api: r.api(),
		MysqlDb: r.MysqlDb,
		TransactionManager: r.transactionManager,
	}
}
