// Copyright 2021 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
)

// Run runs the HTTP server until ctx is canceled.
// The done channel has an error written to when the HTTP server
// is terminated, and can be nil or not nil.
func (s *Server) Run(ctx context.Context, ready chan<- struct{}, done chan<- error) {
	server := http.Server{Addr: s.address, Handler: s.handler}

	crashed := make(chan struct{})
	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		select {
		case <-ctx.Done():
		case <-crashed:
			return
		}

		s.logger.Warn(s.name + " http server shutting down: " + ctx.Err().Error())
		shutdownCtx, cancel := context.WithTimeout(context.Background(),
			s.optional.shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error(s.name + " http server failed shutting down within " +
				s.optional.shutdownTimeout.String())
		}
	}()

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		close(s.addressSet)
		close(crashed) // stop shutdown goroutine
		<-shutdownDone
		done <- err
		return
	}

	s.address = listener.Addr().String()
	close(s.addressSet)

	// note: no further write so no need to mutex
	s.logger.Info(s.name + " http server listening on " + s.address)
	close(ready)

	err = server.Serve(listener)

	if err != nil && !errors.Is(ctx.Err(), context.Canceled) {
		// server crashed
		close(crashed) // stop shutdown goroutine
	} else {
		err = nil
	}
	<-shutdownDone
	done <- err
}
