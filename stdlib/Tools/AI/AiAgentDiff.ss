// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(hwnd, before, after, lib = '', recname = '')
		{
		return Dialog(hwnd, [this, before, after, :lib, :recname],
			keep_size: "AiAgentDiff", closeButton?:,
			title: "Approve Changes" $ Opt(" - ", lib, ':' $ recname))
		}
	New(before, after, lib = '', recname = '')
		{
		super(["Vert",
			["Diff2", after, before, lib, recname, "After", "Before",
				extraControls: AiAgentControl.ApproveButtons, newOnRight?:]
			#Skip,
			#(ScintillaAddons, name: "feedback", wrap:, xstretch: 1, ystretch: 0)])
		}
	On_Allow()
		{
		.Window.Result(Object("allow", feedback: .FindControl("feedback").Get()))
		}
	On_Deny()
		{
		.Window.Result(Object("deny", feedback: .FindControl("feedback").Get()))
		}
	}