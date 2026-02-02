// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name:		"Editor"
	ComponentName:	"Scintilla"
	DefaultFontSize: 11
	Unsortable: true

	New(style/*unused*/ = 0, .readonly = false, lexer/*unused*/ = 'NULL',
		wrap = false, set = false, margin/*unused*/ = 12,
		exStyle/*unused*/ = 0, height = false, .tabthrough = false)
		{
		.setReadOnly = .readonly
		widths = #(tab: 4, scroll: 10)
		.SetTABWIDTH(widths.tab)
		.SetWrap(wrap)
		.Send("Data")

		.Map = Object()
		if set isnt false
			.Defer({ .Set(set) })

		.findreplacedata = Record()

		.ComponentArgs = Object(readonly, height, tabthrough)

.lastEvents = Object().Set_default(Object())
		}

	BaseStyling()
		{
		return []
		}

	wrap: false
	SetWrap(.wrap)
		{
		.Act(#SetlineWrapping, wrap)
		}
	GetWrapMode()
		{
		return .wrap is true ? SC.WRAP_WORD : SC.WRAP_NONE
		}

	setReadOnly: false
	GetReadOnly()
		{
		return .readonly is true or .setReadOnly is true
		}
	// used by Addon_speller
	GetReadonly()
		{
		return .GetReadOnly() ? 1 : 0
		}
	SetReadOnly(readonly = true)
		{
		if .readonly
			return
		if readonly and .Dirty?()
			{
// Extra logging for suggestion 36802
calls = GetCallStack(limit: 99)
RemoveAssertsFromCallStack(calls)
// avoiding broadcast/send loop
calls = FormatCallStack(calls, levels: 99).
	Replace('(?q)Container.Broadcast /* stdlib__webgui method */', '').
	Replace('(?q)Container.SetReadOnly /* stdlib__webgui method */', '').Lines().
	RemoveIf({ it.Blank?() }).Join('\n')
			SuneidoLog('ERROR: ScintillaControl SetReadonly(true) when dirty (36802)',
				:calls
				params: Object(dirty: .dirty, modify: .GetModify(),
					alreadyReadOnly: .GetReadOnly(), name: .Name, uniqueId: .UniqueId).
					Merge(.lastEvents))
			SuRenderBackend().DumpStatus('ScintillaControl SetReadonly(true) when dirty')
			}
		return .Act(#SetReadOnly, .setReadOnly = readonly)
		}
	GETREADONLY()
		{
		return .GetReadOnly() ? 1 : 0
		}

	s: ''
	Get()
		{
		return .s
		}
	Set(s)
		{
		.CancelAct('Set')
		.CancelAct('AppendText')
		.s = EnsureCRLF(s)
		.astWriterManager = false
		.Act(#Set, .s)
		.Dirty?(false)
		}
	Scintilla_UpdateValue(change)
		{
		.s = .s[..change.from] $ change.text $ .s[change.to..]
		}

	wordchars: 'zyxwvutsrqponmlkjihgfedcba_ZYXWVUTSRQPONMLKJIHGFEDCBA?9876543210!'
	SetWordChars(.wordchars) {}
	GetWordChars()
		{
		return .wordchars
		}

	GetLength()
		{
		return .s.Size()
		}

	GetTextLength()
		{
		return .GetLength()
		}

	GetRange(start, end) // gets from start to end - 1
		{
		start = .intoRange(start)
		end = .intoRange(end)
		return .s[start..end]
		}

	intoRange(i)
		{
		return Max(0, Min(i, .GetTextLength()))
		}

	GetCurrentReference()
		{
		pos = .GetSelectionStart()
		line = ('\n' $ .s[..pos]).AfterLast('\n').RightTrim('(')
		return line.Extract("[._[:alnum:]?!]+$")
		}

	GetAt(pos)
		{
		return pos < 0 or pos >= .GetTextLength() ? '\x00' : .GetRange(pos, pos + 1)
		}

	Trim() // remove whitespace from beginning and end
		{
		.s = .s.Trim()
		.Act(#Set, .s)
		return this // to allow chaining
		}

	AppendText(s)
		{
		s = EnsureCRLF(s)
		.Act(#AppendText, s)
		.s $= s
		}

	dirty: false
	Dirty?(dirty = '')
		{
		Assert(dirty is true or dirty is false or dirty is "")
		if dirty isnt ''
			.dirty = dirty
		return .dirty
		}

	modified: false
	GetModify()
		{
		return .modified
		}

	EN_CHANGE()
		{
.lastEvents['EN_CHANGE'] = Object(
	eventId: SuRenderBackend().SuRenderBackend_eventId,
	t: Timestamp())
		.astWriterManager = false
		.Send("EN_CHANGE")
		.dirty = true
		return 0
		}

	SCN_MODIFIED(unused)
		{
.lastEvents['SCN_MODIFIED'] = Object(
	eventId: SuRenderBackend().SuRenderBackend_eventId,
	t: Timestamp(), readonly: .GetReadOnly(), focus: .HasFocus?(), dirty?: .Dirty?())
		// .Paste may change the value without having the focus in the field
		if .GetReadOnly() is false and not .HasFocus?() and .Dirty?()
			.Send("NewValue", .Get())
		}

	LBUTTONUP(pos)
		{
		.Send('Scintilla_LButtonUp', pos)
		}

	Scintilla_SetValue() {}

	SCEN_KILLFOCUS()
		{
.lastEvents['SCEN_KILLFOCUS'] = Object(
	eventId: SuRenderBackend().SuRenderBackend_eventId,
	t: Timestamp(), dirty?: .Dirty?())
		.closeAutoc()
		.Send('Scintilla_KillFocus')
		if (.Dirty?())
			.Send("NewValue", .Get())
		return 0
		}

	KILLFOCUS()
		{
		.SCEN_KILLFOCUS()
		}

	SCEN_SETFOCUS()
		{
		.Send('Scintilla_SetFocus')
		return 0
		}

	SETFOCUS()
		{
		.SCEN_SETFOCUS()
		}

	selection: #(anchor: 0, head: 0)
	SU_UPDATESELECT(.selection)
		{
		.curPos = .selection.head
		.closeAutoc()
		}

	firstVisibleLine: 0
	SU_SYNCFIRSTVISIBLELINE(.firstVisibleLine)
		{
		}

	GetFirstVisibleLine()
		{
		return .firstVisibleLine
		}

	SetFirstVisibleLine(line)
		{
		.firstVisibleLine = line
		.Act('SetFirstVisibleLine', line)
		}

	GetSelectionStart()
		{
		return Min(.selection.anchor, .selection.head)
		}
	GetSelectionEnd()
		{
		return Max(.selection.anchor, .selection.head)
		}
	SelSize()
		{
		return .GetSelectionEnd() - .GetSelectionStart()
		}
	GetSelText()
		{
		sel = .GetSelect()
		return .GetRange(sel.cpMin, sel.cpMax)
		}
	GetSelect()
		{
		return Object(
			cpMin: .GetSelectionStart(),
			cpMax: .GetSelectionEnd())
		}
	SetSelect(i, n = 0)
		{
		.SetSel(i, i + n)
		}
	SetVisibleSelect(i, n)
		{
		.EnsureRangeVisible(i, i + n)
		.SetSelect(i, n)
		}
	AddSelection(head, anchor)
		{
		.SU_UPDATESELECT([:anchor, :head])
		.Act('AddSel', head, anchor)
		}
	DoWithCurrentPos(block)
		{
		.Act('SavePos')
		block()
		.Act('RestorePos')
		}

	SCN_DOUBLECLICK()
		{
		.Send('Scintilla_DoubleClick')
		}

	curPos: false
	GetCurrentPos()
		{
		return .curPos
		}

	GotoLine(line, noFocus? = false)
		{
		if noFocus? is false
			.SetFocus()
		pos = .PositionFromLine(line)
		.SetVisibleSelect(pos, 0)
		}
	LineFromPosition(pos = false)
		{
		if pos is false
			pos = .GetSelectionStart()
		line = 0
		.s.ForEach1of('\n')
			{
			if it > pos
				return line
			line++
			}
		return line
		}

	PositionFromLine(line)
		{
		return .s.StartPositionOfLine(line)
		}

	GetLineCount()
		{
		c = .s.LineCount()
		if .s.Suffix?('\n')
			++c
		return c
		}

	GetColumn(pos = false)
		{
		if pos is false
			pos = .GetSelectionStart()
		lineStart = 0
		.s.ForEach1of('\n')
			{
			if it > pos
				return .s[lineStart..pos].Detab().Size()
			lineStart = it + 1
			}
		return .s[lineStart..pos].Detab().Size()
		}

	GetLineEndPosition(line)
		{
		cur = 0
		.s.ForEach1of('\n')
			{
			if cur is line
				return it - 1
			cur++
			}
		return .s.Size()
		}

	UPDATEUI()
		{
		}

	On_Find()
		{
		s = .GetSelText()
		if s > "" and not s.Has?('\n')
			.findreplacedata.find = s
		_hwnd = .WindowHwnd()
		x = FindDialog(.findreplacedata)
		.DoFind(x)
		}

	DoFind(x)
		{
		if x is 'prev'
			return .On_Find_Previous()
		else if x is 'next'
			return .On_Find_Next()
		}

	On_Find_Next()
		{
		return .findAndMatch()
		}
	On_Find_Previous()
		{
		return .findAndMatch(prev:)
		}

	findAndMatch(prev = false)
		{
		getSelect = prev is false ? .GetSelect().cpMax : .GetSelect().cpMin
		if false is match =
			Find.DoFind(.SearchText(), getSelect, .findreplacedata, :prev)
			{
			.Send('UpdateOccurrence', num: 0, count: 0)
			return false
			}
		findOb = .numOfMatch(match)
		.Send('UpdateOccurrence', num: findOb.num, count: findOb.count)
		.SetVisibleSelect(match[0], match[1])
		return true
		}
	numOfMatch(match)
		{
		matches = Find.FindAll(.SearchText(), .findreplacedata)
		num = matches.FindIf({ it is match }) + 1
		return Object(:num, count: matches.Size())
		}

	SearchText()
		{
		return .Get()
		}

	Context_Menu: (
		"&Undo\tCtrl+Z", "&Redo\tCtrl+Y", "",
		"Cu&t\tCtrl+X", "&Copy\tCtrl+C", "&Paste\tCtrl+V", "&Delete", "",
		"Select &All\tCtrl+A", "Find...\tCtrl+F")
	ContextMenu(x, y)
		{
		i = ContextMenu(.Context_Menu).ShowCall(this, x, y)
		if i is false or i <= 0
			.EnsureSelect()
		return 0
		}

	EnsureSelect()
		{
		.On_Select_All()
		.SetSel(.GetSelectionStart(), .GetSelectionEnd())
		}

	On_Context_Cut()
		{ .On_Cut() }
	On_Context_Copy()
		{ .On_Copy() }
	On_Context_Paste()
		{ .On_Paste() }
	On_Context_Undo()
		{ .On_Undo() }
	On_Context_Redo()
		{ .On_Redo() }
	On_Context_Delete()
		{ .On_Delete() }
	On_Context_Select_All()
		{ .On_Select_All() }
	On_Context_Find()
		{ .On_Find() }

	On_Cut()
		{
		if .GetReadOnly() is true
			return
		.Act('CUT')
		}
	On_Copy()
		{
		.Act('COPY')
		}
	On_Paste()
		{
		if .GetReadOnly() is true
			return
		.Act('PASTE')
		}
	On_Undo()
		{
		if .GetReadOnly() is true
			return
		.Act('UNDO')
		}
	On_Redo()
		{
		if .GetReadOnly() is true
			return
		.Act('REDO')
		}
	On_Delete()
		{
		if .GetReadOnly() is true
			return
		.Act('DELETE')
		}
	On_Select_All()
		{
		.Act('SELECTALL')
		}

	MarkersChanged()
		{
		.Send(#Scintilla_MarkersChanged)
		}

	PasteOverAll(s)
		{
		.On_Select_All()
		.Act('Paste', s)
		.SetSelect(0)
		}

	Paste(s)
		{
		if .GetReadOnly() is true
			return
		.Act('Paste', s)
		}

	listctrl: false
	autocPrefix: 0
	SCIAutocShow(.autocPrefix, matches)
		{
		.closeAutoc()
		.listctrl = AutoChooseList(this, matches)
		.syncAutocStatus(true)
		}

	AutocCancel()
		{
		.closeAutoc()
		}

	AutocSelection(s)
		{
		for .. .autocPrefix
			.CharLeftExtend()
		.ReplaceSel(s)
		}

	closeAutoc()
		{
		if .listctrl isnt false
			{
			.listctrl.DESTROY()
			.listctrl = false
			.syncAutocStatus(false)
			}
		}

	syncAutocStatus(open)
		{
		.Act('SyncAutocStatus', open)
		}

	AUTOC_KEYDOWN(key)
		{
		if .listctrl is false
			return

		switch (key)
			{
		case #ArrowDown:
			.listctrl.Down()
		case #ArrowUp:
			.listctrl.Up()
		case #Tab, #Enter:
			if "" isnt s = .listctrl.Get()
				.Picked(s)
			.closeAutoc()
		case #Escape:
			.closeAutoc()
		default:
			}
		}

	Picked(s) // Sent from AutoChooseList
		{
		.AutocSelection(s)
		.closeAutoc()
		}

	ListClosed() // sent by listbox
		{
		.listctrl = false
		.syncAutocStatus(false)
		}

	Default(@args)
		{
		.Act(@args)
		}

	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}
