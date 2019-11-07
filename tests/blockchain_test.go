package tests

import (
	"testing"

	bc "github.com/BurrowBlocks/blockchain"
	config "github.com/BurrowBlocks/config"
	"github.com/stretchr/testify/require"
)

func TestBlockChain(t *testing.T) {

	gConfig, _ := config.LoadConfigFile(true)
	bc := bc.Gallactic{Config: gConfig}

	clientErr := bc.CreateClient()
	require.NoError(t, clientErr)

	updateErr := bc.Update()
	require.NoError(t, updateErr)

	_, getBlockErr := bc.GetBlock(uint64(14141))
	require.NoError(t, getBlockErr)

	_, getTXsErr := bc.GetTXs(uint64(14141))
	require.NoError(t, getTXsErr)

	/*
		ret1, err := bc.GetAccounts()
		require.NoError(t, err)

		numAccounts := len(ret1.Accounts)
		require.NotEqual(t, numAccounts, 0)

		_, err2 := client.GetGenesis(context.Background(), &pb.Empty{})
		require.NoError(t, err2)

		_, err3 := client.GetBlocks(context.Background(), &pb.BlocksRequest{})
		require.NoError(t, err3)
	*/
}
