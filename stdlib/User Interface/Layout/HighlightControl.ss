// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Value: 'not-empty'
	Color: 'Highlight'
	New(@args)
		{
		super(@args)
		.size = args.GetDefault('size', '')
		.prompt = .FindControl('Static')

		if 0 is .rc = .Send('GetRecordControl')
			throw "HighlightControl must be used in a RecordControl"

		.rc.AddObserver(.observer)
		.rc.AddSetObserver(.highlight)
		.Recalc()
		}
	Recalc()
		{
		.Top = .GetChild().Top
		.Left = .GetChild().Left
		super.Recalc()
		}
	Data(source)
		{
		.field = source
		.Send(#Data, :source)
		}
	observer(member)
		{
		if member is .field.Name
			.highlight()
		}
	hilite: false
	highlight()
		{
		if .Destroyed?()
			return
		value = .field.Get()
		if .Value is 'not-empty'
			hilite = value isnt "" and value isnt false
		else
			hilite = value is .Value

		if hilite is .hilite
			return
		.hilite = hilite
		.prompt.SetFont(weight: hilite ? 900 : 400, size: .size)
		.prompt.SetColor(hilite ? .Color : 0)
		}

	Destroy()
		{
		.rc.RemoveObserver(.observer)
		.rc.RemoveSetObserver(.highlight)
		super.Destroy()
		}
	}