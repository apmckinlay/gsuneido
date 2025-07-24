// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_FileRegx()
		{
		Assert("suneido:/ETAHelp/res/11.png" =~ SuneidoAPP.FileRegx)
		Assert("suneido:/ETAHelp/res/11.png.test" =~ SuneidoAPP.FileRegx)
		Assert("suneido:/ETAHelp/res/11 test.png.test" =~ SuneidoAPP.FileRegx)
		Assert("suneido:/ETAHelp/res/folder/w1.png" =~ SuneidoAPP.FileRegx)
		Assert("suneido:/ETAHelp/res/folder two/w1.png" =~ SuneidoAPP.FileRegx)

		Assert("suneido:/ETAHelp/11.jpg" !~ SuneidoAPP.FileRegx)
		Assert("suneido:/11.jpg" !~ SuneidoAPP.FileRegx)
		Assert("suneido:/ETAHelp/res/11 test.png." !~ SuneidoAPP.FileRegx)
		Assert("suneido:/ETAHelp/res/11 test.png.test two" !~ SuneidoAPP.FileRegx)

		Assert("suneido:/ETAHelp/Folder/Item" !~ SuneidoAPP.FileRegx)
		Assert("suneido:/ETAHelp/Folder/Item Two" !~ SuneidoAPP.FileRegx)
		Assert("suneido:/ETAHelp/Folder Two/Item Two" !~ SuneidoAPP.FileRegx)
		Assert("suneido:/ETAHelp/Folder Two/Sub Folder/Item Two" !~ SuneidoAPP.FileRegx)
		}

	Test_browserHwnd()
		{
		// Empty browserLoads object
		Assert(.simulateHwnd('url_0', #()) is: false)

		browserLoads = Object(
			1000: browser1 = Object(
				url_0: #20250101.010001,
				url_1: #20250101.010002)
			2000: browser2 = Object(
				url_0: #20250101.010002)
			3000: browser3 = Object(
				url_0: #20250101.010003,
				url_1: #20250101.010001,
				url_2: #20250101.010001))

		// URL does not exist in browserLoads
		Assert(.simulateHwnd('url_3', browserLoads) is: false)

		Assert(browser1 members: #(url_0, url_1))
		Assert(.simulateHwnd('url_0', browserLoads) is: 1000)
		Assert(browser1 members: #(url_1))

		Assert(browser2 members: #(url_0))
		Assert(.simulateHwnd('url_0', browserLoads) is: 2000)
		Assert(browser2 members: #())

		Assert(browser3 members: #(url_0, url_1, url_2))
		Assert(.simulateHwnd('url_0', browserLoads) is: 3000)
		Assert(browser3 members: #(url_1, url_2))

		Assert(browser3 members: #(url_1, url_2))
		Assert(.simulateHwnd('url_1', browserLoads) is: 3000)
		Assert(browser3 members: #(url_2))

		Assert(browser1 members: #(url_1))
		Assert(.simulateHwnd('url_1', browserLoads) is: 1000)
		Assert(browser1 members: #())

		Assert(browser3 members: #(url_2))
		Assert(.simulateHwnd('url_2', browserLoads) is: 3000)
		Assert(browser3 members: #())
		}

	simulateHwnd(url, browserLoads)
		{
		if false isnt hwnd = SuneidoAPP.SuneidoAPP_browserHwnd(url, browserLoads)
			browserLoads[hwnd].Delete(url)
		return hwnd
		}
	}