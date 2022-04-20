package pipeline

type Component interface {
	ProcessItem(Item)
}
