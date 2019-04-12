package integration

import (
	"github.com/danielpacak/docker-sdk-experiments/test/common/docker"
	"github.com/danielpacak/docker-sdk-experiments/test/common/io"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strconv"
	"testing"
	"time"
)

func TestK8SIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("This is an integration test")
	}

	const network = "k8s-itest"

	// setup code
	dc, err := docker.NewDockerController()
	require.NoError(t, err)

	_, err = dc.Network().Create(network)
	require.NoError(t, err)

	containerID, err := dc.Container().Builder().
		WithImage("rancher/k3s:v0.3.0").
		WithCmd([]string{"server", "--disable-agent"}).
		WithName("k3s").
		WithEnv("K3S_CLUSTER_SECRET", "somethingtotallyrandom").
		WithEnv("K3S_KUBECONFIG_OUTPUT", "/output/kubeconfig.yaml").
		WithEnv("K3S_KUBECONFIG_MODE", "666").
		WithNetwork(network).
		WithExposedPorts(map[nat.Port]struct{}{"6443/tcp": {}}).
		WithPortBindings(map[nat.Port][]nat.PortBinding{
			"6443/tcp": {
				nat.PortBinding{HostPort: strconv.Itoa(6443)},
			},
		}).
		WithMounts([]mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/tmp",
				Target: "/output",
			},
		}).
		Create()
	require.NoError(t, err)

	err = dc.Container().Start(containerID)
	require.NoError(t, err)

	t.Run("Should list namespaces", func(t *testing.T) {
		err = io.WaitExists(t, "/tmp/kubeconfig.yaml", 30*time.Second)
		require.NoError(t, err)

		config, err := clientcmd.BuildConfigFromFlags("", "/tmp/kubeconfig.yaml")
		require.NoError(t, err)

		clientset, err := kubernetes.NewForConfig(config)
		require.NoError(t, err)

		list, err := clientset.CoreV1().Namespaces().List(v1.ListOptions{})
		require.NoError(t, err)

		var namespaces []string
		for _, n := range list.Items {
			t.Logf("namespace: %s, status: %s", n.Name, n.Status)
			namespaces = append(namespaces, n.Name)
		}
		assert.Contains(t, namespaces, "default")
		assert.Contains(t, namespaces, "kube-system")
		assert.Contains(t, namespaces, "kube-public")
	})

	// tear-down code
	err = dc.Container().Stop(containerID)
	require.NoError(t, err)

	err = dc.Network().Remove(network)
	require.NoError(t, err)
}
