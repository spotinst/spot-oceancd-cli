package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// This code was copied with adjustments from
// https://github.com/argoproj/argo-rollouts/blob/a6dbe0ec2db3f02cf695ba3c972db72cecabaefb/pkg/kubectl-argo-rollouts/cmd/signals/signal.go#L10
func SetupSignalHandler(cancel context.CancelFunc) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()
}
