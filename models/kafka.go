package models

type QueueTopic string

const (
	TopicDeadlineNearby QueueTopic = "deadline_nearby"
	TopicPremiumEnding  QueueTopic = "premium_ending"
)
