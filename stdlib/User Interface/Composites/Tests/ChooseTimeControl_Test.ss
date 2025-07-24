// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ConvertToMilitary()
		{
		ob = Object(minute: 12, suffix: 'am', hour: 2)
		Assert(ChooseTimeControl.ConvertToMilitary(ob) is: 212)

		ob.hour = 12
		Assert(ChooseTimeControl.ConvertToMilitary(ob) is: 12)

		ob.suffix = 'pm'
		Assert(ChooseTimeControl.ConvertToMilitary(ob) is: 1212)

		ob.hour = 2
		Assert(ChooseTimeControl.ConvertToMilitary(ob) is: 1412)
		}

	Test_ValidData?()
		{
		args = Object('')
		Assert(ChooseTimeControl.ValidData?(@args))

		args.mandatory = true
		Assert(ChooseTimeControl.ValidData?(@args) is: false)

		args[0] = 'false'
		Assert(ChooseTimeControl.ValidData?(@args) is: false)

		args[0] = '0000'
		Assert(ChooseTimeControl.ValidData?(@args))

		args[0] = '1:15 AM'
		Assert(ChooseTimeControl.ValidData?(@args) is: false)

		args[0] = '2616'
		Assert(ChooseTimeControl.ValidData?(@args) is: false)

		args[0] = '26:16'
		Assert(ChooseTimeControl.ValidData?(@args) is: false)

		args[0] = '1524'
		Assert(ChooseTimeControl.ValidData?(@args))

		args[0] = '0835'
		Assert(ChooseTimeControl.ValidData?(@args))

		args[0] = '835'
		Assert(ChooseTimeControl.ValidData?(@args))

		args[0] = 835
		Assert(ChooseTimeControl.ValidData?(@args))

		args[0] = '15:24'
		Assert(ChooseTimeControl.ValidData?(@args) is: false)

		args[0] = '15:bob'
		Assert(ChooseTimeControl.ValidData?(@args) is: false)
		}

	Test_SplitTime()
		{
		splitTime = ChooseTimeControl.SplitTime
		Assert(splitTime('') is: false)
		Assert(splitTime('1525') is: #(hour: 15, minute: 25, suffix: ""))
		Assert(splitTime('325') is: #(hour: 3, minute: 25, suffix: ""))
		Assert(splitTime('0325') is: #(hour: 3, minute: 25, suffix: ""))
		Assert(splitTime('1525 pm') is: false)
		Assert(splitTime('25') is: #(hour: 0, minute: 25, suffix: ""))
		Assert(splitTime('15:25') is: #(hour: 15, minute: 25, suffix: ""))
		Assert(splitTime('15:25 pm') is: #(hour: 15, minute: 25, suffix: "pm"))
		Assert(splitTime('3:25 pm') is: #(hour: 3, minute: 25, suffix: "pm"))
		Assert(splitTime('3:25 am') is: #(hour: 3, minute: 25, suffix: "am"))

		Assert(splitTime('00:01 pm') is: #(hour: 0, minute: 1, suffix: "pm"))

		Assert(splitTime('1215') is: #(hour: 12, minute: 15, suffix: ""))
		Assert(splitTime('0015') is: #(hour: 0, minute: 15, suffix: ""))

		Assert(splitTime('12:15') is: #(hour: 12, minute: 15, suffix: ""))
		Assert(splitTime('12:15 am') is: #(hour: 12, minute: 15, suffix: "am"))
		Assert(splitTime('12:15 pm') is: #(hour: 12, minute: 15, suffix: "pm"))

		Assert(splitTime('Fred') is: false)
		Assert(splitTime('Fr15') is: false)
		Assert(splitTime('15ed') is: false)
		Assert(splitTime('Fr:ed') is: false)
		Assert(splitTime('Fr:15') is: false)
		Assert(splitTime('15:ed') is: false)
		Assert(splitTime('15:ed am') is: false)
		}

	Test_timecontrolSplitTime()
		{
		tc = ChooseTimeControl.ChooseTimeControl_timecontrol
		fn = tc[tc.MembersIf({ it.Has?('splitTime') })[0]] // there should only be one
		Assert(fn('15:24 am') is: #(hour: 15, minute: '24', suffix: 'am'))
		Assert(fn('524am') is: #(hour: 5, minute: '24', suffix: 'am'))
		Assert(fn('524') is: #(hour: 5, minute: '24', suffix: ""))
		Assert(fn('8') is: #(hour: 0, minute: '', suffix: ""))
		Assert(fn('') is: #(hour: 0, minute: '', suffix: ""))
		Assert(fn('garbage data pm') is: #(hour: 0, minute: '', suffix: ""))
		}
	}
