package query

// MockWorkerPool mock worker pool implementation for tests
type MockWorkerPool struct {
	resp     *HTTPResponse
	checkReq func(*HTTPRequest)
}

func NewMockWorkerPool() *MockWorkerPool {
	return &MockWorkerPool{}
}

func (mwp *MockWorkerPool) MockResp(resp *HTTPResponse) {
	mwp.resp = resp
}

func (mwp *MockWorkerPool) MockBody(body string) {
	mwp.resp = &HTTPResponse{
		Body: []byte(body),
	}
}
func (mwp *MockWorkerPool) SetRequestAssertions(f func(*HTTPRequest)) {
	mwp.checkReq = f
}

func (mwp *MockWorkerPool) Query(req *HTTPRequest) *HTTPResponse {
	if mwp.checkReq != nil {
		mwp.checkReq(req)
	}
	return mwp.resp
}
