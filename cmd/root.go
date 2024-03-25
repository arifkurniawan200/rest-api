package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"template/cmd/migration"
	"template/config"
	"template/db"
	"template/internal/app"
	"template/internal/repository/items"
	"template/internal/repository/transaction"
	"template/internal/repository/user"
	"template/internal/usecase/item"
	transaction_ucase "template/internal/usecase/transaction"
	user_ucase "template/internal/usecase/user"
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

				userRepo := user.NewUserRepository(dbs)
				transactionRepo := transaction.NewTransactionRepository(dbs)
				userUsecase := user_ucase.NewUserUsecase(userRepo, transactionRepo)
				historyRepo := items.NewHistoryRepository(dbs)
				itemRepository := items.NewItemRepository(dbs)
				itemUsecase := item.NewItemUsecase(transactionRepo, userRepo, historyRepo, itemRepository)
				transactionUcase := transaction_ucase.NewTransactionsUsecase(transactionRepo, userRepo, itemRepository)
				app.Run(userUsecase, transactionUcase, itemUsecase)
			},
		},
	}
	root.AddCommand(cmds...)
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
