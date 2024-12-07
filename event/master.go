package event

import (
	"fmt"

	"github.com/kijimaD/nova/ast"
)

// スライスで保存しつつ、キーで引ける
type LabelMaster struct {
	Labels     []Label
	LabelIndex map[string]int
}

func (master *LabelMaster) GetLabel(key string) (Label, error) {
	idx, ok := master.LabelIndex[key]
	if !ok {
		return Label{}, fmt.Errorf("keyが存在しない")
	}

	return master.Labels[idx], nil
}

func (m *LabelMaster) AddLabel(label Label) {
	_, exists := m.LabelIndex[label.Name]
	if !exists {
		m.Labels = append(m.Labels, label)
		m.LabelIndex[label.Name] = len(m.Labels) - 1
	}
}

type Label struct {
	Name string
	Body *ast.BlockStatement
}
