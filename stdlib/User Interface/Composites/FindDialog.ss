// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// used by ScintillaControl
// see also ReplaceDialog
Controller
	{
	Title: "Find"
	CallClass(data)
		{
		ToolDialog(_hwnd, [this, data])
		}
	New(options)
		{
		.Data.Set(options)
		}
	Commands: (
		(Find_Next,		"F3")
		(Find_Previous,	"Shift+F3"))
	Controls:
		#(Record
			(Horz
				(Vert
					(Pair
						(Static 'Find what')
						(FieldHistory, font: '@mono', size: '+2', width: 30, trim: false,
							name: "find"))
					Skip
					(Horz
						Skip
						(Vert
							(CheckBox, "Match case", name: "case")
							(CheckBox, "Match whole words", name: "word")
							(CheckBox, "Regular expression", name: "regex"))
						)
					)
				Skip
				(Vert
					(Button, "Find Next", xstretch: 0)
					Skip
					(Button, "Find Previous", xstretch: 0)
					Skip
					(Button, "Cancel", xstretch: 0)
					)
				)
			)
	DefaultButton: "Find Next"
	On_Find_Next()
		{
		.Window.Result('next')
		}
	On_Find_Previous()
		{
		.Window.Result('prev')
		}
	}
