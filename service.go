package grpcx

//ServiceOption service options
type ServiceOption func(*serviceOptions)

type serviceOptions struct {
	httpEntrypoints []*httpEntrypoint
}

type httpEntrypoint struct {
	path    string
	method  string
	invoker func([]byte) ([]byte, error)
}

// Service is a service define
type Service struct {
	Name     string
	Metadata interface{}
	opts     *serviceOptions
}

// NewService returns a Service
func NewService(name string, metadata interface{}, opts ...ServiceOption) Service {
	sopts := &serviceOptions{}
	for _, opt := range opts {
		opt(sopts)
	}

	service := Service{
		Name:     name,
		Metadata: metadata,
		opts:     sopts,
	}

	return service
}

// WithAddHTTPEntrypoint add a http service metadata
func WithAddHTTPEntrypoint(path, method string, invoker func([]byte) ([]byte, error)) ServiceOption {
	return func(opt *serviceOptions) {
		opt.httpEntrypoints = append(opt.httpEntrypoints, &httpEntrypoint{
			path:    path,
			method:  method,
			invoker: invoker,
		})
	}
}
