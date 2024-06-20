## Overview

### Scheduler.go:
Reads a sequence of jobs from a YAML configuration file and executes each job concurrently while handling dependencies between jobs. The scheduler ensures that jobs with dependencies are executed only after their dependencies have completed.

### Job Scheduling

- The implementation utilizes goroutines to run jobs concurrently. 
- Each job is launched as a separate goroutine, and synchronization is achieved using a WaitGroup. 
- the "state" map keeps track of the completion status of each job, and "complete" is used to communicate the completion of jobs.

### Handing Dependencies

- Scheduler waits for the completion of all dependencies before executing a job. 
- "dependency" is used to track the number of dependencies for each job. 
- When a job completes, it notifies its dependent jobs by decrementing their dependency count. 
- If a dependent job's dependency count reaches zero, it is marked as completed and can be executed.

### Generating a DOT Graph

- The Graph function in `scheduler.go` generates the DOT graph by iterating over the jobs and their dependencies, writing the nodes and edges in the DOT language format.

### Testing
#### Test Cases
- single job execution
- multiple job execution
- job dependencies
- job delays (slowed down test compilation)
- error cases. 

- Imported bytes for efficient manipulation
- Iported assert/require to write assertions
      - Get expections and actual results
      - Terminates when asserts fails

### Add the delay configuration
- The scheduler waits for the specified delay duration before executing a job

### Debug log flag
- When the flag is set, the scheduler writes debug messages to the console
#### Flag Messages
- when a job is scheduled
- starts executing
- completion of the job.


## Features

- [X] Simple job scheduling
- [X] Job scheduling with dependencies
- [X] DOT Graph generation
- [X] Added tests
- [X] Delay configuration
- [X] Debug log flag


- imported log for debugging
- imported fmt to format string
- imported os for operating system func
- imported sync for WaitGroup
- imported time for delays and debug

## Citations/Sources
### Probably more, but these are the tabs I still have open up
https://github.com/carlescere/scheduler/blob/master/scheduler_test.go
https://medium.com/@sanilkhurana7/understanding-the-go-scheduler-and-looking-at-how-it-works-e431a6daacf
https://gobyexample.com/logging
https://freshman.tech/snippets/go/check-if-slice-contains-element/
https://stackoverflow.com/questions/1494492/graphviz-how-to-go-from-dot-to-a-graph
https://github.com/emicklei/dot
https://avi-nash5.medium.com/implementing-a-delayed-job-scheduler-in-golang-e4d2fecfa797
https://github.com/Flowpack/prunner
https://github.com/jasonlvhit/gocron

## Other Comments
- Not sure why, but doing "make test" works after two attempts (test cases fails but works later)??
- Yaml file, had to change up formatting compared to the one
given in the program specs, cause of input mismatches of the code
I already had.
