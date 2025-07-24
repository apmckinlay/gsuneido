// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
// TODO: if selection is block, convert any params
// TODO: unindent selection as needed so min indent is \t\t
// TODO: handle formatting for partial line selection
// TODO: allow modifying outputs
// TODO: if selection indented, indent call
// TODO: syntax check method
// TODO: allow overriding "result" variable name
// TODO: check that "result" name isn't already used
// TODO: if selection is partial line, not containing ';', add return
// TODO: add option to replace other occurrences of selection

Refactor
	{
	Name: 'Extract Method'
	Desc: 'Convert the selection to a separate method'
	DiffPos: 6
	IdleTime: 250 //ms
	Controls: (Vert
		(Pair (Static 'Name') (Field, name: 'method_name'))
		(Skip 4)
		(Pair (Static 'Inputs') (ChooseTwoList font: '@mono', readonly:,
			listField: inputlist, xstretch: 1 name: 'inputs'))
		(Skip 4)
		(Pair (Static 'Outputs')
			(Field font: '@mono' readonly:, xstretch: 1 name: 'outputs'))
		(Skip 4)
		#(Diff2 '', '', '', '', 'From', 'To')
		name: 'extractVert'
		xmin: 600
		)
	Init(data)
		{
		if .validate(data) is false
			return false
		data.method_name = "???"
		.data = data

		.preview(data)

		data.Observer(.Change)

		return true
		}
	validate(data)
		{
		if data.select.cpMin >= data.select.cpMax
			{
			.Info('Please select the code you want to extract into a method')
			return false
			}
		.selection = data.text[data.select.cpMin :: data.select.cpMax - data.select.cpMin]
		// TODO: also check for invalid break or continue

		if not ClassHelp.Class?(data.text)
			{
			// TODO: handle embedded classes
			.Info('Extract method can only be used on a class')
			return false
			}

		if ScannerHasIf?(.selection,
			{ it is 'return' or it is 'break' or it is 'continue' })
			{
			.Warn('WARNING: The selection contains "break", "continue", or "return"\n\n' $
				'Additional manual changes to the code will be required.')
			}

		// TODO: check that selection isnt "start" of control struct
		// e.g. not just "if (...)" or "for (...)"

		return true
		}
	preview(data)
		{
		pos = data.select.cpMin

		inputs = .Inputs(data.text, pos, .selection)
		data.inputs = data.inputlist = inputs.Join(', ')

		.outputs = .Outputs(data.text, pos, .selection)
		data.outputs = .outputs.Join(', ')

		.initPreview()
		.change_name(true)
		}
	initPreview()
		{
		.vert = .data.ctrl.FindControl('extractVert')
		.vert.Remove(.DiffPos)
		.vert.Insert(.DiffPos, Object('Diff2', .data.text,
			.Extract(.data.text, .inputs(), .data.select.cpMin, .selection,
				.data.method_name), .data.library, .data.name, 'From', 'To'))
		.diff = .data.ctrl.FindControl('Diff')
		.timer = IdleTimer(.IdleTime, .change_name)
		}
	Change(member)
		{
		if member is 'inputs' or
			(member is 'method_name' and .data.method_name =~ .name_pat)
			.resetTimer()
		}
	resetTimer()
		{
		.timer.Reset()
		}
	change_name(fromInit? = false)
		{
		inputs = .inputs()
		.data.call = .make_call(.data.method_name, inputs, .outputs).
			Replace("^\t\t", "")
		.data.method = .make_method(.data.method_name, .selection, inputs, .outputs).
			Replace("^\t", "")

		if fromInit?
			return

		.updateList()
		}
	updateList()
		{
		.diff.UpdateList(.data.text, .Extract(.data.text, .inputs(), .data.select.cpMin,
			.selection, .data.method_name))
		}
	inputs()
		{
		return .data.inputs.Replace(',([^ ])', `, \1`)
		}

	Errors(data)
		{
		if data.method_name !~ .name_pat
			return "Invalid method name"
		if ClassHelp.Methods(data.text).Has?(data.method_name)
			return "Method name already exists"
		// TODO: also look at superclasses
		return ""
		}
	name_pat: '^[[:alpha:]][_[:alpha:][:digit:]]*[?!]?$'

	Process(data)
		{
		data.text = .Extract(data.text, .inputs(),
			data.select.cpMin, .selection, data.method_name)
		return true
		}

	Extract(text, inputs, pos, selection, name)
		{
		outputs = .Outputs(text, pos, selection)

		method = .make_method(name, selection, inputs, outputs)

		text = ClassHelp.AddMethod(text, pos, method)

		call = .make_call(name, inputs, outputs)
		text = text.ReplaceSubstr(pos, selection.Size(),'\t\t' $ call $ '\r\n')

		return text
		}
	make_call(name, inputs, outputs)
		{
		assign = assign2 = ""
		if outputs.Size() is 1
			assign = outputs[0] $ ' = '
		else if outputs.Size() > 1
			{
			assign = "result = "
			for o in outputs
				assign2 $= '\t\t' $ o $ ' = result.' $ o $ '\r\n'
			}
		call = assign $ '.' $ name $ '(' $ inputs $ ')\r\n' $
			assign2
		return call[.. -2]
		}
	make_method(name, selection, inputs, outputs)
		{
		ret = ""
		if outputs.Size() is 1
			ret = "\t\treturn " $ outputs[0] $ '\r\n'
		else if outputs.Size() > 1
			{
			ret = "\t\treturn Object("
			for o in outputs
				ret $= ':' $ o $ ', '
			ret = ret[.. -2] $ ')\r\n'
			}

		return name $ '(' $ inputs $ ')\r\n' $
			'\t\t{\r\n' $
			.stripOuterCurlies(selection) $
			ret $
			'\t\t}'
		}
	stripOuterCurlies(selection)
		{
		opening_curly = '\A[ \t]*{(\r?\n)?'
		closing_curly = '[ \t]*}(\r?\n)?\Z'
		if selection =~ opening_curly and selection =~ closing_curly
			selection = selection.
				Replace(opening_curly, '', 1).
				Replace(closing_curly, '', 1)
		return selection
		}

	Inputs(text, pos, selection)
		{
		inside = ClassHelp.LocalsInputs(selection)

		mr = ClassHelp.MethodRange(text, pos)
		before = text[mr.from :: pos - mr.from]
		i = pos + selection.Size()
		after = text[i :: mr.to - i]
		outside = ClassHelp.LocalsAssigned(before $ '\n' $ after)

		return inside.Intersect(outside)
		}
	Outputs(text, pos, selection)
		{
		mr = ClassHelp.MethodRange(text, pos)
		before = text[mr.from :: pos - mr.from].AfterFirst(')') // skip parameters
		// TODO: handle parameters with nested parenthesis
		i = pos + selection.Size()
		after = text[i :: mr.to - i]
		outside = ClassHelp.LocalsInputs(before $ '\n' $ after)

		return ClassHelp.LocalsModified(selection).
			Intersect(outside)
		}
	}
