## http-pool

A simple http pool build around the net/http.

```
	p := pool.NewPool(10);
	jobs := make(chan pool.JobInterface, 10)
	for i := 0; i < 10; i++ {
		ret, _ := http.NewRequest("GET", "http://www.example.com", nil)
		jobs <- &pool.Job{Request: ret}
	}
	
	close(jobs)
	defer p.Close()
	p.Start(jobs, 0)

	for p.MatchOneState(pool.RUNNING|pool.CHAN_HAS_FAILURE|pool.CHAN_HAS_SUCCESS){
		select {
		case s := <- pool.Success:
			fmt.Printf("Finished url: %s\n", s.GetRequest().URL)
			fmt.Println(s.GetResponse().(*http.Response).StatusCode)
		case e := <- pool.Error:
			fmt.Println(*e.GetError())
		}
	}
```