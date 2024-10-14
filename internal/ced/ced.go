package ced

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	siaproto "github.com/accuknox/shared-informer-agent/protoP"
	"github.com/swarit-pandey/kwoking/internal/common"
	"github.com/swarit-pandey/kwoking/internal/config"
	"github.com/swarit-pandey/kwoking/internal/rmq"
)

type Operation string

const (
	Added   Operation = "added"
	Updated Operation = "updated"
	Delete  Operation = "deleted"
)

type resources struct {
	nodes      []*siaproto.NodeDetails
	namespaces []*siaproto.NamespaceDetails
	pods       []*siaproto.PodDetails
	containers []*siaproto.ContainerDetails
	workloads  []*siaproto.Workload
}

type CEDCore struct {
	ctx       context.Context
	config    *config.CED
	publisher *rmq.Publisher
	resources *resources
	labels    []*siaproto.Labels
}

func NewCEDCore(ctx context.Context, config *config.CED) *CEDCore {
	cc := &CEDCore{
		config:    config,
		ctx:       ctx,
		resources: &resources{},
	}

	if !cc.config.DryRun {
		cc.publisher = cc.setupPublisher()
	}
	cc.generateAndCacheLabels(50)

	return cc
}

func (cc *CEDCore) setupPublisher() *rmq.Publisher {
	slog.Info("Setting up publisher")
	pubOpts := rmq.Opts{
		// Exchange parameters
		ExchangeName: cc.config.Broker.Exchange.Name,
		ExchangeType: cc.config.Broker.Exchange.Type,
		RoutingKey:   cc.config.Broker.Exchange.Key,
		Durable:      cc.config.Broker.Exchange.Durable,
		AutoDeleted:  cc.config.Broker.Exchange.AutoDeleted,
		Internal:     cc.config.Broker.Exchange.Internal,
		ContentType:  cc.config.Broker.Exchange.ContentType,
		NoWait:       cc.config.Broker.Exchange.NoWait,

		// Connection parameters
		Username:   cc.config.Broker.Connection.Username,
		Password:   cc.config.Broker.Connection.Password,
		Connection: cc.config.Broker.Connection.URL,
	}

	return rmq.NewRMQPublisher(&pubOpts)
}

func (cc *CEDCore) Simulate() {
	cc.setupSimulation()

	if !cc.config.DryRun {
		slog.Info("Preparing and flushing data")
		cc.prepareAndFlushWorkloads()
		cc.prepareAndFlushNS()
		cc.prepareAndFlushPods()
		cc.prepareAndFlushNode()
	}
}

func (cc *CEDCore) setupSimulation() {
	if cc.config == nil {
		slog.Error("Config is not initialized")
	}

	simulatedNodes := cc.simulateNodes(Added)
	simulatedNamespaces := cc.simulateNamespaces(Added)

	for _, node := range simulatedNodes {
		for _, ns := range simulatedNamespaces {
			cc.simulatePods(Added, node, ns)
			cc.simulateWorkloads(Added, ns)
		}
	}

	totalNodes := len(cc.resources.nodes)
	totalNamespaces := len(cc.resources.namespaces)
	totalPods := totalNodes * totalNamespaces * cc.config.Simulate.Pods.Base
	totalWorkloads := totalNamespaces * cc.config.Simulate.Workloads.Base
	totalContainers := totalPods * 5
	totalResources := totalNodes + totalNamespaces + totalPods + totalWorkloads + totalContainers

	cc.logSimulationSummary(totalNodes, totalNamespaces, totalPods, totalWorkloads, totalContainers, totalResources)

	if cc.config.Simulate.Log {
		cc.saveSimulationsAsJSON()
	}
}

func (cc *CEDCore) flush(data []byte) {
	if cc.config.DryRun {
		return
	}

	err := cc.publisher.Publish(data)
	if err != nil {
		slog.Error("failed to publish data", err)
	}
}

func (cc *CEDCore) generateAndCacheLabels(numLabels int) {
	generatedLabels := common.GenerateRandomLabels(numLabels, 2)
	cc.labels = []*siaproto.Labels{}

	for _, labelMap := range generatedLabels {
		for k, v := range labelMap {
			cc.labels = append(cc.labels, &siaproto.Labels{Key: k, Value: v})
		}
	}
}

func (cc *CEDCore) getRandomLabels(u, l int) []*siaproto.Labels {
	if l < 1 {
		l = 1
	}
	if u > len(cc.labels) {
		l = len(cc.labels)
	}
	if l > u {
		l = u
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	fn := l + rnd.Intn(u-l+1)

	shuffled := make([]*siaproto.Labels, len(cc.labels))
	copy(shuffled, cc.labels)
	rnd.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
	return shuffled[:fn]
}
