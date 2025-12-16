package domain

type TopicNode struct {
	Topic    Topic
	Children []TopicNode
}
