// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
// used by BookEditControl
Controller
	{
	Title: "Find in Folders"
	CallClass(data, valid = function (unused) { "" })
		{
		ToolDialog(_hwnd, [this, data, valid], keep_size: false)
		}
	New(data, .valid)
		{
		.Data.Set(data)
		}
	Controls: #(Record // mostly duplicate of FindDialog :-(
		(Horz
			(Vert
				(Pair
					(Static 'Find in name')
					(FieldHistory, font: '@mono', size: '+1', width: 30,
						name: "name"))
				(Skip 5)
				(Pair
					(Static 'Find in text')
					(FieldHistory, font: '@mono', size: '+1', width: 30
						name: "find")) // to match ScintillaControl
				(Skip 5),
				(Horz
					Skip
					(Vert
						(CheckBox, "Match case", name: "case")
						(CheckBox, "Match whole words", name: "word")
						(CheckBox, "Regular expression", name: "regex")
						)
					)
				)
			Skip
			(Vert
				(Button, "Find First", xstretch: 0)
				(Skip 8)
				(Button, "Cancel", xstretch: 0))
			)
		)
	DefaultButton: "Find First"
	On_Find_First()
		{
		data = .Data.Get()
		if "" isnt (err = .something_to_find(data)) or
			"" isnt (err = (.valid)(data))
			{
			.AlertWarn('Find in Folders', err)
			return
			}
		.Window.Result(data)
		}
	something_to_find(data)
		{
		return data.find is "" and data.name is ""
			? 'Please enter something to find' : ""
		}
	}
