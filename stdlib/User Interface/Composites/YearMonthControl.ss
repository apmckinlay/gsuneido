// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'YearMonth'
	New(yearmonth = 0)
		{
		super(.controls())
		.lists = .Data.Vert.Horz
		if yearmonth isnt 0 and false isnt Date('#' $ yearmonth $ '01')
			.Set(yearmonth)
		else
			{
			.lists.Year.SetCurSel(.lists.Year.FindString(String(Date().Year())))
			.lists.Month.SetCurSel(0)
			}
		}

	yearsBack: 6
	yearsAhead: 7
	controls()
		{
		yr = Date().Year()
		months = #(Jan Feb Mar Apr May Jun Jul Aug Sep Oct Nov Dec)
		return Object('Record'
			Object('Vert'
				Object('Horz'
					Object('ListBox' Seq(yr - .yearsBack, yr + .yearsAhead) name: 'Year')
					Object('ListBox' months name: 'Month')
					ymin: 175)))
		}

	OK()
		{
		return .Get()
		}

	ListBoxDoubleClick(unused)
		{
		.Send('On_OK')
		}

	Get()
		{
		date = false
		if (false isnt (month = .lists.Month.GetText(.lists.Month.GetSelected())) and
			false isnt (year = .lists.Year.GetText(.lists.Year.GetSelected())))
			date = Date(month $ ' 01 ' $ year)
		return date isnt false ? Number(Display(date.Year()) $ date.Month().Pad(2)) : date
		}

	Set(date)
		{
		date = Date('#' $ date $ '01').Format('yyyyMMM')
		.lists.Year.SetCurSel(.lists.Year.FindString(date[.. 4])) /*= year length */
		.lists.Month.SetCurSel(.lists.Month.FindString(date[-3 ..])) /*= month length */
		}
	}
