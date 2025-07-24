// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	MsgPrefix: 'Server Stats:'
	CallClass()
		{
		return .format(.AsObject())
		}

	AsString(statsOb)
		{
		return .format(statsOb)
		}
	format(stats)
		{
		msg = .MsgPrefix $ "\n"
		msg $= "OS: " $ stats.system $ "\n"
		msg $= "Processor: " $ stats.processor $ "\n"
		msg $= "Number of Logical Processors: " $ stats.num_logical_processors $ '\n'
		msg $= "Number of CPU Cores: " $ stats.num_cpu_cores $ '\n'
		msg $= 'Running since ' $ stats.start_time.StdShortDateTime() $
			' - ' $ .uptime(stats.start_time) $ '\n'
		msg $= 'Last Boot Time: ' $ Display(stats.last_boot_time) $ '\n'
		msg $= Opt("port ", ServerPort(), "\n")
		msg $= 'MAC Address: ' $ stats.mac_address $ "\n"
		msg $= "Disk Free Space: " $ ReadableSize(stats.disk_free_space) $ "\n"
		msg $= "database size: " $ ReadableSize(stats.dbsize) $ "\n"
		msg $= "heap size: " $ ReadableSize(stats.memory) $ "\n"
		msg $= "transactions: " $ stats.transactions $ "\n"
		msg $= "cursors: " $ stats.cursors $ "\n"
		msg $= "connections: (" $ stats.connections.Size() $ ") " $
				stats.connections.Join(', ') $ '\n'
		if stats.Member?(#threads)
			msg $= "threads: (" $ stats.threads.Size() $ ") " $
				stats.threads.Join(', ') $ '\n'
		msg $= "libraries in use: " $ stats.libraries.Join(", ") $ '\n'
		msg $= "All MAC Addresses: " $ stats.all_mac_addresses.Join(", ") $ '\n'
		msg $= "Min Negotiated Network Speed (MB): " $ stats.nic_info.minSpeed $ '\n'
		msg $= stats.nic_info.allInfo $ '\n'
		msg $= "OS Build: " $ stats.build $ '\n'
		msg $= "Bios Version: " $ stats.bios_version $ '\n'
		for c in Contributions('ServerStatsExtra')
			{
			stat = c()
			msg $= stat.prompt $ ': ' $ stats[stat.name] $ '\n'
			}
		msg $= "Disk Information:\n" $ stats.hd_info $ '\n\n'
		return msg
		}

	AsObject()
		{
		sys = .GetSystemInfo()
		stats = Object(start_time: Suneido.start_time,
			system: SystemSummary(),
			last_boot_time: .GetLastBootTime(),
			processor: .GetCpuName(),
			num_logical_processors: .getNumLogicalProcs(),
			num_cpu_cores: .getCpuCores(),
			mac_address: GetMainMacAddressHex(),
			all_mac_addresses: GetMacAddressesHex(),
			memory: MemoryArena(),
			dbsize: Database.CurrentSize(),
			transactions: Database.Transactions().Size(),
			cursors: Database.Cursors(),
			tempdest: Database.TempDest(),
			connections: Sys.Connections(),
			libraries: Libraries(),
			disk_free_space: GetDiskFreeSpace(),
			nic_info: .GetNetworkSpeed(),
			bios_version: sys.bios,
			build: sys.build,
			hd_info: .getHardDriveInfo()
			threads: Thread.List().
				RemoveIf({ it =~ "-connection-|-thread-pool" }).Sort!())
		for c in Contributions('ServerStatsExtra')
			{
			stat = c()
			stats[stat.name] = stat.value
			}
		return stats
		}

	getHardDriveInfo()
		{
		try
			return LocalCmds.GetHardDriveInfo()
		catch (e)
			return e.Tr('\0')
		}

	GetLastBootTime()
		{
		try
			return LocalCmds.GetLastBootTime()
		catch (e)
			return e
		}
	GetCpuName()
		{
		try
			return LocalCmds.GetCpuName()
		catch (e)
			return e
		}
	getNumLogicalProcs()
		{
		try
			return LocalCmds.GetLogicalProcessors()
		catch (e)
			return e
		}
	getCpuCores()
		{
		try
			return LocalCmds.GetCpuCores()
		catch (e)
			return e
		}
	GetNetworkSpeed()
		{
		if Sys.Client?()
			return ServerEval('ServerStats.GetNetworkSpeed')

		minSpeed = 'NOT FOUND'
		try
			{
			allInfoOb = LocalCmds.GetNetworkSpeed().Lines().RemoveIf(#Blank?)
			allInfoOb[0] = 'NIC Info:  ' $ allInfoOb[0]
			for mem in allInfoOb[1..].Members()
				{
				allInfoOb[mem+1] = 'NIC Info' $ (mem+1) $ ': '  $ allInfoOb[mem+1].Trim()
				speed = allInfoOb[mem+1].Extract("\d+$")
				if speed isnt false and Number(speed) < minSpeed
					minSpeed = Number(speed) / 1000000 /*=megs*/
				}
			allInfo = allInfoOb.Join('\n')
			}
		catch (e)
			allInfo = 'NIC Info: ERROR: ' $ e
		if not Number?(minSpeed)
			minSpeed = 0
		return Object(:minSpeed, :allInfo)
		}

	GetSystemInfo()
		{
		bios = SystemInfo().Bios
		build = SystemInfo().Build
		osname = SystemInfo().OSName
		osversion = SystemInfo().OSVersion
		if Sys.Linux?()
			{
			// removed unnecessary command that requires sudo permission
			try
				build = RunPipedOutput('lsb_release -d').AfterFirst(':').Trim()
			catch (err1)
				try
					build = RunPipedOutput('cat /proc/version').Trim()
				catch (err2)
					SuneidoLog('ERROR: failed to get linux system info - ' $
						err1 $ ' - '$ err2)
			}
		if Sys.MacOS?()
			{
			try
				{
				osname $= ' ' $ RunPipedOutput('sw_vers -productVersion').Trim()
				build = RunPipedOutput('sw_vers -buildVersion')
				// macOS system does not have bios
				}
			catch (err)
				{
				SuneidoLog('INFO: failed to get macOS system info - ' $ err)
				}
			}
		return [:bios, :build, :osname, :osversion] // build like "10.0.17384"
		}

	Ipconfig()
		{
		if Sys.Client?()
			return ServerEval('ServerStats.Ipconfig')

		str = RunPipedOutput(LocalCmds.IpInfo).Replace('^\r\n|:|^   ', '')
		str $= .ExternalIp()
		str $= Opt('\r\nInstance Id: ', .instanceId())
		return str.Trim()
		}

	external_IP_URL: 'http://checkip.amazonaws.com/'
	IpProblemMessage: 'Problem getting external IP address'
	ExternalIp(noPrompt = false)
		{
		external_ip = .IpProblemMessage
		try external_ip = Http.Get(.external_IP_URL).Tr('\r\n')
		if not external_ip.Tr('.').Numeric?()
			external_ip = .IpProblemMessage
		if noPrompt
			return external_ip
		return 'External Server IP: \t' $ external_ip
		}

	instanceId()
		{
		return OptContribution('InstanceId', function () { return '' })()
		}

	uptime(start_time, asof = 'now')
		{
		if asof is 'now'
			asof = Date()
		days = asof.MinusDays(start_time)
		if days isnt 0
			return days $ " day(s)"

		seconds = asof.MinusSeconds(start_time)
		secondsPerHour = 3600
		hours = (seconds / secondsPerHour).Round(0)
		return hours $ " hour(s)"
		}
	}