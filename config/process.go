package config

const (
	TERMINATE = "terminate"
	COMPLETE  = "complete"
	REJECT    = "reject"
	START     = "start"
	SAVE      = "save"
)

// Process 流程信息结构
type Process struct {
	Name      string
	Key       string
	Model     string
	ModelName string
	Nodes     map[string]Node
	Seq       int // 序号(用于固定字典顺序)
}

func (p *Process) GetNodes() []string {
	nodes := make([]string, 0)
	for k := range p.Nodes {
		if k != COMPLETE && k != TERMINATE {
			nodes = append(nodes, k)
		}
	}
	return nodes
}

// Node 节点类
type Node struct {
	Title string
	Batch bool // 是否支持批量处理
	Seq   int  // 序号
}

// NodeArr 用来排序
type NodeArr []Node

func (a NodeArr) Len() int           { return len(a) }
func (a NodeArr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NodeArr) Less(i, j int) bool { return a[i].Seq < a[j].Seq }
