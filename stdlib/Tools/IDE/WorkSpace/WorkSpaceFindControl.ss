// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	New()
		{
		.field = .FindControl('text')
		.selectExclude = .FindControl('exclude')
		.selectCase = .FindControl('case')
		.selectWord = .FindControl('word')
		.selectRegex = .FindControl('regex')
		}

	Controls()
		{
		return Object('Horz'
			Object('FieldHistory', font: '@mono', size: '+1', width: 30,
				field: Object('AutoAction', list: .actionList)
				name: 'text') #Skip
			#(CheckBox, "Exclude"	tabover:, name: 'exclude') #Skip
			#(CheckBox, "Case" 		tabover:, name: "case") #Skip
			#(CheckBox, "Word" 		tabover:, name: "word") #Skip
			#(CheckBox, "Regex" 	tabover:, name: "regex"))
		}

	listCanceled?: false
	actionList(text)
		{
		if .listCanceled?
			return #()

		list = Object()
		.buildList(.suggestRegex?(text), .selectRegex.Get(), 'regex', list)
		.buildList(.suggestCase?(text), .selectCase.Get(), 'case', list)
		.buildList(.suggestWord?(text), .selectWord.Get(), 'word', list)
		return list
		}

	buildList(suggest, current, target, list)
		{
		if suggest is current
			return
		list.Add('Toggle ' $ target.Capitalize() $ ': ' $ Display(suggest))
		}

	GetListForEmpty()
		{
		return .selectExclude.Get() or .selectCase.Get() or .selectWord.Get() or
			.selectRegex.Get()
			? #('Clear all')
			: #()
		}

	SelectAction(action, source)
		{
		if action is 'Clear all'
			{
			.clearAll()
			return
			}

		targetName = action.AfterFirst(' ').BeforeFirst(':').Lower()
		value = action.AfterFirst(': ').SafeEval()
		target = .FindControl(targetName)
		if target isnt false and value isnt target.Get()
			{
			target.Toggle()
			.updateHilite()
			}

		newList = .actionList(.field.Get())
		if newList.NotEmpty?()
			.Defer({ source.OpenList(newList) }, uniqueID: "openList")
		}

	clearAll()
		{
		.clearSelect(.selectExclude)
		.clearSelect(.selectCase)
		.clearSelect(.selectWord)
		.clearSelect(.selectRegex)
		.updateHilite()
		}

	clearSelect(target)
		{
		if target.Get() is true
			target.Toggle()
		}

	CancelList()
		{
		.listCanceled? = true
		}

	prev: ''
	Edit_Change(source)
		{
		text = .field.Get()
		if not text.Prefix?(.prev)
			.listCanceled? = false
		.prev = text

		.updateHilite()
		.Send("Edit_Change", :source)
		}
	NewValue(value, source)
		{
		if source in (.selectCase, .selectRegex)
			.updateHilite()
		.Send("NewValue", value, :source)
		}

	suggestRegex?(text)
		{
		if text.Blank?()
			return false

		if not text.Has1of?('.?*+^$|') and
			text !~ `\\d|\\D|\\s|\\S|\\w|\\W|\[:|\\<|\\>`
			return false

		try
			{
			Tdop(text, type: 'expression', symbols: WorkSpaceFindSymbols())
			return false
			}

		try
			{
			// The assignment and the Type call are temporary to suppress the gSuneido
			// check against useless expression
			a = "" =~ text
			Type(a)
			}
		catch
			return false

		return true
		}

	suggestCase?(text)
		{
		return text =~ '[A-Z]'
		}

	suggestWord?(text)
		{
		return text.GlobalName?() and not Uninit?(text)
		}

	updateHilite()
		{
		text = .field.Get()
		.setColor(.selectRegex, .hiliteRegex?(text))
		.setColor(.selectCase, .hiliteCase?(text))
		.setColor(.selectWord, .hiliteWord?(text))
		}

	hiliteRegex?(text)
		{
		return .selectRegex.Get() isnt .suggestRegex?(text)
		}

	hiliteCase?(text)
		{
		return .selectCase.Get() isnt .suggestCase?(text)
		}

	hiliteWord?(text)
		{
		return .selectWord.Get() isnt .suggestWord?(text)
		}

	setColor(ctrl, hilite)
		{
		ctrl.SetColor(hilite ? CLR.purple : CLR.BLACK)
		ctrl.SetFont(weight: hilite ? 'bold' : '')
		ctrl.Repaint()
		}
	}