package shutdown

type Service interface {
	Start()

	Stop()
}
