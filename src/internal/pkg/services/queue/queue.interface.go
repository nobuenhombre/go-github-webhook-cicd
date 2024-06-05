package queue

type Service interface {
	Push(item interface{})
	Run()
	Stop()
}
