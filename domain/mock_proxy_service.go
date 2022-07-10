package domain

type MockProxyService struct {
	ForwardFunc func(pr *ProxyRequest) (*ProxyResponse, error)
}

func (m *MockProxyService) Forward(pr *ProxyRequest) (*ProxyResponse, error) {
	return m.ForwardFunc(pr)
}
