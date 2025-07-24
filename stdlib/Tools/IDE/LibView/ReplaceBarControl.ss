// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// used by LibViewControl
// see also FindBarControl
Controller
	{
	Name: ReplaceBar
	New(data)
		{
		.Data.Set(data)
		.replacetext = .FindControl('replace')
		}
	Controls: (Record
		(Vert
			(ReplaceBarHorz
				(ReplaceStatic 'Replace' weight: semibold color: 0x808080,
					justify: RIGHT)
				(Skip 4)
				(FieldHistory, font: '@mono', size: '+2', width: 10,
					trim: false, name: "replace")
				Skip
				(Button, "Replace Current", tip: "Replace current and find next (F8)")
				Skip
				(Button, "In Select")
				Skip
				(Button, "Replace ALL", tip: "Replace all occurrences")
				Skip)
			(Skip 2))
		)
	GetText()
		{ return .replacetext.Get() }
	SetText(text)
		{ .replacetext.Set(text) }
	Select()
		{
		.SetFocus()
		.replacetext.Field.SelectAll()
		}
	SetFocus()
		{
		.replacetext.SetFocus()
		}
	FieldEscape()
		{
		.Send('On_FindBar_Close')
		}
	On_Replace_Current()
		{
		.Send('On_ReplaceBar_ReplaceCurrent')
		}
	On_In_Select()
		{
		.Send('On_ReplaceBar_ReplaceInSelection')
		}
	On_Replace_ALL()
		{
		.Send('On_ReplaceBar_ReplaceAll')
		}
	}
