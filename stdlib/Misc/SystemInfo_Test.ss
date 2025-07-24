// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		example = ``
		info = Mock(SystemInfo)
		info.Eval(SystemInfo.SystemInfo_parseInfo, example)
		info.Verify.SystemInfo_logErr([anyArgs:])
		Assert(info.OSName is: '')
		Assert(info.Bios is: ' , ')
		Assert(info.OSVersion is: #(0, 0, 0))
		Assert(info.Build is: '')

		example = `

Caption : Microsoft Windows 10 Pro
Version : 10.0.19044





Manufacturer      : American Megatrends Inc.
SMBIOSBIOSVersion : F5
ReleaseDate       : 2019-08-12 6:00:00 PM


`
		info = Mock(SystemInfo)
		info.Eval(SystemInfo.SystemInfo_parseInfo, example)
		info.Verify.Never().SystemInfo_logErr([anyArgs:])
		Assert(info.OSName is: 'Windows 10 Pro')
		Assert(info.Bios is: 'American Megatrends Inc. F5, 2019-08-12 6:00:00 PM')
		Assert(info.OSVersion is: #(10, 0, 19044))
		Assert(info.Build is: '10.0.19044')

		example = `

Caption : Microsoft Windows 11 Pro
Version : 10.0.22000





Manufacturer      : Dell Inc.
SMBIOSBIOSVersion : 1.4.3
ReleaseDate       : 1/6/2022 6:00:00 PM


`
		info = Mock(SystemInfo)
		info.Eval(SystemInfo.SystemInfo_parseInfo, example)
		info.Verify.Never().SystemInfo_logErr([anyArgs:])
		Assert(info.OSName is: 'Windows 11 Pro')
		Assert(info.Bios is: 'Dell Inc. 1.4.3, 1/6/2022 6:00:00 PM')
		Assert(info.OSVersion is: #(10, 0, 22000))
		Assert(info.Build is: '10.0.22000')


		example = `

Caption : Microsoft Windows Server 2016 Standard
Version : 10.0.14393





Manufacturer      : American Megatrends Inc.
SMBIOSBIOSVersion : 3.0a
ReleaseDate       : 2/7/2018 6:00:00 PM


`
		info = Mock(SystemInfo)
		info.Eval(SystemInfo.SystemInfo_parseInfo, example)
		info.Verify.Never().SystemInfo_logErr([anyArgs:])
		Assert(info.OSName is: 'Windows Server 2016 Standard')
		Assert(info.Bios is: 'American Megatrends Inc. 3.0a, 2/7/2018 6:00:00 PM')
		Assert(info.OSVersion is: #(10, 0, 14393))
		Assert(info.Build is: '10.0.14393')

		example =`


Caption : Microsoft Windows Server 2016 Datacenter
Version : 10.0.14393





Manufacturer      : Amazon EC2
SMBIOSBIOSVersion : 1.0
ReleaseDate       : 10/15/2017 8:00:00 PM


`
		info = Mock(SystemInfo)
		info.Eval(SystemInfo.SystemInfo_parseInfo, example)
		info.Verify.Never().SystemInfo_logErr([anyArgs:])
		Assert(info.OSName is: 'Windows Server 2016 Datacenter')
		Assert(info.Bios is: 'Amazon EC2 1.0, 10/15/2017 8:00:00 PM')
		Assert(info.OSVersion is: #(10, 0, 14393))
		Assert(info.Build is: '10.0.14393')

		example =`
Caption : Microsoft Windows Server 2016 Datacenter
Manufacturer      : Amazon EC2
SMBIOSBIOSVersion : 1.0
ReleaseDate       : 10/15/2017 8:00:00 PM
`
		info = Mock(SystemInfo)
		info.Eval(SystemInfo.SystemInfo_parseInfo, example)
		Assert(info.OSName is: 'Windows Server 2016 Datacenter')
		Assert(info.Bios is: 'Amazon EC2 1.0, 10/15/2017 8:00:00 PM')
		Assert(info.OSVersion is: #(0, 0, 0))
		Assert(info.Build is: '')

		example =`
Caption : Microsoft Windows Server 2016 Datacenter
Version : 10.0.cccc
Manufacturer      : Amazon EC2
SMBIOSBIOSVersion : 1.0
ReleaseDate       : 10/15/2017 8:00:00 PM
`
		info = Mock(SystemInfo)
		info.Eval(SystemInfo.SystemInfo_parseInfo, example)
		info.Verify.Times(1).SystemInfo_logErr([anyArgs:])
		Assert(info.OSName is: 'Windows Server 2016 Datacenter')
		Assert(info.Bios is: 'Amazon EC2 1.0, 10/15/2017 8:00:00 PM')
		Assert(info.Build is: '10.0.cccc')
		Assert(info.OSVersion is: #(0, 0, 0))
		}
	}
