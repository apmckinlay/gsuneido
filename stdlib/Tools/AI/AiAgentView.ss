// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(hwnd, code, lib = '', recname = '')
		{
		return Dialog(hwnd, [this, code], keep_size: "AiAgentView",
			closeButton?:, title: "Approve Changes" $ Opt(" - ", lib, ':' $ recname))
		}
	New(code)
		{
		super(["Vert", [CodeViewer, code],
			#Skip,
			AiAgentControl.ApproveButtons,
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