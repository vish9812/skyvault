package main

import (
	"skyvault/api"
	"skyvault/app"
	"skyvault/common"
	"skyvault/infra/store"
	"skyvault/infra/store/db_store"
)

func main() {
	common.LoadConfig("./", "dev", "env")

	newDBStore := db_store.NewDBStore(common.Configs.DB_CONN_STR)
	newStore := store.NewStore(newDBStore)
	newApp := app.NewApp(newStore)
	newAPI := api.NewAPI(newApp)

	newAPI.Run()
}
