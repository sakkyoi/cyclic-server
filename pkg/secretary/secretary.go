package secretary

import (
	"context"
	"cyclic/ent"
	"cyclic/pkg/colonel"
	"cyclic/pkg/scribe"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

var Minute *ent.Client

func Init() {
	driver := colonel.Writ.Database.Driver

	var dsn string

	// Set the DSN based on the driver
	// TODO: Add support for other drivers
	switch driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True",
			colonel.Writ.Database.User,
			colonel.Writ.Database.Password,
			colonel.Writ.Database.Host,
			colonel.Writ.Database.Port,
			colonel.Writ.Database.Name)
	default:
		scribe.Scribe.Fatal("unsupported database driver", zap.String("driver", driver))
	}

	client, err := ent.Open(driver, dsn)
	if err != nil {
		scribe.Scribe.Fatal("failed to open database", zap.Error(err))
	}

	Minute = client

	// do migration
	if err := Migrate(); err != nil {
		scribe.Scribe.Fatal("failed to migrate database", zap.Error(err))
	}
}

func Migrate() error {
	if err := Minute.Schema.Create(context.Background()); err != nil {
		return err
	}

	return nil
}
