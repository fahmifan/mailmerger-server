package main

import (
	"fmt"
	"os"

	"github.com/fahmifan/mailmerger-server/pkg/localfs"
	"github.com/fahmifan/mailmerger-server/pkg/smtp"
	"github.com/fahmifan/mailmerger-server/server"
	"github.com/fahmifan/mailmerger-server/service"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
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
		dsn := "host=localhost user=root password=root dbname=mailmerger port=5432 sslmode=disable"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
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
			Sender:      "admin@mailmerger.com",
			Concurrency: 4,
			Transporter: smptTransporter,
		}
		svc := service.NewService(db, &localFS, &blastEmailCfg)
		srv := server.NewServer(svc)

		srv.Run()
		return
	}
	return &cmd
}
