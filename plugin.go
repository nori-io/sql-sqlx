
package main

import (
"context"

"github.com/jmoiron/sqlx"
"github.com/nori-io/common/v3/config"
"github.com/nori-io/common/v3/logger"
"github.com/nori-io/common/v3/meta"
"github.com/nori-io/common/v3/plugin"
i "github.com/nori-io/interfaces/public/sql/sqlx"

_ "github.com/jinzhu/gorm/dialects/mssql"
_ "github.com/jinzhu/gorm/dialects/mysql"
_ "github.com/jinzhu/gorm/dialects/postgres"
_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type service struct {
	db     *sqlx.DB
	config *pluginConfig
	logger logger.FieldLogger
}

type pluginConfig struct {
	dsn     string
	dialect string
	logMode bool
}

var (
	Plugin plugin.Plugin = &service{}
)

func (p *service) Init(ctx context.Context, config config.Config, log logger.FieldLogger) error {
	p.logger = log
	p.config.logMode = config.Bool("logMode", "log mode: true or false")()
	p.config.dsn = config.String("dsn", "database connection string")()
	p.config.dialect = config.String("dialect", "sql dialect: mssql, mysql, postgres, sqlite")()
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
	p.db, err = sqlx.Open(p.config.dialect, p.config.dsn)
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
