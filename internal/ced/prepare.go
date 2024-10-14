package ced

import (
	"log/slog"

	siaproto "github.com/accuknox/shared-informer-agent/protoP"
	"google.golang.org/protobuf/proto"
)

func (cc *CEDCore) prepareAndFlushNode() {
	for _, node := range cc.resources.nodes {
		singleNode := []*siaproto.NodeDetails{node}
		details := siaproto.EntityDetails{
			EntityType:          "node",
			NamespaceList:       nil,
			PodList:             nil,
			NodeList:            &siaproto.MultipleNodeDetails{Nodes: singleNode},
			CurrentResourceList: nil,
			Workloads:           nil,
			Services:            nil,
			Heartbeat:           nil,
		}

		detailsInBytes, err := proto.Marshal(&details)
		if err != nil {
			slog.Error("failed to flush event", slog.Any("details", node), slog.Any("error", err))
			continue
		}

		cc.flush(detailsInBytes)
	}
}

func (cc *CEDCore) prepareAndFlushNS() {
	for _, namespace := range cc.resources.namespaces {
		singleNS := []*siaproto.NamespaceDetails{namespace}
		details := siaproto.EntityDetails{
			EntityType:          "namespace",
			NamespaceList:       &siaproto.MultipleNamespaceDetails{Namespaces: singleNS},
			PodList:             nil,
			NodeList:            nil,
			CurrentResourceList: nil,
			Workloads:           nil,
			Services:            nil,
			Heartbeat:           nil,
		}

		detailsInBytes, err := proto.Marshal(&details)
		if err != nil {
			slog.Error("failed to flush event", slog.Any("details", namespace), slog.Any("error", err))
			continue
		}

		cc.flush(detailsInBytes)
	}
}

func (cc *CEDCore) prepareAndFlushWorkloads() {
	for _, workload := range cc.resources.workloads {
		singleWorkload := []*siaproto.Workload{workload}
		details := siaproto.EntityDetails{
			EntityType:          "workload",
			NamespaceList:       nil,
			PodList:             nil,
			NodeList:            nil,
			CurrentResourceList: nil,
			Workloads:           singleWorkload,
			Services:            nil,
			Heartbeat:           nil,
		}

		detailsInBytes, err := proto.Marshal(&details)
		if err != nil {
			slog.Error("failed to flush event", slog.Any("details", singleWorkload), slog.Any("error", err))
			continue
		}

		cc.flush(detailsInBytes)
	}
}

func (cc *CEDCore) prepareAndFlushPods() {
	for _, pod := range cc.resources.pods {
		singlePod := []*siaproto.PodDetails{pod}
		details := siaproto.EntityDetails{
			EntityType:          "pod",
			NamespaceList:       nil,
			PodList:             &siaproto.MultiplePodDetails{Pods: singlePod},
			NodeList:            nil,
			CurrentResourceList: nil,
			Workloads:           nil,
			Services:            nil,
			Heartbeat:           nil,
		}

		detailsInBytes, err := proto.Marshal(&details)
		if err != nil {
			slog.Error("failed to flush event", slog.Any("details", singlePod), slog.Any("error", err))
			continue
		}

		cc.flush(detailsInBytes)
	}
}

func (cc *CEDCore) flushRefresh() {
	resourceList := &siaproto.ResourceList{
		ClusterId:      cc.config.Simulate.ClusterID,
		WorkspaceId:    cc.config.Simulate.WorkspaceID,
		NodeArray:      cc.resources.nodesList,
		NamespaceArray: cc.resources.namespacesList,
		PodArray:       cc.resources.podsList,
	}

	slog.Info("", slog.Any("len(ns)", len(cc.resources.namespacesList)), slog.Any("len(node)", len(cc.resources.nodesList)), slog.Any("len(pods)", len(cc.resources.podsList)))

	details := siaproto.EntityDetails{
		EntityType:          "currentResource",
		NamespaceList:       nil,
		PodList:             nil,
		CurrentResourceList: resourceList,
		Workloads:           nil,
		Services:            nil,
		Heartbeat:           nil,
	}

	detailsInBytes, err := proto.Marshal(&details)
	if err != nil {
		slog.Error("failed to flush event", slog.Any("details", resourceList), slog.Any("error", err))
		return
	}

	cc.flush(detailsInBytes)
}
