package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster describes a cluster with kubernetes and addon
//
// Cluster are non-namespaced; the id of the cluster
// according to etcd is in ObjectMeta.Name.
type Cluster struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   ClusterSpec   `json:"spec"`
	Status ClusterStatus `json:"status"`
}

type ClusterSpec struct {
	DisplayName      string                     `json:"displayName"`
	Provider         CloudProvider              `json:"provider"`
	ProviderConfig   ClusterCloudProviderConfig `json:"providerConfig"`
	IsControlCluster bool                       `json:"isControlCluster"`
	Network          ClusterNetwork             `json:"network"`
	IsHighAvailable  bool                       `json:"isHighAvailable"`
	MastersVIP       string                     `json:"mastersVIP"`
	Auth             ClusterAuth                `json:"auth"`
	Versions         *ClusterVersions           `json:"versions,omitempty"`
	Masters          []string                   `json:"masters"`
	Nodes            []string                   `json:"nodes"`
	// deploy
	DeployToolsExternalVars map[string]string `json:"deployToolsExternalVars"`
	// adapt expired
	ClusterToken string       `json:"clusterToken"`
	Ratio        ClusterRatio `json:"ratio"`
}

type ClusterStatus struct {
	Phase         ClusterPhase                       `json:"phase"`
	Conditions    []ClusterCondition                 `json:"conditions"`
	Masters       []MachineThumbnail                 `json:"masters"`
	Nodes         []MachineThumbnail                 `json:"nodes"`
	Capacity      map[ResourceName]resource.Quantity `json:"capacity"`
	OperationLogs []OperationLog                     `json:"operationLogs,omitempty"`
	AutoScaling   ClusterAutoScalingStatus           `json:"autoScaling,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList is a collection of clusters
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of Clusters
	Items []Cluster `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Machine describes real machines
//
// Machine are non-namespaced; the id of the machine
// according to etcd is in ObjectMeta.Name.
type Machine struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   MachineSpec   `json:"spec"`
	Status MachineStatus `json:"status"`
}

type MachineSpec struct {
	Provider       CloudProvider              `json:"provider"`
	ProviderConfig MachineCloudProviderConfig `json:"providerConfig,omitempty"`
	Address        []NodeAddress              `json:"address"`
	SshPort        string                     `json:"sshPort"`
	Auth           MachineAuth                `json:"auth"`
	Versions       MachineVersions            `json:"versions,omitempty"`
	Cluster        string                     `json:"cluster"`
	IsMaster       bool                       `json:"isMaster"`
	Tags           map[string]string          `json:"tags"`
}

