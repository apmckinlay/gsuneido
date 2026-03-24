// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(hwnd, before, after, buttons)
		{
		return Dialog(hwnd, [this, before, after, buttons],
			keep_size: "AiAgentDiff", closeButton?:, title: "Approve Changes")
		}
	New(before, after, buttons)
		{
		super(["Diff2", before, after, "", "", "Before", "After", extraControls: buttons])
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