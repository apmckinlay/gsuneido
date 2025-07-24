// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'MonthDay'
	dayColumnWidth: 50
	New(monthday = 0)
		{
		.lists = .Data.Vert.Horz
		if ChooseMonthDayControl.DateFromMonthDay(monthday) isnt false
			.Set(monthday)
		else
			{
			.lists.Month.SetCurSel(0)
			.lists.Day.SetCurSel(0)
			}
		.lists.Day.SetColumnWidth(.dayColumnWidth)
		}

	Controls:
		#(Record
			(Vert
				(Horz
					(ListBox #('Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug',
						'Sep', 'Oct', 'Nov', 'Dec') name: 'Month')
					(ListBox #('01', '02', '03', '04', '05', '06', '07', '08', '09', '10',
						'11', '12', '13', '14', '15', '16', '17' ,'18' ,'19', '20', '21',
						'22', '23', '24', '25', '26', '27', '28', '29', '30', '31')
						name: 'Day' multicolumn: , xmin: 170 ymin: 170))))
	OK()
		{
		return .Get()
		}
	ListBoxDoubleClick(@unused)
		{
		.Send('On_OK')
		}

	Get()
		{
		monthIdx = .lists.Month.GetCurSel()
		Assert(monthIdx isnt false)
		++monthIdx
		month = monthIdx.Pad(2)
		day = .lists.Day.Get()
		date = ChooseMonthDayControl.DateFromMonthDay(month $ day)
		return date isnt false ? date.Format('MMdd') : false
		}

	Set(date)
		{
		date = ChooseMonthDayControl.DateFromMonthDay(date).Format('MMMdd')
		.setnumdays(date[.. -2])
		.lists.Month.SetCurSel(.lists.Month.FindString(date[.. -2]))
		dayPos = 3
		.lists.Day.SetCurSel(.lists.Day.FindString(date[dayPos ..]))
		}

	ListBoxSelect(item, source)
		{
		if item is -1 or source is .lists.Day
			return
		day = String(.Get())[2 ..]
		.setnumdays(source.GetText(item))
		.lists.Day.SetCurSel(.lists.Day.FindString(day))
		.Send("NewValue", .Get())
		}

	setnumdays(data)
		{
		numday = 0
		if #(Sep Apr Jun Nov).Has?(data)
			numday = 30
		else if data is 'Feb'
			numday = 28
		else
			numday = 31

		.lists.Day.DeleteAll()
		for(i = 1; i <= numday; ++i)
			.lists.Day.InsertItem(i.Pad(2), i - 1)
		}
	}
