package utils

type contextKey string

const (
	AuthedUsername  contextKey = "authedUsername"
	AuthedUserIdHex contextKey = "authedUserIdHex"
	RequestIdKey    contextKey = "requestIdKey"
)
