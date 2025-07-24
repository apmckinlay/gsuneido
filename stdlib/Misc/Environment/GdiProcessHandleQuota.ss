// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Singleton
	{
	New()
		{
		.readQuota()
		}

	reg: "'HKLM:\HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Windows'"
	quota: 10000
	readQuota()
		{
		result = false
		try
			{
			result = RunPipedOutput.WithExitValue(
				PowerShell() $ " Get-ItemPropertyValue -Path " $
				.reg $ " -Name USERProcessHandleQuota")
			if result.exitValue is 0
				{
				value = Number(result.output.BeforeLast('\r\n'))
				// https://docs.microsoft.com/en-us/windows/win32/sysinfo/gdi-objects: value should be between 256 and 65536
				if value < 256 or value > 65536  /*= based on msdn above */
					throw "Invalid value for USERProcessHandleQuota: " $ result.output
				.quota = value
				}
			}
		catch (err)
			{
			SuneidoLog('ERRATIC: Cannot get gdi handle quota, ' $ 'checking skipped - ' $
				err, params: result)
			}
		}

	gdiUsageFactor: 0.9
	CheckOverLimit()
		{
		usage = Suneido.Info("windows.nGdiObject")
		if .quota * .gdiUsageFactor <= usage
			return 'windows GDI resource over threshold: ' $
				'threshold is ' $  .quota * .gdiUsageFactor $
				' (limit: ' $ .quota $ ', factor: ' $ .gdiUsageFactor $ ')' $
				' and the actual usage is ' $ usage

		return ''
		}
	}
