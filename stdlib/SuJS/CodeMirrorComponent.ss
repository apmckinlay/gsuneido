// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
Component
	{
	Xstretch: 1
	Ystretch: 1
	New(.readonly = false)
		{
		.CreateElement('textarea')

		.CodeMirror = SuUI.GetCodeMirror()
		.cm = .CodeMirror.FromTextArea(.El)
		.SetEl(.cm.GetWrapperElement())

		.El.SetStyle("border", "1px solid gray")
		.initModeAndTheme()
		.initExtraKeys()
		.initStyles()

		.cm.On("changes", .EN_CHANGE)
		.SetMinSize()
		}

	initModeAndTheme()
		{
		.cm.SetOption("mode", "suneido")
		.cm.SetOption("theme", "suneido")
		}

	initExtraKeys()
		{
		.cm.SetOption("extraKeys", Object(
			"Ctrl-Q": .toggleFold,
			"Ctrl-/": "toggleComment",
			"Shift-Ctrl-/": .toggleBlockComment
			"Ctrl-G": .goto
			))
		}

	toggleFold(cm)
		{
		Print(cm.GetCursor())
		cm.FoldCode(cm.GetCursor())
		}

	toggleBlockComment(cm)
		{
		from = cm.GetCursor("from")
		to = cm.GetCursor("to")
		if cm.Uncomment(from, to) is false
			{
			cm.BlockComment(from, to, Object(fullLines: false))
			}
		}

	goto(cm)
		{
		selection = cm.GetSelection()
		if selection.Blank?()
			return
		.Event('Send', 'On_GO_TO', selection)
		}

	initStyles()
		{
		.cm.SetOption("lineNumbers", true)
		.cm.SetOption("lineWrapping", true)
		.cm.SetOption("indentUnit", 4/*=indent unit*/)
		.cm.SetOption("indentWithTabs", true)
		.cm.SetOption("matchBrackets", true)
		.cm.SetOption("foldGutter", true)
		.cm.SetOption("gutters", [
			"Suneido-searchgutter", "CodeMirror-foldgutter", "CodeMirror-linenumbers"])
		.cm.SetOption("autoCloseBrackets", true)
		.cm.SetOption("styleActiveLine", true)
		.cm.SetOption("highlightWords", true)
		.cm.SetOption("readOnly", .readonly)
		}

	Set(value)
		{
		if (not String?(value))
			value = Display(value)
		.cm.SetValue(value)
		}

	Get()
		{
		return .cm.GetValue()
		}

	EN_CHANGE(@unused)
		{
		.Event('EN_CHANGE', .Get())
		return 0
		}

	SetReadOnly(readOnly)
		{
		if (.readonly)
			return
		.cm.SetOption("readOnly", readOnly is true)
		super.SetReadOnly(readOnly)
		}

	GetReadOnly()
		{
		return .cm.GetOption("readOnly")
		}

	SetCursor(line, ch)
		{
		.cm.SetCursor(line, ch)
		}
	}
