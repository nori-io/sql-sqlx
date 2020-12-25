package plugin

import (
	"context"
	"errors"

	"github.com/nori-io/common/v4/pkg/domain/registry"

	"github.com/jmoiron/sqlx"
	"github.com/nori-io/common/v4/pkg/domain/config"
	em "github.com/nori-io/common/v4/pkg/domain/enum/meta"
	"github.com/nori-io/common/v4/pkg/domain/logger"
	"github.com/nori-io/common/v4/pkg/domain/meta"
	p "github.com/nori-io/common/v4/pkg/domain/plugin"
	m "github.com/nori-io/common/v4/pkg/meta"

	i "github.com/nori-io/interfaces/database/sql/sqlx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Plugin  p.Plugin = plugin{}
	drivers          = [3]string{"mysql", "postgres", "sqlite3"}
)

type plugin struct {
	db     *sqlx.DB
	config conf
	logger logger.FieldLogger
}

type conf struct {
	dsn     string
	driver  string
	logMode bool
}

func (p plugin) Init(ctx context.Context, config config.Config, log logger.FieldLogger) error {
	p.logger = log
	p.config.dsn = config.String("sql.sqlx.dsn", "database connection string")()
	p.config.driver = config.String("sql.sqlx.driver", "sql driver: postgres")()

	var isValidDriver bool

	for _, v := range drivers {
		if v == p.config.driver {
			isValidDriver = true
		}
	}

	if !isValidDriver {
		return errors.New("Driver is wrong. You should use on of sql drivers: mysql, postgres, sqlite3")
	}

	return nil
}

func (p plugin) Instance() interface{} {
	return p.db
}

func (p plugin) Meta() meta.Meta {
	return m.Meta{
		ID: m.ID{
			ID:      "sql/sqlx",
			Version: "1.2.0",
		},
		Author: m.Author{
			Name: "Nori.io",
			URL:  "https://nori.io/",
		},
		Dependencies: []meta.Dependency{},
		Description: m.Description{
			Title:       "",
			Description: "This plugin implements instance of Sqlx",
		},
		Interface: i.SqlxInterface,
		License:   []meta.License{},
		Links:     []meta.Link{},
		Repository: m.Repository{
			Type: em.Git,
			URL:  "https://github.com/nori-io/sql-sqlx",
		},
		Tags: []string{"sqlx", "sql", "database", "db"},
	}
}

func (p plugin) Start(ctx context.Context, registry registry.Registry) error {
	var err error
	p.db, err = sqlx.Open(p.config.driver, p.config.dsn)
	if err != nil {
		p.logger.Error(err.Error())
	}
	return err
}

func (p plugin) Stop(ctx context.Context, registry registry.Registry) error {
	err := p.db.Close()
	if err != nil {
		p.logger.Error(err.Error())
	}
	return err
}
