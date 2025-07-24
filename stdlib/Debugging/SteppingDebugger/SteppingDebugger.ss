// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Debugger
	{
	Title: "Stepping Debugger"
	ErrorMsg: "Stepping Break Point"
	CallClass(num)
		{
		if not SteppingDebuggerManager().Break?(num)
			return

		.CallStepping(GetActiveWindow(), .ErrorMsg, GetCallStack(skip: 1),
			debuggerNum: num)
		}

	CallStepping(hwnd, err, calls = false, debuggerNum = false)
		{
		if Suneido.Member?(#SteppingDebugger)
			Suneido.SteppingDebugger.Reload(hwnd, err, calls, debuggerNum)
		else
			ToolDialog(hwnd, Object(this, hwnd, err, calls, debuggerNum), border: 0)
		}

	New(hwnd, .err, calls, .debuggerNum = false)
		{
		super(hwnd, err, calls)
		Suneido.SteppingDebugger = this
		.continueBtn = .FindControl('Continue')
		.nextStepBtn = .FindControl('Next_Step')
		}

	Reload(hwnd, .err, calls, .debuggerNum = false)
		{
		.Init(hwnd, err, calls)
		.continueBtn.SetEnabled(true)
		.nextStepBtn.SetEnabled(true)
		.Startup()
		.Window.ActivateDialog()
		}

	ExtraButtons()
		{
		return #(#(Button "Next Step") #Skip #(Button Continue) #Skip)
		}

	ExtraAddons()
		{
		return Object(Addon_stepping_debugger:)
		}

	ShowSource(frame, disassembled = false)
		{
		if false isnt info = .getFrameInfo()
			{
			.SetDebugView(frame.fn, info.source)
			.Source.Set(info.source)
			.Source.SetSelect(info.i, info.n)
			return
			}
		super.ShowSource(frame, disassembled)
		}

	getFrameInfo()
		{
		if .CurCall isnt 0
			return false

		if false is debugger = SteppingDebuggerManager().GetDebugger(.debuggerNum)
			return false

		if false is source = SteppingDebuggerManager().
			GetSource(debugger.lib, debugger.name)
			return false

		return Object(:source, lib: debugger.lib, name: debugger.name,
			i: debugger.stmtNode.Position - 1, n: debugger.stmtNode.Length)
		}

	CurrentTable()
		{
		if false is info = .getFrameInfo()
			return ''
		return info.lib
		}

	CurrentName()
		{
		if false is info = .getFrameInfo()
			return ''
		return info.name
		}

	RemoveAllBreakPointsOnDestroy?()
		{
		return false
		}

	On_Next_Step()
		{
		SteppingDebuggerManager().Step? = true
		.unblock()
		}

	On_Continue()
		{
		SteppingDebuggerManager().Step? = false
		.unblock()
		}

	unblock()
		{
		.SetupError('Idling...', color: CLR.GREEN)
		.continueBtn.SetEnabled(false)
		.nextStepBtn.SetEnabled(false)
		.Window.Unblock()
		}

	Destroy()
		{
		SteppingDebuggerManager().Step? = false
		Suneido.Delete(#SteppingDebugger)
		super.Destroy()
		}
	}
