package concurrency

type ChannelResponse[T any] struct {
	Error  error
	Result T
}
