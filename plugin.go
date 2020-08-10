package main

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/nori-io/common/v3/config"
	"github.com/nori-io/common/v3/logger"
	"github.com/nori-io/common/v3/meta"
	"github.com/nori-io/common/v3/plugin"
	i "github.com/nori-io/interfaces/public/sql/sqlx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type service struct {
	db     *sqlx.DB
	config *pluginConfig
	logger logger.FieldLogger
}

type pluginConfig struct {
	dsn     string
	driver  string
	logMode bool
}

var (
	Plugin  plugin.Plugin = &service{}
	drivers               = [3]string{"mysql", "postgres", "sqlite3"}
)

func (p *service) Init(ctx context.Context, config config.Config, log logger.FieldLogger) error {
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

func (p *service) Instance() interface{} {
	return p.db
}

func (p *service) Meta() meta.Meta {
	return &meta.Data{
		ID: meta.ID{
			ID:      "sql/sqlx",
			Version: "1.2.0",
		},
		Author: meta.Author{
			Name: "Nori.io",
			URI:  "https://nori.io/",
		},
		Core: meta.Core{
			VersionConstraint: "=0.2.0",
		},
		Dependencies: []meta.Dependency{},
		Description: meta.Description{
			Name:        "Nori: Sqlx",
			Description: "This plugin implements instance of Sqlx",
		},
		Interface: i.SqlxInterface,
		License: []meta.License{
			{
				Title: "GPLv3",
				Type:  "GPLv3",
				URI:   "https://www.gnu.org/licenses/"},
		},
		Tags: []string{"sqlx", "sql", "database", "db"},
	}
}

func (p *service) Start(ctx context.Context, registry plugin.Registry) error {
	var err error
	p.db, err = sqlx.Open(p.config.driver, p.config.dsn)
	if err != nil {
		p.logger.Error(err.Error())
	}
	return err
}

func (p *service) Stop(ctx context.Context, registry plugin.Registry) error {
	err := p.db.Close()
	if err != nil {
		p.logger.Error(err.Error())
	}
	return err
}
