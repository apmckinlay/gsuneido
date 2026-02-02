// Copyright (C) 2015 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass()
		{
		Suneido.RunningHttpServer = true
		with = [RackResponseHeaders, RackContentType].MergeUnion(
			GetContributions('RackWithExtra'))
		RackServer(app: RackRouter(.Routes()), :with, port: HttpPort())
		}

	Routes()
		{
		return OptContribution('RackRoutesOverride', GetContributions('RackRoutes'))
		}
	}
