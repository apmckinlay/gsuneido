class
	{
	Netstat: 'netstat -apn'
	PidKillRegex: '(\d\d\d\d(\d)?)\/'
	IpInfo: 'ip addr show'
	XCopy: 'cp -p -u '
	Nl: '\n'
	ScriptHeader: '#!/bin/bash'
	Unzip: 'unzip'
	ScriptExt: ''
	HttpTestThrowMatch: 'socket open failed'

	GetHardDriveInfo()
		{
		return RunPipedOutput('df -h').Trim()
		}
	Taskkill(pid)
		{
		return Spawn(P.WAIT, 'kill', '-9', pid)
		}
	GetLastBootTime()
		{
		s = ''
		try
			{
			if false isnt date = Date(RunPipedOutput('uptime -s').Trim())
				return date
			s = RunPipedOutput('uptime').Trim()
			return .parseUpTime(s)

			}
		catch (err)
			{
			if err.Has?('executable file not found') or
				err.Has?('No such file or directory')
				return false
			SuneidoLog('ERROR: LocalCmds.linux.GetLastBootTime ' $ err $ ' - ' $ s)
			return false
			}
		}

	parseUpTime(s)
		{
		s1 = s.BeforeFirst(' up')
		hours = Number(s1.BeforeFirst(':'))
		minutes = Number(s1.AfterFirst(':').BeforeFirst(':'))
		days = upHours = upMins = 0
		if s.Has?(' days')
			{
			days = Number(s.AfterFirst('up ').BeforeFirst(' days'))
			s2 = s.AfterFirst(',').BeforeFirst(',').Trim()
			upHours = Number(s2.BeforeFirst(':'))
			upMins = Number(s2.AfterFirst(':'))
			}
		else if s.Has?(' min')
			{
			days = 0
			upHours = 0
			upMins = Number(s.AfterFirst('up ').BeforeFirst(' min'))
			}
		else
			{
			s2 = s.AfterFirst('up ').BeforeFirst(',')
			upHours = Number(s2.BeforeFirst(':'))
			upMins = Number(s2.AfterFirst(':'))
			}
		return Date().NoTime().Plus(:hours, :minutes).
			Minus(:days, hours: upHours, minutes: upMins)
		}

	arm: #('CPU implementer', 'CPU architecture', 'CPU variant')
	GetCpuName()
		{
		s = .cpuinfo()
		if s.Has?('model name')
			return s.Extract('model name\s\:\s(\w+).+', 0).Trim('model name\t:')
		if s.Has?('CPU implementer')
			{
			// nightly checks file checksum relies on "0x41" as arm architecture
			cpu = .arm.Map({ s.Extract(it $ '\s*\:\s*(\w+)') }).Remove(false).Join(',')
			if false isnt v = s.Extract('BogoMIPS\s*\:\s*(.+)')
				cpu $= ' @ ' $ v $ ' BogoMIPS'
			return cpu
			}
		throw "unhandled cpuinfo"
		}
	cpuinfo()
		{
		return RunPipedOutput('cat /proc/cpuinfo')
		}
	GetLogicalProcessors()
		{
		return RunPipedOutput('nproc').Tr('\n')
		}
	GetCpuCores()
		{
		if false isnt cores = .cpuinfo().Extract('cpu cores\s\:\s.+', 0)
			return Number(cores.Trim('cpu cores\t:'))
		return 'not available'
		}
	GetNetworkSpeed()
		{
		return 'not implemented'
		}
	Robocopy(currentFolder, newFolder)
		{
		return 0 is Spawn(P.WAIT, 'cp', '-r', '-T', currentFolder, newFolder)
		}
	}
