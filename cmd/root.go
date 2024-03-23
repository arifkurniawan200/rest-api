package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"template/cmd/migration"
	"template/config"
	"template/db"
	"template/internal/app"
	"template/internal/repository"
	"template/internal/usecase"
)

func Start() {
	cfg := config.ReadConfig()
	// root command
	root := &cobra.Command{}

	// command allowed
	cmds := []*cobra.Command{
		{
			Use:   "db:migrate",
			Short: "database migration",
			Run: func(cmd *cobra.Command, args []string) {
				migration.RunMigration(cfg)
			},
		},
		{
			Use:   "api",
			Short: "run api server",
			Run: func(cmd *cobra.Command, args []string) {
				dbs, err := db.NewDatabase(cfg.DB)
				if err != nil {
					log.Fatal(err)
				}

				userRepo := repository.NewUserRepository(dbs)
				transactionRepo := repository.NewTransactionRepository(dbs)
				userUsecase := usecase.NewUserUsecase(userRepo, transactionRepo)
				transactionUcase := usecase.NewTransactionsUsecase(transactionRepo, userRepo)
				app.Run(userUsecase, transactionUcase)
			},
		},
	}
	root.AddCommand(cmds...)
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
