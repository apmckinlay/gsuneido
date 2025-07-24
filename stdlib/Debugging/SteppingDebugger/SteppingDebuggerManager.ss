// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Singleton
	{
	New()
		{
		//#(debuggerNum:)
		.breakPoints = Object().Set_default(false)

		//#(debuggerNum: #(lib: String, name: String, stmtNode: Tdop))
		.debuggerMap = Object()

		//#(lib: #(
		//	name: #(
		//		src:			String,
		//		convertedSrc: 	String,
		//		node:			Tdop,
		//		debuggerNums: 	Number[],
		//		breakPointNums:	Number[])))
		.sourceMap = Object().Set_default(Object())
		.nextDebuggerNum = 0
		}

	registerDebugger(lib, name, stmtNode)
		{
		.debuggerMap.Add(Object(:lib, :name, :stmtNode), at: .nextDebuggerNum)
		return .nextDebuggerNum++
		}

	ToggleBreakPoint(lib, name, src, pos)
		{
		if lib.Blank?() or name.Blank?()
			return false

		if not .sourceMap[lib].Member?(name)
			.createNewSource(lib, name, src)

		if false is debuggerNum = .findDebugger(.sourceMap[lib][name], pos)
			return false

		.sourceMap[lib][name].breakPointNums.Has?(debuggerNum)
			? .removeBreakPoint(lib, name, debuggerNum)
			: .addBreakPoint(lib, name, debuggerNum)

		return true
		}

	createNewSource(lib, name, src)
		{
		debuggerNums = Object()
		breakPointNums = Object()
		node = Tdop(src)
		writeMgr = AstWriteManager(src, node)
		convertedSrc = AddDebugger(src, node, writeMgr,
			{ |stmtNode|
				num = .registerDebugger(lib, name, stmtNode)
				debuggerNums.Add(num)
				num
			})
		.sourceMap[lib][name] = Object(:src, :convertedSrc, :node, :debuggerNums,
			:breakPointNums)
		}

	findDebugger(sourceItem, pos)
		{
		debuggers = sourceItem.debuggerNums.FindAllIf({
			node = .debuggerMap[it].stmtNode
			node.Position - 1 <= pos and node.Position - 1 + node.Length > pos}).
			Map!({ sourceItem.debuggerNums[it] })

		if debuggers.Empty?()
			return false

		return debuggers.MinWith({ .debuggerMap[it].stmtNode.Length })
		}

	addBreakPoint(lib, name, debuggerNum)
		{
		if .sourceMap[lib][name].breakPointNums.Empty?()
			LibraryOverride(lib, name, .sourceMap[lib][name].convertedSrc)

		.breakPoints[debuggerNum] = true
		.sourceMap[lib][name].breakPointNums.Add(debuggerNum)
		}

	removeBreakPoint(lib, name, debuggerNum)
		{
		.breakPoints.Erase(debuggerNum)
		.sourceMap[lib][name].breakPointNums.Remove(debuggerNum)

		if .sourceMap[lib][name].breakPointNums.Empty?()
			LibraryOverride(lib, name)
		}

	RemoveAllBreakPoints(lib, name)
		{
		if not .HasBreakPoints?(lib, name)
			return

		.breakPoints.Erase(@.sourceMap[lib][name].breakPointNums)
		.sourceMap[lib][name].breakPointNums = Object()
		LibraryOverride(lib, name)
		}

	ClearSource(lib, name)
		{
		if not .sourceMap[lib].Member?(name)
			return

		.RemoveAllBreakPoints(lib, name)
		.debuggerMap.Erase(@.sourceMap[lib][name].debuggerNums)
		.sourceMap[lib].Delete(name)
		}

	HasBreakPoints?(lib, name)
		{
		if not .sourceMap[lib].Member?(name)
			return false

		return .sourceMap[lib][name].breakPointNums.NotEmpty?()
		}

	GetBreakPointRanges(lib, name)
		{
		if not .HasBreakPoints?(lib, name)
			return #()

		return .sourceMap[lib][name].breakPointNums.Map({ Object(
				i: .debuggerMap[it].stmtNode.Position - 1,
				n: .debuggerMap[it].stmtNode.Length) })
		}

	Step?: false
	Break?(num)
		{
		return .Step? is true or .breakPoints[num]
		}

	GetSource(lib, name)
		{
		if not .sourceMap[lib].Member?(name)
			return false

		return .sourceMap[lib][name].src
		}

	GetDebugger(debuggerNum)
		{
		return .debuggerMap.GetDefault(debuggerNum, false)
		}
	}
