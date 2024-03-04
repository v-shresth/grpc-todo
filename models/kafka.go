package models

type QueueTopics string

const (
	TopicDeadlineNearby QueueTopics = "deadline_nearby"
	TopicPremiumEnding  QueueTopics = "premium_ending"
)
