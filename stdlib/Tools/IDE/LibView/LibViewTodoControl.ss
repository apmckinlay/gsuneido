// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
// TODO allow sorting on "type" as well as line
Controller
	{
	Xmin: 10
	Ystretch: 1
	Name: 'LibViewTodo'
	New()
		{
		.todo = .FindControl(#todo)
		.todo.SetWordChars(
			"_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ?!:")
		}

	Controls: (TodoOutput readonly:, name: 'todo', height: 0 ystretch: .1)
	Set(data)
		{
		lib = data.GetDefault(#table, '')
		name = data.GetDefault(#name, '')
		text = data.GetDefault(#text, '')
		qcText = data.GetDefault(#qcText, '')
		if true is data.GetDefault(#group, false)
			{
			.folder(lib, ClassBrowserModel.Get_num(data.num),
				lines = Object()).SortWith!(.name)
			if lines.Size() > .MAXLINES
				lines.Add('... more...')
			}
		else
			{
			lines = Object()
			if '' isnt otherDefs = .getOtherDefinitions(name, lib)
				lines.Add(otherDefs)
			lines.Append(.extract(lib, name, text))
			}
		if data.Member?('qcText') and data.qcText isnt ""
			lines.Add(qcText)

		results = lines.Join('\n')
		results.Trim()
		if results[0] is '\n'
			results = results [1 ..]

		.todo.Set(results)
		}

	getOtherDefinitions(name, curLib)
		{
		libs = Object()
		for lib in LibraryTables().Remove(curLib)
			if not QueryEmpty?(lib, :name, group: -1)
				libs.Add(lib)
		msg = Opt(.otherDefinitionsWarning(name), libs.Join(', '))
		for m in LastContribution('Svc_TrialTags').Members()
			{
			trialName = name $ (name.Suffix?('__webgui') ? '_' : '__') $ m
			if name.Suffix?('_' $ m)
				trialName = name.Replace('(_|__)' $ m, '')
			if not QueryEmpty?(curLib, name: trialName, group: -1)
				msg = Opt(msg, '\n') $ 'WARNING: ' $ trialName $ ' also exists'
			}
		if BuiltinNames().BinarySearch?(name)
			msg $= "\nWARNING: Built-in on gSuneido"
		return msg
		}

	otherDefinitionsWarning(name)
		{
		return 'WARNING: ' $ name $ ' also exists in '
		}

	name(s)
		{
		return s.AfterFirst(':').BeforeFirst(':')
		}

	MAXLINES: 100
	folder(lib, num, lines)
		{
		if lines.Size() > .MAXLINES
			return lines
		QueryApply(lib, parent: num)
			{|x|
			if x.group is -1
				lines.Append(.extract(lib, x.name, x.text))
			else
				.folder(lib, x.num, lines) // recursive
			}
		return lines
		}

	extract(lib, name, text)
		{
		return text.Lines().Grep(`/[*/] *[A-Z][A-Z][A-Z]+`,
			{|i, line| lib $ ':' $ name $ ':' $ (i+1) $ ' ' $ line.Trim() })
		}

	Scintilla_DoubleClick()
		{
		line = .todo.GetLine()
		libview = .Send(#CurrentLibView)
		name = libview.CurrentName()
		if line =~ '^.+?:[[:alpha:]][_[:alnum:]]*[?!]?\>'
			{
			.todo.Home()
			.todo.CharRight()
			.todo.Recv(#On_Go_To_Definition, :libview)
			}
		else if line.Prefix?(.otherDefinitionsWarning(name))
			{
			lib = libview.CurrentTable()
			text = libview.Editor.Get()
			LibDiffOverriddenControl(name, lib, text)
			}
		else if line =~ `^WARNING: \S+ also exists$`
			{
			lib = libview.CurrentTable()
			trialName = line.RemovePrefix('WARNING: ').BeforeLast(' also exists')
			LibDiff(lib, name, lib, trialName)
			}
		}
	}
