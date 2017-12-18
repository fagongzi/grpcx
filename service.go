package grpcx

import "io"

//ServiceOption service options
type ServiceOption func(*serviceOptions)

type serviceOptions struct {
	httpPath   string
	httpMethod string

	customMetadata interface{}

	invoker func(io.Reader) ([]byte, error)
}

// HTTPService http service
type HTTPService struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

// ServiceMetadata service meta data
type ServiceMetadata struct {
	HTTP   *HTTPService `json:"http, omitempty"`
	Custom interface{}  `json:"custom, omitempty"`
}

// Service is a service define
type Service struct {
	Name     string
	Metadata ServiceMetadata
	invoker  func(io.Reader) ([]byte, error)
}

// NewService returns a Service
func NewService(name string, opts ...ServiceOption) Service {
	sopts := &serviceOptions{}
	for _, opt := range opts {
		opt(sopts)
	}

	service := Service{
		Name: name,
	}

	if sopts.httpMethod != "" && sopts.httpPath != "" {
		service.Metadata.HTTP = &HTTPService{
			Path:   sopts.httpPath,
			Method: sopts.httpMethod,
		}
	}

	if sopts.invoker != nil {
		service.invoker = sopts.invoker
	}

	if sopts.customMetadata != nil {
		service.Metadata.Custom = sopts.customMetadata
	}

	return service
}

// WithCustomMetadata returns a custom meta data
func WithCustomMetadata(metadata interface{}) ServiceOption {
	return func(opt *serviceOptions) {
		opt.customMetadata = metadata
	}
}

// WithHTTPMetadata returnd a http service metadata
func WithHTTPMetadata(path, method string) ServiceOption {
	return func(opt *serviceOptions) {
		opt.httpPath = path
		opt.httpMethod = method
	}
}

// WithGRPCInvoker with GRPC invoker
func WithGRPCInvoker(invoker func(io.Reader) ([]byte, error)) ServiceOption {
	return func(opt *serviceOptions) {
		opt.invoker = invoker
	}
}
