// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		mask = '#####.####'
		Assert(HoursControl.HoursControl_convertToHours('', mask) is: '')
		Assert(HoursControl.HoursControl_convertToHours(1, mask) is: '1:00')
		Assert(HoursControl.HoursControl_convertToHours(1.5, mask) is: '1:30')

		Assert(HoursControl.HoursControl_convertToHours(1.01666666666, mask) is: '1:01')
		Assert(HoursControl.HoursControl_convertToHours(1.02, mask) is: '1:01')
		Assert(HoursControl.HoursControl_convertToHours(1.03, mask) is: '1:02')
		Assert(HoursControl.HoursControl_convertToHours(1.04, mask) is: '1:02')
		Assert(HoursControl.HoursControl_convertToHours(1.05, mask) is: '1:03')
		Assert(HoursControl.HoursControl_convertToHours(1.06, mask) is: '1:04')
		Assert(HoursControl.HoursControl_convertToHours(1.07, mask) is: '1:04')
		Assert(HoursControl.HoursControl_convertToHours(1.08, mask) is: '1:05')
		Assert(HoursControl.HoursControl_convertToHours(1.09, mask) is: '1:05')
		Assert(HoursControl.HoursControl_convertToHours(1.1, mask) is: '1:06')

		Assert(HoursControl.HoursControl_convertToHours(1.5, mask) is: '1:30')
		Assert(HoursControl.HoursControl_convertToHours(2.75, mask) is: '2:45')
		Assert(HoursControl.HoursControl_convertToHours(3.6, mask) is: '3:36')
		Assert(HoursControl.HoursControl_convertToHours('3.6', mask) is: '3:36')
		Assert(HoursControl.HoursControl_convertToHours('3:36', mask) is: '3:36')
		Assert(HoursControl.HoursControl_convertToDecimal('', mask) is: '')
		Assert(HoursControl.HoursControl_convertToDecimal(1, mask) is: 1)
		Assert(HoursControl.HoursControl_convertToDecimal('1:00', mask) is: 1)
		Assert(HoursControl.HoursControl_convertToDecimal('1:30', mask) is: 1.5)
		Assert(HoursControl.HoursControl_convertToDecimal('2:45', mask) is: 2.75)
		Assert(HoursControl.HoursControl_convertToDecimal('3:36', mask) is: 3.6)
		Assert(HoursControl.HoursControl_convertToDecimal('5.:00', mask) is: 5)
		Assert(HoursControl.HoursControl_convertToDecimal('5:.6', mask) is: 5.01)
		Assert(HoursControl.HoursControl_convertToDecimal('5.:.6', mask) is: 5.01)
		Assert(HoursControl.HoursControl_convertToDecimal('5.5:.6', mask) is: 5.51)
		Assert(HoursControl.HoursControl_convertToDecimal('5:300.6', mask) is: 10.01)

		for (i = 0; i < 60; i++)
			{
			orig = '1:' $ i.Pad(2, '0')
			decimal = HoursControl.HoursControl_convertToDecimal(orig, mask)
			hours = HoursControl.HoursControl_convertToHours(decimal, mask)
			Assert(hours is: orig
				msg: 'Converted does not match (orig: ' $ orig $ ', converted: ' $ hours)
			}

		Assert(HoursControl.HoursControl_convertToHours('1:30', mask) is: '1:30')
		Assert(HoursControl.HoursControl_convertToDecimal('1.5', mask) is: 1.5)
		Assert(HoursControl.HoursControl_convertToDecimal('1:', mask) is: '')
		Assert(HoursControl.HoursControl_convertToDecimal(':45', mask) is: .75)
		// 75min
		Assert(HoursControl.HoursControl_convertToDecimal(':75', mask) is: 1.25)
		// 3 hrs, 75 mins
		Assert(HoursControl.HoursControl_convertToDecimal('3:75', mask) is: 4.25)

		// test bad values, should all just get cleared out
		Assert(HoursControl.HoursControl_convertToDecimal('5..6', mask) is: '')
		Assert(HoursControl.HoursControl_convertToHours('fred', mask) is: '')
		Assert(HoursControl.HoursControl_convertToDecimal('fred', mask) is: '')
		Assert(HoursControl.HoursControl_convertToHours('fred:30', mask) is: '')
		Assert(HoursControl.HoursControl_convertToHours('1:fred', mask) is: '')
		// hours too big for mask
		Assert(HoursControl.HoursControl_convertToHours('123456:45', mask) is: '')
		Assert(HoursControl.HoursControl_convertToDecimal('99999:60', mask) is: '')
		Assert(HoursControl.HoursControl_convertToDecimal(':9999999', mask) is: '')
		}
	}