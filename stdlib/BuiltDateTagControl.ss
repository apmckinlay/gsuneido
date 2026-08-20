// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(cur)
		{
		.Data.SetField(#dir,
			cur isnt false and cur.dir is '<' ? #before : "on or after")
		.Data.SetField(#date,
			cur isnt false and cur.date isnt false ? cur.date : BuiltDate().NoTime())

		if cur is false
			.FindControl(#Remove).SetEnabled(false)
		}

	Controls()
		{
		return [#Record,
			[#Vert,
				#(Static, "Skip code checking and test running unless the exe was built"),
				#Skip,
				#(Horz, (Skip, 16),
					(RadioButtons, "on or after", before, horz:, name: dir), Fill),
				#Skip,
				#(Horz, (Skip, 16), (Date, name: date, mandatory:), Fill),
				#Skip,
				#(Horz, (Skip, 16), (MonthCal, name: date), Fill),
				#(Skip, 5),
				#(Horz, Fill, (Button, Add, xmin: 50), Skip,
					(Button, Remove, xmin: 50), Skip, (Button, Cancel, xmin: 50))]]
		}

	On_OK()
		{
		.On_Add()
		}

	On_Add()
		{
		data = .Data.Get()
		if Date?(data.date)
			.Window.Result([dir: data.dir is #before ? '<' : '>', date: data.date])
		}

	On_Remove()
		{
		.Window.Result(#remove)
		}
	}
