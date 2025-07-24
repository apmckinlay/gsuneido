// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{

	Name: "Accordion"

	New(@controls)
		{
		super(.tabs(controls))
		}
	tabs(controls)
		{
		tabs = Object('Vert')
		open = controls.GetDefault(#open, false)
		controls.Delete(#open)

		for m in controls.Members(list:)
			{
			c = controls[m]
			et = c.Copy()
			et.Add('Expand', at: 0)
			et.leftAlign = true
			et.saveExpandName = c.GetDefault(#saveExpandName, false)
			et.open = c.GetDefault(#open, open)
			tabs.Add(et)
			}
		tabs.Add('EtchedLine')
		return tabs
		}

	ExpandAll()
		{
		for c in .getExpands()
			c.Expand()
		}
	ContractAll()
		{
		for c in .getExpands()
			c.Contract()
		}

	Static_ContextMenu(x, y, source)
		{
		if source.Send('SourceFromHeading?', evtSource: source) is true
			return .contextMenu(x, y)
		return 0
		}
	EnhancedButton_ContextMenu(x, y)
		{ .contextMenu(x, y) }
	contextMenu(x, y)
		{
		ContextMenu(#("Expand All", "Contract All")).ShowCall(this, x, y)
		return true // to avoid extra context menu from static
		}
	On_Context_Expand_All()
		{ .ExpandAll() }
	On_Context_Contract_All()
		{ .ContractAll() }

	Resize(x, y, w, h)
		{
		labels = .getExpands().Map({ it.FindControl(#tabname) })
		xmin = 0
		for c in labels
			if c.Xmin > xmin
				xmin = c.Xmin
		for c in labels
			c.Xmin = xmin
		super.Resize(x, y, w, h)
		}
	getExpands()
		{
		return .Vert.GetChildren()[..-1] // remove etched line
		}
	ExpandByName(name)
		{
		for expand in .Vert.GetChildren()[..-1]
			{
			if expand.Label is name
				{
				expand.Expand()
				return expand
				}
			}
		}
	}
