// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(hwnd, code, buttons)
		{
		return Dialog(hwnd, [this, code, buttons],
			keep_size: "AiAgentDiff", closeButton?:, title: "Approve Changes")
		}
	New(code, buttons)
		{
		super(["Vert", [CodeViewer, code], buttons])
		}
	On_Allow()
		{
		.Window.Result("allow")
		}
	On_Deny()
		{
		.Window.Result("deny")
		}
	}