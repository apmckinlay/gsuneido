// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Singleton
	{
	OSName: ''
	Build: ''
	OSVersion: #(0, 0, 0)
	Bios: ''
	New()
		{
		if Sys.Linux?() or Sys.MacOS?()
			{
			if Sys.Linux?()
				.Bios = .biosInfo()
			.OSName = OSName()
			return
			}

		try
			s = RunPipedOutput(.Script())
		catch (err)
			{
			.logErr('INFO: failed to get system info - ' $ err)
			return
			}
		.parseInfo(s)
		}

	parseInfo(s)
		{
		info = Object().Set_default('')
		for l in s.Trim().Lines()
			{
			m = l.BeforeFirst(':').Trim()
			if m isnt ''
				info[m] = l.AfterFirst(':').Trim()
			}
		.OSName = info.Caption.Replace('Microsoft', '').Trim()
		.Build = info.Version
		try
			.OSVersion = .Build.Split('.').Map!(Number)
		catch (err)
			{
			.logErr('ERROR: unexpected system version number - ' $ .Build $ ' - ' $ err)
			}
		if .OSName is '' or .OSVersion.Size() < 3 /*= major, minor, build*/
			.logErr('INFO: Missing operating system info - ' $ s)
		if .OSVersion.Size() < 3 /*= major, minor, build*/
			.OSVersion = .OSVersion.Copy().MergeNew(#(0, 0, 0))
		.Bios = info.Manufacturer $ ' ' $ info.SMBIOSBIOSVersion $ ', ' $ info.ReleaseDate
		}

	biosInfo()
		{
		try
			return RunPipedOutput('cat /sys/class/dmi/id/bios_vendor ' $
				'/sys/class/dmi/id/bios_version ' $
				'/sys/class/dmi/id/bios_date ' $
				'/sys/class/dmi/id/product_name').Tr('\r\n', ' ').Trim()
		catch
			return 'not implemented'
		}

	Script()
		{
		return PowerShell() $ ` "Get-CIMInstance Win32_OperatingSystem` $
			` | select -Property Caption,Version | Format-List; ` $
			`Get-CIMInstance -class Win32_bios | ` $
			`select -Property Manufacturer,SMBIOSBIOSVersion,ReleaseDate ` $
			`| Format-List"`
		}

	ShowScript()
		{
		PutFile('axon_system_info.bat', .Script())
		System('start /w axon_system_info.bat')
		DeleteFile('axon_system_info.bat')
		}

	logErr(err)
		{
		if not TestRunner.RunningTests?()
			ErrorLog(err)
		}
	}
