package main

import (
	"embed"

	"github.com/fahmifan/mailmerger-server/config"
	"github.com/fahmifan/mailmerger-server/db/gorms"
	"github.com/fahmifan/mailmerger-server/pkg/localfs"
	"github.com/fahmifan/mailmerger-server/pkg/smtp"
	"github.com/fahmifan/mailmerger-server/service"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"gorm.io/gorm"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	cfg := config.GetConfig()
	sqlite := gorms.SQLite{
		DSN:              cfg.SQLite.GetDBName(),
		EnableWAL:        cfg.SQLite.EnableWAL,
		EnableForeignKey: cfg.SQLite.EnableForeignKey,
	}
	postgres := gorms.Postgres{DSN: cfg.Postgres.DSN()}

	var db *gorm.DB
	switch cfg.DbDriver {
	case config.DbDriverPostgres:
		db = postgres.MustOpen()
	case config.DbDriverSQLite:
		db = sqlite.MustOpen()
	}
	defer gorms.Close(db)

	localFS := localfs.Storage{
		RootDir: "private",
	}
	smptTransporter, err := smtp.NewSmtpClient(&smtp.Config{
		Host: "0.0.0.0",
		Port: 1025,
	})
	if err != nil {
		return
	}

	blastEmailCfg := service.BlastEmailConfig{
		Sender:      cfg.Mailer.SenderAddress,
		Concurrency: cfg.Mailer.Concurrency,
		Transporter: smptTransporter,
	}
	svc := service.NewService(db, &localFS, &blastEmailCfg)

	// Create an instance of the app structure
	app := NewApp(svc)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "mailmerger-app",
		Width:  options.Default.Width,
		Height: options.Default.Height,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
