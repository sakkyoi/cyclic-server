package secretary

import (
	"context"
	"cyclic/ent"
	"cyclic/pkg/colonel"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var Minute *ent.Client

func Init() error {
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
		return fmt.Errorf("unsupported database driver: %s", driver)
	}

	client, err := ent.Open(driver, dsn)
	if err != nil {
		return err
	}

	Minute = client

	// do migration
	if err := Migrate(); err != nil {
		return err
	}

	return nil
}

func Migrate() error {
	if err := Minute.Schema.Create(context.Background()); err != nil {
		return err
	}

	return nil
}
