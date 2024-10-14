package ced

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
)

func getUniqueName(prefix string) string {
	uid := uuid.New()
	uidStr := uid.String()[:8]
	return fmt.Sprintf("%s-%s", prefix, uidStr)
}

func getOneOf(states []string) string {
	if len(states) == 0 {
		return ""
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return states[rnd.Intn(len(states))]
}

func (cc *CEDCore) saveSimulationsAsJSON() {
	if cc.resources == nil {
		fmt.Println("No resources available to save simulations")
		return
	}

	filename := fmt.Sprintf("simulated_resources_%s.json", time.Now().Format("20060102-150405"))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create file: %s\n", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "{\n")
	saveResourceAsJSON(file, "nodes", cc.resources.nodes)
	saveResourceAsJSON(file, "namespaces", cc.resources.namespaces)
	saveResourceAsJSON(file, "pods", cc.resources.pods)
	saveResourceAsJSON(file, "containers", cc.resources.containers)
	saveResourceAsJSON(file, "workloads", cc.resources.workloads)
	fmt.Fprintf(file, "}\n")
}

func saveResourceAsJSON(file *os.File, resourceType string, resources interface{}) {
	jsonData, err := json.MarshalIndent(resources, "", "    ")
	if err != nil {
		fmt.Printf("Failed to marshal %s: %s\n", resourceType, err)
		return
	}

	if _, err := fmt.Fprintf(file, "    \"%s\": ", resourceType); err != nil {
		fmt.Printf("Failed to write key for %s: %s\n", resourceType, err)
		return
	}

	if _, err := file.Write(jsonData); err != nil {
		fmt.Printf("Failed to write JSON for %s: %s\n", resourceType, err)
		return
	}

	fmt.Fprintf(file, ",\n")
}

func (cc *CEDCore) logSimulationSummary(nodes, namespaces, pods, workloads, containers, total int) {
	data := [][]string{
		{"Nodes", fmt.Sprintf("%d", nodes)},
		{"Namespaces", fmt.Sprintf("%d", namespaces)},
		{"Pods", fmt.Sprintf("%d", pods)},
		{"Workloads", fmt.Sprintf("%d", workloads)},
		{"Containers", fmt.Sprintf("%d", containers)},
		{"Total Resources", fmt.Sprintf("%d", total)},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Resource", "Count"})
	table.SetBorder(true)
	table.AppendBulk(data)
	table.Render()

	slog.Info("Simulation Summary", slog.String("total_resources", fmt.Sprintf("%d", total)))
}
