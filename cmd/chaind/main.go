package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/kingblockio/chain/app"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "kingchaind",
		Short:             "kingchain Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, server.DefaultAppInit,
		server.ConstructAppCreator(newApp, "kingchain"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "kingchain"))

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.kingchaind")
	executor := cli.PrepareBaseCmd(rootCmd, "BC", rootDir)

	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	return app.NewKingBlockApp(logger, db, baseapp.SetPruning(viper.GetString("pruning")))
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	bapp := app.NewKingBlockApp(logger, db)
	return bapp.ExportAppStateAndValidators()
}
