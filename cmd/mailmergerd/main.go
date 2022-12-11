package main

import (
	"fmt"
	"os"

	"github.com/fahmifan/mailmerger-server/config"
	"github.com/fahmifan/mailmerger-server/db/gorms"
	"github.com/fahmifan/mailmerger-server/pkg/localfs"
	"github.com/fahmifan/mailmerger-server/pkg/smtp"
	"github.com/fahmifan/mailmerger-server/server"
	"github.com/fahmifan/mailmerger-server/service"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
		fmt.Println(err)
	}
}

func run() error {
	cmd := cobra.Command{}
	cmd.AddCommand(runServer())
	return cmd.Execute()
}

func runServer() *cobra.Command {
	cmd := cobra.Command{
		Use: "server",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
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
		srv := server.NewServer(svc)

		srv.Run()
		return
	}

	return &cmd
}