type MachineStatus struct {
	Phase MachinePhase `json:"phase"`
	// env
	Environment MachineEnvironment `json:"environment"`
	// node about
	NodeRefer  string                             `json:"nodeRefer"`
	Capacity   map[ResourceName]resource.Quantity `json:"capacity"`
	NodeStatus MachineNodeStatus                  `json:"nodeStatus"`
	// other
	OperationLogs []OperationLog `json:"operationLogs,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineList is a collection of machine.
type MachineList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of Machines
	Items []Machine `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// cloud provider

type MachineCloudProviderConfig struct {
	// auto scaling

	// auto scaling group name which this machine belongs to
	// empty means not belongs to any auto scaling group
	AutoScalingGroup string `json:"autoScalingGroup,omitempty"`

	Azure *AzureMachineCloudProviderConfig `json:"azure,omitempty"`
}

type ClusterCloudProviderConfig struct {
	// auto scaling

	// cluster level auto scaling setting
	// maybe nil in old or control cluster, need to be inited with a default setting by controller
	AutoScalingSetting *ClusterAutoScalingSetting `json:"autoScalingSetting,omitempty"`

	Azure *AzureClusterCloudProviderConfig `json:"azure,omitempty"`
}

// azure

type AzureObjectMeta struct {
	// ID - Resource ID.
	ID string `json:"id,omitempty"`
	// Name - Resource name.
	Name string `json:"name,omitempty"`
	// Location - Resource location.
	Location string `json:"location,omitempty"`
	// ResourceGroupName - Resource group name
	GroupName string `json:"groupName,omitempty"`
}

type AzureClusterCloudProviderConfig struct {
	// Location - cluster azure resource location.
	Location string `json:"location,omitempty"`
	// VirtualNetwork - cluster azure virtual network
	VirtualNetwork AzureVirtualNetwork `json:"virtualNetwork"`
}

type AzureMachineCloudProviderConfig struct {
	AzureObjectMeta
	VirtualNetwork    AzureVirtualNetwork     `json:"virtualNetwork"`
	LoginUser         string                  `json:"loginUser"`
	LoginPassword     string                  `json:"loginPassword"`
	VMSize            AzureVMSize             `json:"vmSize"`
	ImageReference    AzureImageReference     `json:"imageReference"`
	OSDisk            AzureDisk               `json:"osDisk"`
	DataDisks         []AzureDisk             `json:"dataDisks"`
	NetworkInterfaces []AzureNetworkInterface `json:"networkInterfaces"`
}

type AzureVirtualNetwork struct {
	AzureObjectMeta
}
type AzureSubnet struct {
	AzureObjectMeta
}
type AzureSecurityGroup struct {
	AzureObjectMeta
}
type AzureVMSize struct {
	AzureObjectMeta
}
type AzureImageReference struct {
	AzureObjectMeta
	Publisher string `json:"publisher,omitempty"`
	Offer     string `json:"offer,omitempty"`
	Sku       string `json:"sku,omitempty"`
	Version   string `json:"version,omitempty"`
}
type AzureNetworkInterface struct {
	AzureObjectMeta
	Primary          bool               `json:"primary"`
	SecurityGroup    AzureSecurityGroup `json:"securityGroup"`
	IPConfigurations []AzureIpConfig    `json:"ipConfigurations"`
}
type AzurePublicIP struct {
	AzureObjectMeta
	PublicIPAddress string `json:"publicIPAddress"`
}
type AzureIpConfig struct {
	AzureObjectMeta
	Primary          bool           `json:"primary"`
	Subnet           AzureSubnet    `json:"subnet"`
	PrivateIPAddress string         `json:"privateIPAddress"`
	PublicIP         *AzurePublicIP `json:"publicIPAddress,omitempty"`
}
type AzureDisk struct {
	AzureObjectMeta
	SizeGB  int32  `json:"sizeGB"`
	SkuName string `json:"skuName"`         // ssd/hdd and theirs upper type
	Owner   string `json:"owner,omitempty"` // when cleanup, only controller created can be delete
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Tag describes machine tags history
//
// Tag are non-namespaced; the id of the tag
// according to etcd is in ObjectMeta.Name.
type Tag struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Values []string `json:"values"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TagList is a collection of tag.
type TagList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of Tags
	Items []Tag `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Config describes login information, ssh keys
//
// Config are non-namespaced; the id of the login
// according to etcd is in ObjectMeta.Name.
type Config struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Values map[string][]byte `json:"values"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigList is a collection of login.
type ConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of Configs
	Items []Config `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// outside

type ClusterNetwork struct {
	Type        NetworkType `json:"type"`
	ClusterCIDR string      `json:"clusterCIDR"`
}

type ClusterAuth struct {
	KubeUser     string `json:"kubeUser"`
	KubePassword string `json:"kubePassword"`
	KubeToken    string `json:"kubeToken"`
	KubeCertPath string `json:"kubeCertPath"`
	KubeCAData   []byte `json:"kubeCAData"`
	EndpointIP   string `json:"endpointIP"`
	EndpointPort string `json:"endpointPort"`
}

type ClusterVersions struct {
	MasterSets map[string]string
	NodeSets   MachineVersions
}

type ClusterCondition struct {
	Type               ClusterConditionType `json:"type"`
	Status             ConditionStatus      `json:"status"`
	LastHeartbeatTime  metav1.Time          `json:"lastHeartbeatTime"`
	LastTransitionTime metav1.Time          `json:"lastTransitionTime"`
	Reason             string               `json:"reason"`
	Message            string               `json:"message"`
}

type ClusterRatio struct {
	CpuOverCommitRatio    float64 `json:"cpuOverCommitRatio"`
	MemoryOverCommitRatio float64 `json:"memoryOverCommitRatio"`
}

type MachineThumbnail struct {
	Name   string       `json:"name"`
	Status MachinePhase `json:"status"`
}

type MachineAuth struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Key      string `json:"key"`
}

type NodeAddress struct {
	Type    NodeAddressType `json:"type"`
	Address string          `json:"address"`
}

type MachineNodeStatus struct {
	Unschedulable bool            `json:"unschedulable"`
	Conditions    []NodeCondition `json:"conditions"`
}

type NodeCondition = corev1.NodeCondition

type OperationLog struct {
	Time     metav1.Time       `json:"time"`
	Operator string            `json:"operator"`
	Type     OperationLogType  `json:"type"`
	Field    OperationLogField `json:"field"`
	Value    string            `json:"value"`
	Detail   string            `json:"detail"`
}

// little types

type CloudProvider string

type ClusterPhase string
type MachinePhase string
type PodPhase string
type MASGPhase string // machine auto scaling group phase

type ClusterConditionType string
type NodeConditionType = corev1.NodeConditionType
type ConditionStatus = corev1.ConditionStatus

type NetworkType string

type NodeAddressType string

type ResourceName string

type OperationLogType string
type OperationLogField string

type MachineVersions map[string]string // TODO
type MachineHardware map[string]string // TODO

// agent

type MachineEnvironment struct {
	SystemInfo         MachineSystemInfo   `json:"systemInfo"`
	HardwareInfo       MachineHardwareInfo `json:"hardwareInfo"`
	DiskInfo           []MachineDiskInfo   `json:"diskInfo"`
	NicInfo            []MachineNicInfo    `json:"nicInfo"`
	GPUInfo            []MachineGPUInfo    `json:"gpuInfo"`
	LastTransitionTime metav1.Time         `json:"lastTransitionTime"`
}

type MachineSystemInfo struct {
	BootTime        uint64 `json:"bootTime"`
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformFamily  string `json:"platformFamily"`
	PlatformVersion string `json:"platformVersion"`
	KernelVersion   string `json:"kernelVersion"`
}

type MachineHardwareInfo struct {
	CPUModel         string  `json:"cpuModel"`
	CPUArch          string  `json:"cpuArch"`
	CPUMHz           float64 `json:"cpuMHz"`
	CPUCores         int     `json:"cpuCores"`
	CPUPhysicalCores int     `json:"cpuPhysicalCores"`
	MemoryTotal      uint64  `json:"memoryTotal"`
}

type MachineDiskInfo struct {
	Device     string `json:"device"`
	Capacity   uint64 `json:"capacity"`
	Type       string `json:"type"`
	MountPoint string `json:"mountPoint"`
}

type MachineNicInfo struct {
	Name         string   `json:"name"`
	MTU          string   `json:"mtu"`
	Speed        string   `json:"speed"`
	HardwareAddr string   `json:"hardwareAddr"`
	Status       string   `json:"status"`
	Addrs        []string `json:"addrs"`
}

type MachineGPUInfo struct {
	UUID             string `json:"uuid"`
	ProductName      string `json:"productName"`
	ProductBrand     string `json:"productBrand"`
	PCIeGen          string `json:"pcieGen"`
	PCILinkWidths    string `json:"pciLinkWidths"`
	MemoryTotal      string `json:"memoryTotal"`
	MemoryClock      string `json:"memoryClock"`
	GraphicsAppClock string `json:"graphicsAppClock"`
	GraphicsMaxClock string `json:"graphicsMaxClock"`
}

// auto scaling

// ClusterScaleUpSetting describe cluster scale up setting
type ClusterScaleUpSetting struct {
	Algorithm            string `json:"algorithm"`
	IsQuotaUpdateEnabled bool   `json:"isQuotaUpdateEnabled"`
}

// ClusterScaleDownSetting describe cluster scale down setting
type ClusterScaleDownSetting struct {
	// is scale down enabled
	IsEnabled bool `json:"isEnabled"`
	// cool down time after any cluster scale up action
	CoolDown metav1.Duration `json:"coolDown"`
	// machine continues idle time threshold
	IdleTime metav1.Duration `json:"idleTime"`
	// machine idle threshold, percent of cpu/mem usage
	IdleThreshold int `json:"idleThreshold"`
}

// AutoScalingNotifySetting describe notify about setting
type AutoScalingNotifySetting struct {
	// notify methods
	Methods []string `json:"methods"`
	// notify user ids
	Users []string `json:"users"`
}

// ClusterAutoScalingSetting describe a cluster auto scaling setting
// maybe nil in old or control cluster, need to be inited with a default setting by controller
type ClusterAutoScalingSetting struct {
	// scale up setting
	ScaleUpSetting ClusterScaleUpSetting `json:"scaleUpSetting"`
	// scale down setting
	ScaleDownSetting ClusterScaleDownSetting `json:"scaleDownSetting"`
	// cluster level warning message notify setting
	NotifySetting AutoScalingNotifySetting `json:"notifySetting"`
}

// ClusterAutoScalingStatus describe cluster auto scaling operate status
type ClusterAutoScalingStatus struct {
	// last scale up operation time
	LastScaleUpTime metav1.Time `json:"lastScaleUpTime,omitempty"`
	// last selected scale up group name
	LastScaleUpGroup string `json:"lastScaleUpGroup,omitempty"`
	// last scale down operation time
	LastScaleDownTime metav1.Time `json:"lastScaleDownTime,omitempty"`
	// last selected scale down group name
	LastScaleDownGroup string `json:"lastScaleDownGroup,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineAutoScalingGroup describe a machine auto scaling group
type MachineAutoScalingGroup struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   MASGSpec   `json:"spec"`
	Status MASGStatus `json:"status"`
}

// MASGSpec describe MachineAutoScalingGroup spec
type MASGSpec struct {
	// is this auto scaling group enabled
	IsEnabled bool `json:"isEnabled"`
	// machine auto scaling group cloud provider
	Provider CloudProvider `json:"provider"`
	// machine auto scaling group cloud provider config
	ProviderConfig MASGProviderConfig `json:"providerConfig"`
	// tags of scaled machine
	Tags map[string]string `json:"tags"`
	// cluster name which group belongs to
	Cluster string `json:"cluster"`
	// group min machine num
	MinNum int `json:"minNum"`
	// group max machine num
	MaxNum int `json:"maxNum"`
	// group level warning message notify setting
	NotifySetting AutoScalingNotifySetting `json:"notifySetting"`
}

// MASGProviderConfig describe MachineAutoScalingGroup provider config
type MASGProviderConfig struct {
	Azure *AzureMASGProviderConfig `json:"azure"`
}

// MASGProviderAzureConfig describe machine auto scaling group provider config for azure
// similar with AzureMachineCloudProviderConfig, but no nic inside
type AzureMASGProviderConfig struct {
	AzureObjectMeta
	VirtualNetwork AzureVirtualNetwork `json:"virtualNetwork"`
	Subnet         AzureSubnet         `json:"subnet"`
	LoginUser      string              `json:"loginUser"`
	LoginPassword  string              `json:"loginPassword"`
	VMSize         AzureVMSize         `json:"vmSize"`
	ImageReference AzureImageReference `json:"imageReference"`
	OSDisk         AzureDisk           `json:"osDisk"`
	DataDisks      []AzureDisk         `json:"dataDisks"`
	SecurityGroup  AzureSecurityGroup  `json:"securityGroup"`
}

// MASGStatus describe MachineAutoScalingGroup status
type MASGStatus struct {
	// machine auto scaling group status phase
	Phase MASGPhase `json:"phase"`
	// info of machines belong to this group
	Machines []MASGMachineInfo `json:"machines"`
	// last scale up operation time
	LastScaleUpTime metav1.Time `json:"lastScaleUpTime,omitempty"`
	// last selected scale up machine name
	LastScaleUpMachine string `json:"lastScaleUpMachine,omitempty"`
	// last scale down operation time
	LastScaleDownTime metav1.Time `json:"lastScaleDownTime,omitempty"`
	// last selected scale down machine name
	LastScaleDownMachine string `json:"lastScaleDownMachine,omitempty"`
}

// MASGMachineInfo saves info of machine which belongs to this MachineAutoScalingGroup
type MASGMachineInfo struct {
	// name of related machine
	Name string `json:"name"`

	// machine provider config
	ProviderConfig MASGMachineProviderConfig `json:"providerConfig"`

	// scaling up about

	// timestamp when vm created
	// if nil or 0, means machine not created yet
	CreatedTime *metav1.Time `json:"createdTime,omitempty"`
	// timestamp when vm bounded to cluster
	// if nil or 0, means machine not bound to cluster yet
	BoundTime *metav1.Time `json:"boundTime,omitempty"`

	// scaling down about

	// lastBusyTime mark the last time when enough pods run on this machine
	// ignore if nil or not greater than boundTime
	LastBusyTime *metav1.Time `json:"lastBusyTime,omitempty"`

	// timestamp when set unbound
	// if nil or 0, means machine still bound, creating or failed
	UnboundTime *metav1.Time `json:"unboundTime,omitempty"`
}

// MASGMachineProviderConfig saves inited provider config of machine in auto scaling group
type MASGMachineProviderConfig struct {
	// provider config for azure
	Azure *MASGMachineAzureProviderConfig `json:"azure,omitempty"`
}

// MASGMachineAzureProviderConfig is MASGMachineProviderConfig in azure
type MASGMachineAzureProviderConfig struct {
	// azure machine object name is generated by vm resource group and name, so we need save vm name first
	VMName string `json:"vmName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineAutoScalingGroupList is a collection of machine auto scaling groups
type MachineAutoScalingGroupList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of MachineAutoScalingGroups
	Items []MachineAutoScalingGroup `json:"items" protobuf:"bytes,2,rep,name=items"`
}
