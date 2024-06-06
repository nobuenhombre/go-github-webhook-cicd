package queue

type Service interface {
	Push(item interface{}) error
	Run() error
	Stop()
}
