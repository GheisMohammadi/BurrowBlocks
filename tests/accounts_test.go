package tests

import (
	"testing"

	config "github.com/BurrowBlocks/config"
	db "github.com/BurrowBlocks/database"
	"github.com/stretchr/testify/require"
)

func TestAccountsInDataBase(t *testing.T) {

	gConfig, _ := config.LoadConfigFile(true)
	dbe := db.Postgre{Config: gConfig}
	connErr := dbe.Connect()
	require.NoError(t, connErr)

	defer dbe.Disconnect()

	/*
		acc := db.Account{Address: "Addr123456", PublicKey: "ABC", Balance: 1234.56, Permission: "Perm456", Sequence: 2, Code: "CodeF1F2"}
		insertErr := dbe.InsertAccount(&acc)
		require.NoError(t, insertErr)

		sAcc, GAccErr := dbe.GetAccount(7)
		require.NoError(t, GAccErr)
		require.Equal(t, sAcc.ID, 7)

		sAcc, GNoAccErr := dbe.GetAccount(10000000)
		require.Error(t, GNoAccErr)
	*/
}
