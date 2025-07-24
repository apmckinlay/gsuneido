// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'Tags'
	Xmin: 200
	Ymin: false
	New(tags = #())
		{
		super(['Flow', tags.Map({ ['Tag' it] }).Add(#(HandleEnter width: 10))])
		.field = .FindControl('Field')
		if .Ymin is false
			.Ymin = (2.5 /*= baseYMin*/ * .field.Ymin).Int()
		}
	Tag_Remove(text)
		{
		ctrls = .Flow.GetChildren()
		i = ctrls.FindIf({ it.Get() is text })
		.Flow.Remove(i)
		}
	NewValue(text)
		{
		.Flow.Insert(.Flow.Tally() - 1, ['Tag' text])
		.field.Set("")
		}
	}