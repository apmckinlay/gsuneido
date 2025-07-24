// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Profile Results"
	Xstretch: 1
	Ystretch: 1
	CallClass(results, reps = 1)
		{
		Window([this, results, :reps], keep_placement:)
		}
	New(.results, .reps)
		{
		data = .BuildDataValues(.results, reps)
		.list = .FindControl("List")
		.list.Set(data)
		}
	columns: (profile_name, profile_self, profile_total, profile_calls)
	Controls()
		{
		return Object('Vert',
			['Border',
				['ListStretch', .columns, defWidth: 160,
					stretchColumn: 'profile_name', columnsSaveName: .Title],
				border: 5],
			['Horz',
				'Skip',
				#(Static, "double-click to go to definition"),
				'Fill',
				['Static', .reps $ ' reps'],
				'Skip']
			#(Skip, medium:)
			)
		}
	List_WantNewRow()
		{ return false }
	List_WantEditField(unused)
		{ return false }
	List_DoubleClick(@unused)
		{
		.goto()
		return true
		}
	goto()
		{
		if false isnt x = .get_selected()
			GotoLibView(x.profile_name.BeforeFirst(' '))
		}
	get_selected()
		{
		sel = .list.GetSelection()
		data = .list.Get()
		return data[sel[0]]
		}
	BuildDataValues(results, reps)
		{
		if results.Empty?()
			return results
		results.RemoveIf({ it.calls < reps })
		results.Sort!({|x,y| x.self > y.self }) // reverse
		totalTime = results.MaxWith({ it.total }).total
		for ob in results
			{
			ob.profile_name = Unprivatize(ob.name)
			ob.profile_calls = (ob.calls / reps).Round(0)
			ob.profile_self = .percent(ob.self, totalTime)
			ob.profile_total = .percent(ob.total, totalTime)
			}
		results.RemoveIf({ it.profile_self < .1 /*= anything less is irrelevant */ })
		return results
		}
	percent(n, tot)
		{
		precision = 2
		return (100 * n / tot).RoundToPrecision(precision) /*= percentage */
		}
	}