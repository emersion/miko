package action

type Request struct {
	*Action
	wait chan error
}

func (r *Request) Wait() error {
	return <-r.wait
}

func (r *Request) Done(err error) {
	select {
	case r.wait <- err:
	default:
	}
}

func newRequest(action *Action) *Request {
	return &Request{
		Action: action,
		wait:   make(chan error),
	}
}
