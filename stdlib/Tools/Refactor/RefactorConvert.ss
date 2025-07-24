// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
// abstract class for simple conversion refactors
Refactor
	{
	DiffPos: 1
	Controls: (Vert
		(Skip 3)
		#(Diff2 '', '', '', '', 'From', 'To')
		name: 'convertVert'
		xmin: 600
		)

	Init(data)
		{
		if true isnt (result = .CanConvert(data))
			{
			.alert(result)
			return false
			}

		.vert = data.ctrl.FindControl('convertVert')
		.vert.Remove(.DiffPos)
		.vert.Insert(.DiffPos, Object('Diff2', data.text,
			.Convert(data.text), data.library, data.name, 'From', 'To'))

		return true
		}

	// return true if can be converted, otherwise a string explaining why it can't
	CanConvert(data /*unused*/)
		{
		return true
		}

	alert(msg, warning = false)
		{
		Alert(msg, .Name, flags: warning ? MB.ICONWARNING : MB.ICONINFORMATION)
		}

	Process(data)
		{
		data.text = .Convert(data.text)
		return true
		}

	Convert(text /*unused*/)
		{
		throw 'must be defined by derived class'
		}
	}