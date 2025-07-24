// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_buildSumFields()
		{
		data = [function1: "max", function0: "total", fieldlist1: "Date",
			fieldlist0: "Name"]
		m = SummarizeControl.SummarizeControl_buildSumFields
		ctrl = Object()
		ctrl.SummarizeControl_sf = SelectFields(#(name, date, time))
		ctrl.SummarizeControl_numrows = 8
		Assert(ctrl.Eval(m, [], Object()) is: '')
		Assert(ctrl.Eval(m, data, Object()) is: ', total name, max date')
		}
	}