// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidData()
		{
		Assert(UOMControl.ValidData?(''))
		Assert(UOMControl.ValidData?(' lbs') is: false)
		Assert(UOMControl.ValidData?(5) is: false)
		Assert(UOMControl.ValidData?(5, uom_optional:))
		Assert(UOMControl.ValidData?('5 lbs'))
		Assert(UOMControl.ValidData?('5') is: false)
		}

	Test_validCheck()
		{
		Assert(UOMControl.UOMControl_validCheck?("", true, false))
		Assert(UOMControl.UOMControl_validCheck?(" ", true, false) is: false)
		Assert(UOMControl.UOMControl_validCheck?("5", true, false) is: false)
		Assert(UOMControl.UOMControl_validCheck?("lbs", true, false) is: false)
		Assert(UOMControl.UOMControl_validCheck?("5 lbs", true, false))

		Assert(UOMControl.UOMControl_validCheck?("", false, true))
		Assert(UOMControl.UOMControl_validCheck?(" ", false, true))
		Assert(UOMControl.UOMControl_validCheck?("5", false, true))
		Assert(UOMControl.UOMControl_validCheck?("lbs", false, true))
		Assert(UOMControl.UOMControl_validCheck?("5 lbs", false, true))

		Assert(UOMControl.UOMControl_validCheck?("", true, true))
		Assert(UOMControl.UOMControl_validCheck?(" ", true, true) is: false)
		Assert(UOMControl.UOMControl_validCheck?("5", true, true))
		Assert(UOMControl.UOMControl_validCheck?("lbs", true, true) is: false)
		Assert(UOMControl.UOMControl_validCheck?("5 lbs", true, true))

		Assert(UOMControl.UOMControl_validCheck?("", false, false))
		Assert(UOMControl.UOMControl_validCheck?(" ", false, false) is: false)
		Assert(UOMControl.UOMControl_validCheck?("5", false, false) is: false)
		Assert(UOMControl.UOMControl_validCheck?("lbs", false, false) is: false)
		Assert(UOMControl.UOMControl_validCheck?("5 lbs", false, false))
		}

	Test_Get()
		{
		mock = Mock()
		mock.UOMControl_value = FakeObject(Get: "12")
		mock.UOMControl_uom = FakeObject(Get: "lb")
		mock.UOMControl_flat_amt = false
		Assert(mock.Eval(UOMControl.Get) is: "12 lb")

		mock.UOMControl_value = FakeObject(Get: "")
		mock.UOMControl_uom = FakeObject(Get: "")
		mock.UOMControl_flat_amt = false
		Assert(mock.Eval(UOMControl.Get) is: "")

		mock.UOMControl_value = FakeObject(Get: "5")
		mock.UOMControl_uom = FakeObject(Get: "")
		mock.UOMControl_flat_amt = false
		Assert(mock.Eval(UOMControl.Get) is: "5")

		mock.UOMControl_value = FakeObject(Get: "")
		mock.UOMControl_uom = FakeObject(Get: "")
		mock.UOMControl_flat_amt = 100
		Assert(mock.Eval(UOMControl.Get) is: 100)

		mock.UOMControl_value = FakeObject(Get: "")
		mock.UOMControl_uom = FakeObject(Get: 'miles')
		mock.UOMControl_flat_amt = false
		Assert(mock.Eval(UOMControl.Get) is: ' miles')
		}
	}