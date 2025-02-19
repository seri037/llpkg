package config

import (
	"fmt"
	"time"
)

type LoadingSpinner struct {
	stopChan chan struct{}
	message  string
}

func NewLoadingSpinner(message string) *LoadingSpinner {
	return &LoadingSpinner{
		stopChan: make(chan struct{}),
		message:  message,
	}
}

func (s *LoadingSpinner) Start() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	go func() {
		for i := 0; ; i++ {
			select {
			case <-s.stopChan:
				fmt.Printf("\r%s... Done!    \n", s.message)
				return
			default:
				frame := frames[i%len(frames)]
				fmt.Printf("\r%s %s ", frame, s.message)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (s *LoadingSpinner) Stop() {
	s.stopChan <- struct{}{}
}
