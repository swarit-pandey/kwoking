package ced

import (
	"time"

	siaproto "github.com/accuknox/shared-informer-agent/protoP"
)

func (cc *CEDCore) simulateNodes(op Operation) []string {
	var simulatedNodes []*siaproto.NodeDetails
	resourcePrefix := "simulated-node"
	numSimulatedNodes := cc.config.Simulate.Nodes.Base
	if numSimulatedNodes == 0 {
		numSimulatedNodes = 50 // By default we simulate a cluster of 50 nodes
	}
	var allSimulatedNodes []string

	for i := 0; i < numSimulatedNodes; i++ {
		simulateNodeName := getUniqueName(resourcePrefix)
		simulateOldNodeName := ""
		if op == Updated {
			simulateOldNodeName = getUniqueName(resourcePrefix)
		}
		simulatedNodes = append(simulatedNodes, &siaproto.NodeDetails{
			ClusterId:       cc.config.Simulate.ClusterID,
			WorkspaceId:     cc.config.Simulate.WorkspaceID,
			Operation:       string(op),
			Labels:          cc.getRandomLabels(5, 10),
			LastUpdatedTime: time.Now().UTC().String(),
			NewNodeName:     simulateNodeName,
			OldNodeName:     simulateOldNodeName,
		})
		allSimulatedNodes = append(allSimulatedNodes, simulateNodeName)
	}

	cc.resources.nodes = simulatedNodes
	return allSimulatedNodes
}

func (cc *CEDCore) simulateNamespaces(op Operation) []string {
	var simulatedNamespaces []*siaproto.NamespaceDetails
	resourcePrefix := "simulated-namespace"
	numSimulatedNamespaces := cc.config.Simulate.Namespaces.Base
	if numSimulatedNamespaces == 0 {
		numSimulatedNamespaces = 50
	}
	var allSimulatedNamespaces []string

	for i := 0; i < numSimulatedNamespaces; i++ {
		simulatedNSName := getUniqueName(resourcePrefix)
		simulateOldNSName := ""
		if op == Updated {
			simulateOldNSName = getUniqueName(resourcePrefix)
		}
		allSimulatedNamespaces = append(allSimulatedNamespaces, simulatedNSName)
		simulatedNamespaces = append(simulatedNamespaces, &siaproto.NamespaceDetails{
			ClusterId:               cc.config.Simulate.ClusterID,
			WorkspaceId:             cc.config.Simulate.WorkspaceID,
			Operation:               string(op),
			Labels:                  cc.getRandomLabels(1, 5),
			LastUpdatedTime:         time.Now().UTC().String(),
			NewNamespaceName:        simulatedNSName,
			OldNamespaceName:        simulateOldNSName,
			KubearmorFilePosture:    "kubearmor-file-posture",
			KubearmorNetworkPosture: "kubearmor-network-posture",
		})
	}

	cc.resources.namespaces = simulatedNamespaces
	return allSimulatedNamespaces
}

func (cc *CEDCore) simulatePods(op Operation, nsName, nodeName string) {
	var simualatedPods []*siaproto.PodDetails
	resourcePrefix := "simulated-pods"
	numSimulatedPods := cc.config.Simulate.Pods.Base
	if numSimulatedPods == 0 {
		numSimulatedPods = 50
	}

	for i := 0; i < numSimulatedPods; i++ {
		simulatedContainers := cc.resources.containers
		if len(simulatedContainers) == 0 {
			simulatedContainers = cc.simulatedContainers()
		}
		simulatedPodName := getUniqueName(resourcePrefix)
		simulateOldPodName := ""
		if op == Updated {
			simulateOldPodName = getUniqueName(resourcePrefix)
		}
		simualatedPods = append(simualatedPods, &siaproto.PodDetails{
			ClusterId:       cc.config.Simulate.ClusterID,
			WorkspaceId:     cc.config.Simulate.WorkspaceID,
			Operation:       string(op),
			Labels:          cc.getRandomLabels(5, 10),
			LastUpdatedTime: time.Now().UTC().String(),
			NewPodName:      simulatedPodName,
			OldPodName:      simulateOldPodName,
			Namespace:       nsName,
			NodeName:        nodeName,
			Container:       simulatedContainers,
			WorkloadType:    getOneOf([]string{"job", "cronjob", "deployment", "replicaset"}),
		})
	}

	cc.resources.pods = simualatedPods
}

func (cc *CEDCore) simulateWorkloads(op Operation, namespaceName string) {
	var simulatedWorkloads []*siaproto.Workload
	resourcePrefix := "simulated-workload"
	numSimulatedWorkloads := cc.config.Simulate.Workloads.Base
	if numSimulatedWorkloads == 0 {
		numSimulatedWorkloads = 10
	}

	for i := 0; i < numSimulatedWorkloads; i++ {
		simulatedWorloadName := getUniqueName(resourcePrefix)
		simulateOldWorkloadName := ""
		if op == Updated {
			simulateOldWorkloadName = getUniqueName(resourcePrefix)
		}

		simulatedWorkloads = append(simulatedWorkloads, &siaproto.Workload{
			ClusterId:       cc.config.Simulate.ClusterID,
			WorkspaceId:     cc.config.Simulate.WorkspaceID,
			Namespace:       namespaceName,
			Type:            getOneOf([]string{"job", "cronjob", "deployment", "replicaset"}),
			NewName:         simulatedWorloadName,
			OldName:         simulateOldWorkloadName,
			Operation:       string(op),
			Labels:          cc.getRandomLabels(5, 10),
			LastUpdatedTime: time.Now().UTC().String(),
		})
	}

	cc.resources.workloads = simulatedWorkloads
}

func (cc *CEDCore) simulatedContainers() []*siaproto.ContainerDetails {
	var simulatedContainers []*siaproto.ContainerDetails
	resourcePrefix := "simulated-container"
	numSimulatedContainers := 5

	for i := 0; i < numSimulatedContainers; i++ {
		simulatedContainerName := getUniqueName(resourcePrefix)
		simulatedContainers = append(simulatedContainers, &siaproto.ContainerDetails{
			ContainerName: simulatedContainerName,
			NameOfService: getUniqueName("simulated-container-service"),
			Image:         getUniqueName("simulated-container-image"),
			Status:        getOneOf([]string{"running", "suspended", "terminated"}),
			ContainerId:   getUniqueName("simulated-container-id"),
		})
	}

	cc.resources.containers = simulatedContainers
	return simulatedContainers
}
