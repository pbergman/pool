package pool

import (
	"testing"
	"net/http"
)

func BenchmarkPool(b *testing.B) {

	pool := NewPool(uint64(b.N));
	jobs := make(chan JobInterface, b.N*2)

	for i := 0; i < b.N*2; i++ {
		ret, _ := http.NewRequest("GET", "http://www.google.com", nil)
		jobs <- &Job{Request: ret}
	}

	close(jobs)
	defer pool.Close()
	pool.Start(jobs, 0)

	for pool.MatchOneState(RUNNING|CHAN_HAS_FAILURE|CHAN_HAS_SUCCESS){
		select {
		case s := <- pool.Success:
			b.Logf("Finished url: %s(%d)\n", s.GetRequest().URL, s.GetResponse().(*http.Response).StatusCode)
		case e := <- pool.Error:
			b.Log(*e.GetError())
		}
	}
}
