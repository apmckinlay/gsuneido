// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// used by ScintillaControl
// see also FindDialog
// editor must support .On_Find_Next() and .ReplaceOne()
Controller
	{
	Title: "Replace"
	CallClass(editor, data)
		{
		ToolDialog(_hwnd, [this, editor, data])
		}
	New(.editor, data)
		{
		.Data.Set(data)
		}
	Controls:
		#(Record
			(Horz
				(Vert
					(Pair
						(Static 'Find what')
						(FieldHistory, font: '@mono', width: 30, trim: false,
							name: "find"))
					(Skip 6)
					(Pair
						(Static 'Replace with')
						(FieldHistory, font: '@mono', width: 30, trim: false,
							name: "replace"))
					Skip
					(Horz
						Skip
						(Vert
							Fill
							(CheckBox, "Match case", name: "case")
							(CheckBox, "Match whole words", name: "word")
							(CheckBox, "Regular expression", name: "regex")
							Fill
							)
						Fill
						(RadioButtons, "Selection", "Entire text",
							label: "Replace in:", name: 'replaceIn')
						Fill
						)
					)
				Skip
				(Vert
					(Button, "Find Next", xstretch: 0)
					(Skip 6)
					(Button, "Replace", xstretch: 0)
					(Skip 6)
					(Button, "Replace All", xstretch: 0)
					(Skip 6)
					(Button, "Cancel", xstretch: 0)
					)
				)
			)
	On_Find_Next()
		{
		.Data.HandleFocus()
		.editor.On_Find_Next()
		}
	On_Replace()
		{
		.Data.HandleFocus()
		.editor.ReplaceOne()
		if .Data.Get().replaceIn isnt "Selection"
			.editor.On_Find_Next()
		}
	On_Replace_All()
		{
		.Window.Result('all')
		}
	}
