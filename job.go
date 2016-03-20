package pool

import "net/http"

// JobInterface basic interface for a pool queue job
type JobInterface interface {
	GetError() *error
	SetError(e *error)
	GetResponse() ResponseInterface
	SetResponse(r ResponseInterface)
	GetRequest() *http.Request
	SetRequest(r *http.Request)
}

// Job a basic implementation of the JobInterface for the request pool
type Job struct {
	error    *error
	response ResponseInterface
	Request  *http.Request
}
func (j *Job) GetError() *error                { return j.error }
func (j *Job) SetError(e *error)               { j.error = e }
func (j *Job) GetResponse() ResponseInterface  { return j.response }
func (j *Job) SetResponse(r ResponseInterface) { j.response = r }
func (j *Job) GetRequest() *http.Request       { return j.Request }
func (j *Job) SetRequest(r *http.Request)      { j.Request = r }
