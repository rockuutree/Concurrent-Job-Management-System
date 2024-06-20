package scheduler

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"uw.cse374.go/scheduler/internal/config"
)

// Scheduler schedules and executes a collection
// of one or more jobs.
type Scheduler struct {
	config *config.Config
}

// New returns a new Scheduler with the
// given configuration.
func New(config *config.Config) *Scheduler {
	return &Scheduler{
		config: config,
	}
}

// Run executes the configured jobs and writes their names
// Added the parameter debug to enable debug logging
func (s *Scheduler) Run(writer io.Writer, debug bool) error {
	//  Create a wait group to wait for all jobs to complete
	var wait sync.WaitGroup
	//  Create a map to keep track of the job status
	state := make(map[string]bool)
	//  Create a channel to notify when a job is done
	complete := make(chan string)

	//  Create a logger for debug output, basically getting the time stamps
	logger := log.New(os.Stderr, "", log.LstdFlags)

	//  Initialize job status map to false
	for _, job := range s.config.Jobs {
		state[job.Name] = false
	}

	//  Start a goroutine for each job
	for _, job := range s.config.Jobs {
		wait.Add(1)
		//  Start a goroutine for each job
		go func(job config.Job) {
			defer wait.Done()

			//  Log the job execution
			if debug {
				logger.Printf("Scheduling job %q\n", job.Name)
			}

			//  Wait for the job's dependencies to complete
			for _, dependency := range job.DependsOn {
				for !state[dependency] {
					time.Sleep(100 * time.Millisecond)
				}
			}

			//  Delay the job execution if specified
			if job.Delay > 0 {
				time.Sleep(time.Duration(job.Delay) * time.Second)
			}

			//  Execute the job
			if debug {
				//  Log the job execution, and delay if specified
				if job.Delay > 0 {
					logger.Printf("Delaying job %q for %d seconds\n", job.Name, job.Delay)
				}
				//  Log the job execution
				logger.Printf("Executing job %q\n", job.Name)
			}
			fmt.Fprintln(writer, job.Name)
			//  Mark the job as completed
			state[job.Name] = true
			if debug {
				logger.Printf("Completed job %q\n", job.Name)
			}

			//  Notify that the job is done
			complete <- job.Name
		}(job) //  Pass the job to the goroutine
	}

	// Wait for all jobs to complete
	go func() {
		wait.Wait()
		close(complete)
	}()

	//  Collect the completed jobs
	var jobsCompleted []string
	for range s.config.Jobs {
		name := <-complete
		jobsCompleted = append(jobsCompleted, name)
	}

	// Ensure the jobs are printed in the order they appear
	for _, job := range s.config.Jobs {
		if !contains(jobsCompleted, job.Name) {
			fmt.Fprintln(writer, job.Name)
		}
	}

	return nil
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, val := range slice {
		if val == item {
			return true
		}
	}
	return false
}

// Graph writes a DOT-graph containing the configured jobs
// to the given writer.
func (s *Scheduler) Graph(writer io.Writer) error {
	//  Write the graph header
	_, err := io.WriteString(writer, "digraph G {\n")
	if err != nil {
		return err
	}

	// Nodes
	for _, job := range s.config.Jobs {
		//  Write the job name as a node
		_, err := fmt.Fprintf(writer, "  %q;\n", job.Name)
		if err != nil {
			return err
		}
	}

	// Edges
	for _, job := range s.config.Jobs {
		//  Write the job dependencies as edges
		for _, dependency := range job.DependsOn {
			_, err := fmt.Fprintf(writer, "  %q -> %q;\n", dependency, job.Name)
			if err != nil {
				return err
			}
		}
	}

	//  Close the graph
	_, err = io.WriteString(writer, "}\n")
	if err != nil {
		return err
	}

	//  Return nil if no error
	return nil
}
