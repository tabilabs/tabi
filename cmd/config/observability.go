// Copyright 2024 Tabi Foundation
// This file is part of the Tabi Network packages.
//
// Tabi is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Tabi packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

package config

func EnableObservability() error {
	return nil
	// if true {
	// Temporarily disabling this until we can configure out port reuse
	// fast enough or enabling observability through the config.
	// return nil
	// }

	// pe, err := prometheus.NewExporter(prometheus.Options{
	// 	Namespace: "tabid",
	// })
	// if err != nil {
	// 	return fmt.Errorf("cmd/config: failed to create the OpenCensus Prometheus exporter: %w", err)
	// }

	// views := app.ObservabilityViews()
	// if err := view.Register(views...); err != nil {
	// 	return fmt.Errorf("cmd/config: failed to register OpenCensus views: %w", err)
	// }
	// view.RegisterExporter(pe)

	// mux := http.NewServeMux()
	// mux.Handle("/metrics", pe)

	// // TODO: Derive the Prometheus observability exporter from the Tabi config.
	// addr := ":8877"
	// go func() {
	// 	println("Serving the Prometheus observability exporter at", addr)
	// 	if err := http.ListenAndServe(addr, mux); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// return nil
}
