// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Html_ahref_Control
	{
	New(text, .href, .prefix = 'http://', .valign = false)
		{
		super(text, href)
		.ToolTip(.prefix $ .href)
		.valign = .valign in ('bottom', 'top', 'vcenter') ? DT[.valign.Upper()] : 0
		}
	LBUTTONUP()
		{
		ShellExecute(.WindowHwnd(), 'open', .prefix $ .href)
		return 0
		}
	DrawFlags()
		{
		return super.DrawFlags() + .valign
		}
	Set(href)
		{
		.href = href
		.ToolTip(.prefix $ .href)
		}
	// Links should be usable, regardless of parent controls being set read only
	SetReadOnly(unused) { }
	}
