// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_OtherIP()
		{
		spy = .SpyOn(NetworkService.IPAddrs)
		spy.Return(Object(),
			Object('ip1', 'ip2'),

			Object('ip1', 'ip2'),
			Object('ip1', 'ip2')

			Object('ip1'),

			Object('ip1', 'ip2', 'ip3'),
			Object('ip1', 'ip2', 'ip3'),
			Object('ip1', 'ip2', 'ip3')
			)

		// empty list
		Assert(NetworkService.OtherIP('ipx') is: 'ipx')

		// two ips list
		Assert(NetworkService.OtherIP('ip1') is: 'ip2')
		Assert(NetworkService.OtherIP('ip2') is: 'ip1')

		// not found
		Assert(NetworkService.OtherIP('ipx') is: 'ip1')

		// not found and the list size is 1
		Assert(NetworkService.OtherIP('ipx') is: 'ip1')

		// three ips list
		Assert(NetworkService.OtherIP('ip1') is: 'ip2')
		Assert(NetworkService.OtherIP('ip2') is: 'ip3')
		Assert(NetworkService.OtherIP('ip3') is: 'ip1')
		}

	Test_IPAddress()
		{
		.SpyOn(ServerSuneido.Get).Return(false, 'ip1',
			Object(ip: 'ip2', expires: Date().Plus(days: -1)),
			Object(ip: 'ip3', expires: Date().Plus(days: 1)))
		.SpyOn(NetworkService.IPAddrs).Return(#('ip4'))

		Assert(NetworkService.IPAddress() is: 'ip4')
		Assert(NetworkService.IPAddress() is: 'ip1')
		Assert(NetworkService.IPAddress() is: 'ip4')
		Assert(NetworkService.IPAddress() is: 'ip3')
		}
	}