// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(port = false, extraRoutes = #())
		{
		if port is false
			port = HttpPort() + 1

		// Sys.Server?() is whether the current "system" is acting as a server
		// A Suneido.js server is acting as a server even if the exe is in standalone
		Sys.SetServer()
		LibraryTags.AddMode('webgui')

		// Preload the class to prevent duplicate token issues caused by concurrent
		// .CreateToken calls. If the class code isn't already loaded and multiple calls
		// occur at the same time, .Synchronized cannot enforce execution order reliably.
		// This is because synchronization applies only within a single internal Suneido
		// class, while multiple internal classes may be loaded and returned concurrently.
		Global('JsSessionToken')

		RackServer(app: RackRouter(SuJsRackRoutes().Append(extraRoutes)),
			with: [RackResponseHeaders, RackContentType],
			:port)
		}
	}