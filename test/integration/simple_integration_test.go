package integration

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntegrationSimple(t *testing.T) {
	if testing.Short() {
		t.Skip("This is an integration test")
	}

	t.Run("Should list containers", func(t *testing.T) {
		cli, err := client.NewEnvClient()
		require.NoError(t, err)

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
		require.NoError(t, err)

		for _, container := range containers {
			t.Logf("Container.ID=%s", container.ID)
		}
		assert.Empty(t, containers)
	})

}
