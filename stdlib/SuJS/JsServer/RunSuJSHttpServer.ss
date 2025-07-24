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
		RackServer(app: RackRouter(SuJsRackRoutes().Append(extraRoutes)),
			with: [RackResponseHeaders, RackContentType],
			:port, name: 'SuJS Server')
		}
	}