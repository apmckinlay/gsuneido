// Copyright (C) 2015 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass()
		{
		Suneido.RunningHttpServer = true
		RackServer(app: RackRouter(.Routes()),
			with: [RackResponseHeaders, RackContentType],
			port: HttpPort())
		}

	Routes()
		{
		return OptContribution('RackRoutesOverride', GetContributions('RackRoutes'))
		}
	}
