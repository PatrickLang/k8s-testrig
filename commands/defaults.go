package commands

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type properties struct {
	OrchestratorProfile *orchestratorProfile `json:"orchestratorProfile"`
	MasterProfile       *masterProfile       `json:"masterProfile"`
	AgentPoolProfiles   []agentPoolProfile   `json:"agentPoolProfiles"`
	LinuxProfile        *linuxProfile        `json:"linuxProfile,omitempty"`
	WindowsProfile      *windowsProfile      `json:"windowsProfile",omitempty`
}

type orchestratorProfile struct {
	OrchestratorType    string            `json:"orchestratorType"`
	OrchestratorRelease string            `json:"orchestratorRelease"`
	KubernetesConfig    *kubernetesConfig `json:"kubernetesConfig"`
}

type kubernetesConfig struct {
	UseManagedIdentity bool   `json:"useManagedIdentity,omitempty"`
	NetworkPlugin      string `json:"networkPlugin,omitempty"`
	NetworkPolicy      string `json:"networkPolicy,omitempty"`
	ContainerRuntime   string `json:"containerRuntime,omitempty"`
}

type masterProfile struct {
	Count          int    `json:"count"`
	VMSize         string `json:"vmSize"`
	OSDiskSizeGB   int    `json:"osDiskSizeGB,omitempty"`
	StorageProfile string `json:"storageProfile,omitempty"`
	DNSPrefix      string `json:"dnsPrefix"`
	Distro         string `json:"distro,omitempty"`
}

type agentPoolProfile struct {
	Name                         string `json:"name"`
	Count                        int    `json:"count"`
	VMSize                       string `json:"vmSize"`
	OSDiskSizeGB                 int    `json:"osDiskSizeGB,omitempty"`
	StorageProfile               string `json:"storageProfile,omitempty"`
	AcceleratedNetworkingEnabled *bool  `json:"acceleratedNetworkingEnabled,omitempty"`
	OSType                       string `json:"osType"`
	AvailabilityProfile          string `json:"availabilityProfile,omitempty"`
	Distro                       string `json:"distro,omitempty"`
}

type linuxProfile struct {
	AdminUsername string    `json:"adminUsername"`
	SSH           sshConfig `json:"ssh"`
}

type windowsProfile struct {
	AdminUsername string `json:"adminUsername"`
	AdminPassword string `json:"adminPassword"`
}

type sshConfig struct {
	PublicKeys []sshKey `json:"publicKeys"`
}

type sshKey struct {
	KeyData string `json:"keyData"`
}

type apiModel struct {
	APIVersion string      `json:"apiVersion"`
	Properties *properties `json:"properties"`
}

// Defaults creates the `generate-defaults` subcommand.
// Use this to generate a default API model.
func Defaults() *cobra.Command {
	return &cobra.Command{
		Use: "generate-defaults",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}

func defaultModel() *apiModel {
	return &apiModel{
		APIVersion: "vlabs",
		Properties: &properties{
			OrchestratorProfile: &orchestratorProfile{
				OrchestratorType:    "Kubernetes",
				OrchestratorRelease: "1.10",
				KubernetesConfig: &kubernetesConfig{
					UseManagedIdentity: true,
					NetworkPlugin:      "azure",
					NetworkPolicy:      "azure",
				},
			},
			MasterProfile: &masterProfile{
				Count:          3,
				VMSize:         "Standard_DS2_v2",
				StorageProfile: "ManagedDisks",
				OSDiskSizeGB:   200,
			},
			AgentPoolProfiles: []agentPoolProfile{
				agentPoolProfile{
					Name:                         "linuxpool1",
					Count:                        3,
					VMSize:                       "Standard_DS2_v2",
					StorageProfile:               "ManagedDisks",
					OSDiskSizeGB:                 200,
					AvailabilityProfile:          "VirtualMachineScaleSets",
					AcceleratedNetworkingEnabled: boolPtr(true),
					OSType:                       "Linux",
				},
				agentPoolProfile{
					Name:                "windowspool1",
					Count:               0,
					VMSize:              "Standard_DS2_v3",
					StorageProfile:      "ManagedDisks",
					OSDiskSizeGB:        200,
					AvailabilityProfile: "VirtualMachineScaleSets",
					OSType:              "Windows",
				},
			},
			LinuxProfile: &linuxProfile{
				AdminUsername: "azureuser",
			},
			WindowsProfile: &windowsProfile{
				AdminUsername: "azureuser",
			},
		},
	}
}

func overrideModelDefaults(m *apiModel, cfg *UserConfig) error {
	if cfg == nil {
		return nil
	}

	if cfg.Profile.Leader.Linux.Count != nil {
		m.Properties.MasterProfile.Count = *cfg.Profile.Leader.Linux.Count
	}
	if cfg.Profile.Leader.Linux.SKU != "" {
		m.Properties.MasterProfile.VMSize = cfg.Profile.Leader.Linux.SKU
	}

	if cfg.Profile.Agent.Linux.Count != nil {
		m.Properties.AgentPoolProfiles[0].Count = *cfg.Profile.Agent.Linux.Count
	}
	if cfg.Profile.Agent.Linux.SKU != "" {
		m.Properties.AgentPoolProfiles[0].VMSize = cfg.Profile.Agent.Linux.SKU
	}

	if cfg.Profile.Agent.Windows.Count != nil {
		m.Properties.AgentPoolProfiles[1].Count = *cfg.Profile.Agent.Windows.Count
	}
	if cfg.Profile.Agent.Windows.SKU != "" {
		m.Properties.AgentPoolProfiles[1].VMSize = cfg.Profile.Agent.Windows.SKU
	}

	if cfg.Profile.Auth.Linux.User != "" {
		m.Properties.LinuxProfile.AdminUsername = cfg.Profile.Auth.Linux.User
	}
	if cfg.Profile.Auth.Linux.PublicKeyFile != "" {
		keyData, err := ioutil.ReadFile(cfg.Profile.Auth.Linux.PublicKeyFile)
		if err != nil {
			return errors.Wrap(err, "error reading user supplied linux public ssh key file")
		}
		m.Properties.LinuxProfile.SSH.PublicKeys = append(m.Properties.LinuxProfile.SSH.PublicKeys, sshKey{KeyData: string(keyData)})
	}

	if cfg.Profile.Auth.Windows.User != "" {
		m.Properties.WindowsProfile.AdminUsername = cfg.Profile.Auth.Windows.User
	}
	if cfg.Profile.Auth.Windows.PasswordFile != "" {
		pData, err := ioutil.ReadFile(cfg.Profile.Auth.Windows.PasswordFile)
		if err != nil {
			return errors.Wrap(err, "error reading user supplied windows admin password file")
		}
		m.Properties.WindowsProfile.AdminPassword = string(pData)
	}

	if cfg.Profile.Kubernetes.Version != "" {
		m.Properties.OrchestratorProfile.OrchestratorRelease = cfg.Profile.Kubernetes.Version
	}
	if cfg.Profile.Kubernetes.NetworkPlugin != "" {
		m.Properties.OrchestratorProfile.KubernetesConfig.NetworkPlugin = cfg.Profile.Kubernetes.NetworkPlugin
	}
	if cfg.Profile.Kubernetes.NetworkPolicy != "" {
		m.Properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy = cfg.Profile.Kubernetes.NetworkPolicy
	}

	return nil
}
