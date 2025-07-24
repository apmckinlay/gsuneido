// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Getter_(mem)
		{
		return .getClass()[mem]
		}
	Default(@args)
		{
		return .getClass()[args[0]](@+1 args)
		}

	getClass()
		{
		switch
			{
		case Sys.Windows?(): 	return .windows
		case Sys.Linux?():		return LocalCmds_Linux
		case Sys.MacOS?():		return .mac
			}
		}

	windows: class
		{
		Netstat: 'netstat -aon'
		PidKillRegex: '(\d\d\d\d(\d)?)$'
		IpInfo: 'ipconfig /all'		//TODO: make this a method
		XCopy: 'xcopy /D /I /Y ' 	//TODO: change code to use this
		Nl: '\r\n'
		ScriptHeader: ''
		Unzip: 'unzip.exe'
		ScriptExt: '.bat'
		HttpTestThrowMatch: 'unknown address|could not resolve|no such host'

		GetHardDriveInfo()
			{
			return RunPipedOutput(
				'powershell get-physicaldisk | Select FriendlyName, MediaType').Trim()
			}
		Taskkill(pid)
			{
			return Spawn(P.WAIT, 'taskkill', '/pid', pid, '/f', '/t')
			}
		GetLastBootTime()
			{
			// cmd returns value like "20160126073017.491649-360" where
			// part before "." is date/time, part between "." and "-" is the
			// microseconds, and the final part is the offSet from UTC in mins
			bootTime = .wmiQuery('LastBootUpTime', 'win32_operatingsystem')
			if bootTime.Empty?()
				return ''
			ymd = `(\d\d\d\d\d\d\d\d)`
			hmss = `(\d\d\d\d\d\d)\.(\d\d\d)`
			formattedTime = Date(bootTime[0].Replace(ymd $ hmss $ '(.*)', '\1.\2\3'))
			return formattedTime is false
				? bootTime[0]
				: formattedTime
			}
		GetCpuName()
			{
			return .wmiQuery('Name')[0]
			}
		GetLogicalProcessors()
			{
			logicProcessors =  .wmiQuery('NumberOfLogicalProcessors')
			return logicProcessors.Filter({ not it.Has?("NumberOfLogicalProcessors") }).
					SumWith(Number)
			}
		GetCpuCores()
			{
			cores = enabled = logic = 0
			queryResults = ''
			try
				{
				queryResults = .wmiQuery('NumberOfCores')

				cores = queryResults.Filter({ not it.Has?('NumberOfCores') }).
					SumWith(Number)
				}
			catch (err)
				throw err $ ' - ' $ Display(queryResults)
			try
				enabled = .wmiQuery('NumberOfEnabledCore').
					Filter({ not it.Has?('NumberOfEnabledCore') }).SumWith(Number)
			catch(err)
				SuneidoLog("INFO: cannot query WMI to get NumberOfEnabledCore - " $ err)

			cores = Max(enabled, cores)
			try
				logic = .GetLogicalProcessors()
			catch (err)
				SuneidoLog("INFO: cannot get logical processors - " $ err)
			// cpu cores should not be more than logic processors

			return logic is 0 ? cores : Min(cores, logic)
			}

		GetNetworkSpeed()
			{
			return .wmiQuery('Name,Speed',
				'win32_networkadapter -Filter NetEnabled=true', false).Each(#RightTrim).
					Delete(1).Join('\r\r\n')
			}
		Robocopy(currentFolder, newFolder)
			{
			exitCode = Spawn(P.WAIT, 'robocopy', '/E', '/R:10', currentFolder, newFolder)
			return exitCode < 8 /*= min code */
			}

		wmiQuery(property, propertyClass = 'Win32_Processor', expandProperty? = true)
			{
			propertyType = expandProperty? ? '-ExpandProperty ' : '-Property '
			x = RunPipedOutput(PowerShell() $ ` -Command "Get-WmiObject -class ` $
				propertyClass $ ` -Property ` $ property $
				` | Select-Object ` $ propertyType $ property $ `"`).Trim('\n\r ')
			return x.Split('\r\n')
			}
		}

	mac: LocalCmds_Linux
		{
		GetCpuName()
			{
			return RunPipedOutput('sysctl -n machdep.cpu.brand_string').Trim()
			}
		GetLogicalProcessors()
			{
			return RunPipedOutput('sysctl -n hw.logicalcpu').Trim()
			}
		GetCpuCores()
			{
			return Number(RunPipedOutput('sysctl -n hw.physicalcpu').Trim())
			}
		}
	}
