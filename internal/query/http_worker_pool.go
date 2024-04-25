package query

// max amount of tasks in worker pool queue
const maxTasksQueue = 10

// WorkerPool interface for any Query Engine worker pools
type WorkerPool interface {
	Query(req *HTTPRequest) *HTTPResponse
}

// HTTPWorkerPool structure that contain WokerPool HTTP implementation
// It implements worker pool that will do real HTTP calls to resources using `query.MakeHTTPRequest`
type HTTPWorkerPool struct {
	workerCount int
	input       chan *asyncHTTPRequest
}

type asyncHTTPRequest struct {
	request  *HTTPRequest
	response chan *HTTPResponse
}

// NewHTTPWorkerPool create new worker pool for queries
func NewHTTPWorkerPool(workerCount int) *HTTPWorkerPool {
	wp := &HTTPWorkerPool{
		workerCount: workerCount,
		input:       make(chan *asyncHTTPRequest, maxTasksQueue),
	}

	for w := 0; w < wp.workerCount; w++ {
		go wp.worker()
	}

	return wp
}

// Query makes request to given Request
// Under the hood it will wrap everything to async query and execute it using
// worker pool.
func (wp *HTTPWorkerPool) Query(req *HTTPRequest) *HTTPResponse {
	asyncReq := &asyncHTTPRequest{
		request:  req,
		response: make(chan *HTTPResponse),
	}
	// Sending request
	wp.input <- asyncReq
	// Waiting for response
	res := <-asyncReq.response
	// Have to close channel
	close(asyncReq.response)
	return res
}

func (wp *HTTPWorkerPool) worker() {
	for req := range wp.input {
		req.response <- MakeHTTPRequest(req.request)
	}
}
