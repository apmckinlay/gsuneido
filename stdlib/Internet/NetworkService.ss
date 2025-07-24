// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
class
	{
	IPAddrs()
		{
		return GetContributions('NetworkServiceIPAddrs')
		}
	IPAddress()
		{
		ips = .IPAddrs()
		if ips.Empty?()
			throw 'No network service IP addresses have been defined'
		if false isnt value = ServerSuneido.Get('NetworkServiceIPAddress', false)
			{
			if not Object?(value)
				return value
			if value.GetDefault('expires', Date.Begin()) >= Date()
				return value.ip
			}
		return ips[Random(ips.Size())]
		}

	SetIPAddress(value)
		{
		ServerSuneido.Set('NetworkServiceIPAddress',
			Object(ip: value, expires: .expires()))
		}
	baseInterval: 30
	expires()
		{
		return Date().Plus(minutes: .baseInterval + Random(.baseInterval))
		}
	OtherIP(ip)
		{
		ips = .IPAddrs()
		if ips.Empty?()
			return ip
		index = ips.Find(ip)
		next = index is false or index is ips.Size() - 1 ? 0 : ++index
		return ips[next]
		}
	RegisteredForService?()
		{
		func = OptContribution('RegisteredForNetworkService', function(){ return true })
		return func()
		}
	NotRegisteredMessage: 'Not Registered For Network Services'

	// expects true or false from block
	// returns false if CircuitBreaker is open
	RunWithService(block)
		{
		return CircuitBreaker('NetworkService', .runWithCircuitBreaker, block)
		}
	runWithCircuitBreaker(block)
		{
		ip = NetworkService.IPAddress()
		if false is block(ip) and
			false is block(ip = NetworkService.OtherIP(ip))
				return false
		NetworkService.SetIPAddress(ip)
		return true
		}
	LogMsg(ip1 = "", ip2 = "")
		{
		return Opt('unable to connect to ', ip1) $ Opt(', switch to ', ip2)
		}
	}
