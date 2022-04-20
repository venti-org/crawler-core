package pipeline

type Pipeline interface {
	AppendComponents(components ...Component) Pipeline
	ProcessItem(item Item)
	GetComponents() []Component
	Empty() bool
}

type DefaultPipeline struct {
	components []Component
}

func NewDefaultPipeline(components ...Component) *DefaultPipeline {
	p := &DefaultPipeline{}
	return p
}

func (p *DefaultPipeline) AppendComponents(components ...Component) Pipeline {
	p.components = append(p.components, components...)
	return p
}

func (p *DefaultPipeline) ProcessItem(item Item) {
	for _, component := range p.components {
		component.ProcessItem(item)
	}
}

func (p *DefaultPipeline) GetComponents() []Component {
	return p.components
}

func (p *DefaultPipeline) Empty() bool {
	return len(p.components) == 0
}
