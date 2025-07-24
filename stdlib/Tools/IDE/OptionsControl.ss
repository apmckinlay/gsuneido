// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(options, parent)
		{
		super(.layout(parent))
		for field in options.Members()
			.Data.SetField(field, options[field])
		}
	layout(parent)
		{
		.Title = parent $ ' Options'
		ob = Object('Vert')
		Plugins().ForeachContribution('Options', parent)
			{ |x|
			ob.Add(Object(x.control, x.prompt, name: x.name))
			}
		ob.Add('Skip', 'OkCancel')
		return Object('Record', ob)
		}
	On_OK()
		{
		.Window.Result(.Data.Get())
		}
	}