// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.

// TODO handle adding braces when necessary

Refactor
	{
	Name: 'Extract Variable'
	Desc: 'Extract the selected expression to a new local variable'
	IdleTime: 250 //ms
	DiffPos: 2
	Controls: (Vert
		(Pair (Static 'Name') (Field, name: 'var_name'))
		(Skip 3)
		(Diff2 '', '', '', '', 'From', 'To')
		name: 'extractVert'
		xmin: 600
		)
	Init(data)
		{
		if .validate(data) is false
			return false
		.vert = data.ctrl.FindControl('extractVert')
		.setup(data)
		.timer = IdleTimer(.IdleTime, .setPreview)
		return true
		}
	validate(data)
		{
		.selection = data.text[data.select.cpMin ..	data.select.cpMax]
		if not .isExpression?(.selection)
			{
			.alert('Please select a valid expression')
			return false
			}
		return true
		}
	setup(data)
		{
		.lineStartPos = data.editor.PositionFromLine(-1)
		after = data.text[.lineStartPos ..]
		lineEndPos = after.Find('\n', data.select.cpMax - .lineStartPos)
		.line = after[.. lineEndPos]
		.indent = .line.Extract('^[ \t]*')
		data.var_name = "???"
		.data = data
		.change_name(true)
		data.Observer(.Change)
		}
	isExpression?(expr)
		{
		try return expr isnt "" and
			Function?(("function () { " $ expr $ " }").Compile())
		return false
		}
	Change(member)
		{
		if member is 'var_name' and .data.var_name =~ .name_pat
			.change_name()
		}
	change_name(fromInit? = false)
		{
		linePos = .data.select.cpMin - .lineStartPos
		selSize = .selection.Size()
		.afterChange =
			.make_assign(.data.var_name, .selection).Replace("^\t", "") $ '\n' $
			.line.ReplaceSubstr(linePos, selSize, .data.var_name).Trim()
		if fromInit?
			{
			.vert.Remove(.DiffPos)
			.vert.Insert(.DiffPos, Object('Diff2', .line, .afterChange, .data.library,
				.data.name, 'From', 'To'))
			.diff = .data.ctrl.FindControl('Diff')
			.setPreview()
			}
		else
			.timer.Reset()
		}
	setPreview()
		{
		.diff.UpdateList(.line, .afterChange)
		}
	alert(msg, warning = false)
		{
		Alert(msg, .Name, flags: warning ? MB.ICONWARNING : MB.ICONINFORMATION)
		}

	Errors(data)
		{
		if data.var_name !~ .name_pat
			return "Invalid method name"
		if .var_used?(data.text, .lineStartPos, data.var_name)
			return "Variable name already used"
		return ""
		}
	var_used?(text, pos, var)
		{
		range = ClassHelp.MethodRange(text, pos)
		method_text = text[range.from .. range.to]
		return ScannerHas?(method_text, var)
		}
	name_pat: '^[[:lower:]][_[:alpha:][:digit:]]*[?!]?$'

	Process(data)
		{
		data.text = .Extract(data.text, data.select.cpMin, .selection, data.var_name)
		return true
		}

	Extract(text, pos, selection, name)
		{
		assign = .indent $ .make_assign(name, selection) $ '\n'
		text = text.
			ReplaceSubstr(pos, selection.Size(), name).
			ReplaceSubstr(.lineStartPos, 0, assign)

		return text
		}
	make_assign(name, selection)
		{
		return name $ ' = ' $ selection
		}
	}
