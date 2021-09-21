package schemas

import (
	"bytes"
	"hash/fnv"
	"sort"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/desc/protoprint"
	"google.golang.org/protobuf/types/descriptorpb"
)

type DAGsBuilder struct {
	mergedDAGs map[string]Node
}

func NewDAGsBuilder() *DAGsBuilder {
	mergedDAGs := make(map[string]Node)
	return &DAGsBuilder{mergedDAGs: mergedDAGs}
}

func (d *DAGsBuilder) AddDag(fileDescriptor *desc.FileDescriptor, subject string) error {
	nodes, err := ConstructDAG(fileDescriptor, subject)
	if err != nil {
		return err
	}

	d.mergeDAG(nodes)
	return nil
}

func (d *DAGsBuilder) mergeDAG(dag map[string]Node) {
	for name, node := range dag {
		if mergedNode, ok := d.mergedDAGs[name]; ok {
			if mergedNode.Level >= node.Level {
				continue
			}
		}
		d.mergedDAGs[name] = node
	}
}

func (d *DAGsBuilder) Records() (records Records, err error) {
	nodeList := GetDepOrder(d.mergedDAGs)

	records, err = GetResults(nodeList)
	if err != nil {
		return records, err
	}

	return records, nil
}

type Node struct {
	Subject    string
	Descriptor *desc.FileDescriptor
	References []string
	// Is this node a dependency of another node?
	Dependency bool
	// Level in the graph
	Level int
	// LevelHash contains the hash of the subject + filepath of the content so we can make the level unique
	// in order to achieve a stable sort
	// We use subject + filepath because in the same level we might have the same file twice one having the
	// subject being the topic name and the other being the filepath itself, which could lead to the same hash.
	LevelHash uint32
}

func ConstructDAG(fileDescriptor *desc.FileDescriptor, subject string) (nodes map[string]Node, err error) {
	nodes = make(map[string]Node)
	return nodes, constructDAGRecursively(fileDescriptor, nodes, subject, false, 0)
}

func constructDAGRecursively(fileDescriptor *desc.FileDescriptor, nodes map[string]Node, subject string, dependency bool, level int) (err error) {
	depsName := []string{}

	// Check if we already have it
	if node, ok := nodes[subject]; ok {
		if node.Level >= level {
			return nil
		}
		node.Level = level
		nodes[subject] = node
		return nil
	}

	deps := fileDescriptor.GetDependencies()
	for _, dep := range deps {
		depsName = append(depsName, dep.GetName())
		err = constructDAGRecursively(dep, nodes, dep.GetName(), true, level+1)
		if err != nil {
			return err
		}
	}

	fileBuilder, err := builder.FromFile(fileDescriptor)
	if err != nil {
		return err
	}

	// fileBuilder.SetOptions(&descriptorpb.FileOptions{})

	fileDescriptor, err = fileBuilder.Build()
	if err != nil {
		return err
	}

	nodes[subject] = Node{
		Subject:    subject,
		Descriptor: fileDescriptor,
		References: depsName,
		Dependency: dependency,
		Level:      level,
		LevelHash:  hash(subject + fileDescriptor.GetName()),
	}
	return nil
}

func MergeDAGs(dags []map[string]Node) map[string]Node {
	mergedResult := make(map[string]Node)

	for _, dag := range dags {
		for name, node := range dag {
			if mergedNode, ok := mergedResult[name]; ok {
				if mergedNode.Level >= node.Level {
					continue
				}
			}
			mergedResult[name] = node
		}
	}

	return mergedResult
}

// GetDepOrder returns a list of Nodes in order they should be inserted.
func GetDepOrder(nodes map[string]Node) []Node {
	nodeList := make([]Node, 0, len(nodes))

	for _, node := range nodes {
		nodeList = append(nodeList, node)
	}

	sort.SliceStable(nodeList, func(i, j int) bool {
		if nodeList[i].Level > nodeList[j].Level {
			return true
		} else if nodeList[i].Level == nodeList[j].Level {
			return int(nodeList[i].LevelHash) > int(nodeList[j].LevelHash)
		}
		return false
	})

	return nodeList
}

func GetResults(nodeList []Node) (Records, error) {
	records := Records{}

	// Sort fields to achieve consistent output
	// Force FQDNs so we can properly map dependencies
	printer := protoprint.Printer{
		SortElements:             true,
		OmitComments:             protoprint.CommentsAll,
		Compact:                  true,
		ForceFullyQualifiedNames: false,
	}

	for _, node := range nodeList {
		buf := bytes.Buffer{}

		err := printer.PrintProtoFile(node.Descriptor, &buf)
		if err != nil {
			return records, err
		}

		record := Record{Subject: node.Subject, Schema: buf.String(), References: node.References}
		records = append(records, record)
	}

	return records, nil
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
