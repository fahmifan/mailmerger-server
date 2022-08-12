package main

import (
	"fmt"
	"os"

	"github.com/fahmifan/mailmerger-server/localfs"
	"github.com/fahmifan/mailmerger-server/server"
	"github.com/fahmifan/mailmerger-server/service"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
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
		db, err := bbolt.Open("./mailmerger.bbolt.db", 0666, nil)
		if err != nil {
			return
		}
		localFS := localfs.Storage{
			RootDir: "private",
		}
		svc := service.NewService(db, &localFS)
		srv := server.NewServer(svc)

		srv.Run()
		return
	}
	return &cmd
}
