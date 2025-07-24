// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.

// TODO: add option to flag all changes

Refactor
	{
	SelectWord: true
	IdleTime: 250 //ms
	DiffPos: 3
	Init(data)
		{
		if (data.select.cpMin < data.select.cpMax)
			{
			selected = data.text[data.select.cpMin ::
				data.select.cpMax - data.select.cpMin]
			if selected =~ .NamePattern
				data.from = selected
			}
		if data.from isnt ""
			data.ctrl.Data.GetControl('to').SetFocus()
		return true
		}
	Controls()
		{
		return Object('Vert'
			#(Pair (Static From) (Field font: '@mono', name: from))
			'Skip'
			#(Pair (Static To) (Horz
				(Field font: '@mono', name: to)
				Skip
				(Pair (Static '') (CheckBox 'Replace in Comments' name: inComments))))
			#(Diff2 '', '', '', '', 'From', 'To')
			name: 'renameVert')
		}
	Errors(data)
		{
		errs = Object()
		if data.from is ""
			errs.Add("From is required")
		else if data.from !~ .NamePattern
			errs.Add("From must be a valid name")
		if data.to is ""
			errs.Add("To is required")
		else if data.to !~ .NamePattern
			errs.Add("To must be a valid name")
		return errs.Join('\n')
		}
	Warnings(data)
		{
		if .ToExists?(data.text, data.to)
			return '"' $ data.to $ '" is already used'
		return ""
		}
	ToExists?(text, name)
		{
		ScannerMap(text)
			{ |prev2, prev, token, next|
			if .Matches?(name, prev2, prev, token, next)
				return true
			token // returned
			}
		return false
		}

	Process(data)
		{
		data.text = .Rename(data.text, data.from, data.to, data.inComments)
		return true
		}

	Rename(text, from, to, inComments = false)
		{
		return ScannerMap(text)
			{ |prev2, prev, token, next|
			if .Matches?(from, prev2, prev, token, next)
				token = to
			else if inComments is true and token =~ "^(//|/[*])"
				token = token.Replace('\<' $ from $ '\>', to)
			token // returned
			}
		}
	FindMatch(text, name)
		{
		prev2 = prev = ''
		scanner = Scanner(text)
		pos = nextPos = 0
		if scanner is token = scanner.Next()
			return -1
		while scanner isnt next = scanner.Next()
			{
			if .Matches?(name, prev2, prev, token, next)
				return pos
			prev2 = prev
			prev = token
			token = next
			pos = nextPos
			nextPos = scanner.Position()
			}
		return -1
		}

	InitializePreview(data, methodText)
		{
		.vert = data.ctrl.FindControl('renameVert')
		.vert.Remove(.DiffPos)
		.vert.Insert(.DiffPos, Object('Diff2', methodText, .GetRenamed(data, methodText),
			data.library, data.name, 'From', 'To'))
		data.Observer(.Change)
		.diff = data.ctrl.FindControl('Diff')
		.Timer = IdleTimer(.IdleTime, { .setPreview(data, methodText) })
		}

	GetRenamed(data, methodText)
		{
		return .Rename(methodText, data.from, data.to, data.inComments)
		}

	setPreview(data, methodText)
		{
		.diff.UpdateList(methodText,
			.Rename(methodText, data.from, data.to, data.inComments))
		}
	}
