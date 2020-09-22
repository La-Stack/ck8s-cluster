package azure

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/elastisys/ck8s/api"
)

func TestAddMachine(t *testing.T) {
	cluster := Default(-1, "testName")

	want := &api.Machine{
		NodeType: api.Master,
		Size:     "Standard_D2_v3",
		Image:    api.NewImage("test", "v1.2.3"),
	}

	for _, clusterType := range []api.ClusterType{
		api.ServiceCluster,
		api.WorkloadCluster,
	} {
		cluster.config.ClusterType = clusterType

		name, err := cluster.AddMachine("", want)
		if err != nil {
			t.Fatalf(
				"error while adding %s machine: %s",
				clusterType.String(), err,
			)
		}

		machines := cluster.Machines()

		got, ok := machines[name]
		if !ok {
			t.Errorf(
				"machine missing: %s", name,
			)
		}

		if diff := cmp.Diff(
			want,
			got,
			cmpopts.IgnoreFields(api.Image{}, "KubeletVersion"),
		); diff != "" {
			t.Errorf("machine mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestAddMachineWithName(t *testing.T) {
	name := "foo"

	cluster := Default(-1, "testName")

	want := &api.Machine{
		NodeType: api.Worker,
		Size:     "Standard_D2_v3",
		Image:    api.NewImage("test", "v1.2.3"),
	}

	for _, clusterType := range []api.ClusterType{
		api.ServiceCluster,
		api.WorkloadCluster,
	} {
		cluster.config.ClusterType = clusterType

		_, err := cluster.AddMachine(name, want)
		if err != nil {
			t.Fatalf(
				"error while adding %s machine: %s",
				clusterType.String(), err,
			)
		}

		machines := cluster.Machines()

		got, ok := machines[name]
		if !ok {
			t.Errorf(
				"machine missing: %s", name,
			)
		}

		if diff := cmp.Diff(
			want,
			got,
			cmpopts.IgnoreFields(api.Image{}, "KubeletVersion"),
		); diff != "" {
			t.Errorf("machine mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestRemoveMachine(t *testing.T) {
	testName := "bar"

	got, want := Default(-1, "testName"), Default(-1, "testName")

	got.tfvars = AzureTFVars{
		MachinesSC: map[string]*api.Machine{
			testName: {
				NodeType: api.Master,
				Size:     "Standard_D2_v3",
				Image:    api.NewImage("test", "v1.2.3"),
			},
		},
		MachinesWC: map[string]*api.Machine{
			testName: {
				NodeType: api.Worker,
				Size:     "Standard_D4_v3",
				Image:    api.NewImage("test", "v1.2.3"),
			},
		},
	}

	want.tfvars = AzureTFVars{
		MachinesSC: map[string]*api.Machine{},
		MachinesWC: map[string]*api.Machine{},
	}

	for _, clusterType := range []api.ClusterType{
		api.ServiceCluster,
		api.WorkloadCluster,
	} {
		got.config.ClusterType = clusterType

		if err := got.RemoveMachine(testName); err != nil {
			t.Fatalf(
				"error while removing %s machine: %s",
				clusterType.String(), err,
			)
		}
	}

	if diff := cmp.Diff(want.tfvars, got.tfvars); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
