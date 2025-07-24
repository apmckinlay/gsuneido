// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	ContextMenu()
		{
		return #("Go To Definition\tF12")
		}
	On_Go_To_Definition(libview = false)
		{
		ref = .selectRef()
		if not .validRef(ref)
			return
		.SendToAddons(#GoToDef)
		if .isSuper(ref)
			return .goToSuper(ref, libview)

		// Attempt to run the control specific "Goto"
		ctrlHandled? = not (.Send(#Goto, ref) in (0, false))

		// If the "Goto" is not received by the control, run the standard "Goto"
		goto? = ctrlHandled? ? false : .Goto(ref, libview)

		// If neither "Goto" is successful, attempt "GotoLibView"
		if not ctrlHandled? and not goto?
			GotoLibView(ref, this, :libview)
		}
	goToSuper(ref, libview)
		{
		libs = false
		superClass = ClassHelp.SuperClass(.Get())
		ref = ref.Replace('^super', superClass.LeftTrim('_'))
		if superClass.Prefix?('_')
			libs = LibraryTags.GetTagFromName(.Send(#CurrentName)) isnt ''
				? [.table()]
				: Libraries()[.. Libraries().Find(.table())]
		GotoLibView(ref, this, :libview, :libs)
		return true
		}

	table()
		{ return .Send(#CurrentTable) }

	selectRef()
		{
		sel = .GetSelect()
		if sel.cpMin < sel.cpMax // there is a selection
			return .GetSelText()
		wc = .GetWordChars() $ ':'
		org = end = sel.cpMin
		end = .findUntilNotIn(end, wc)
		org = .findUntilNotIn(org, wc, dir: -1)
		afterFirstDot = org
		if .GetAt(org - 1) is '.'
			org = .findBeforeDot(org, wc)

		isJsOrCss = false
		if .GetAt(end) is '.'
			{
			afterEnd = .findUntilNotIn(end + 1, wc)
			if .Get()[end .. afterEnd] =~ "^\.(js|css)"
				{
				isJsOrCss = true
				end = afterEnd
				}
			}
		else if .Get()[afterFirstDot..end] =~ "^(js|css)"
			isJsOrCss = true

		orgToEndStr = .Get()[org .. end]

		if .skipFirstDot?(org, orgToEndStr, isJsOrCss)
			org = afterFirstDot
		.SetSelect(org, end - org)
		return .isSuper(orgToEndStr) ? orgToEndStr : .GetSelText().Trim(':')
		}
	findUntilNotIn(pos, wc, dir = 1)
		{
		offset = dir is 1 ? 0 : -1
		while .GetAt(pos + offset).In?(wc)
			pos += dir
		return pos
		}
	findBeforeDot(org, wc)
		{
		--org
		if .GetAt(org - 1).In?(wc)
			org = .findUntilNotIn(org, wc, dir: -1)
		else if .GetAt(org - 1).In?(')]`"\'')
			--org
		return org
		}
	isSuper(ref)
		{
		return ref.Prefix?('super.') or ref is 'super'
		}
	skipFirstDot?(org, orgToEndStr, isJsOrCss)
		{
		return .GetAt(org).Lower?() and not .isSuper(orgToEndStr) and not isJsOrCss
		}
	validRef(ref)
		{
		id = '[_a-zA-Z][_a-zA-Z0-9]*[?!]?'
		public = GlobalRegExForGoTo
		find = id $ ':' $ public $ '(:[0-9]+:?)?'
		return ref =~ '^(' $
			'(super\.?' $ id $ ')|' $
			'(\.?' $ id $ ')|' $
			'((' $ public $ '\.)?' $ public $ ')|' $
			'(' $ find $ ')' $
			')$'
		}

	// Called locally and by Addon_class_outline.ClassOutline_SelectItem
	Goto(ref, libview = false)
		{
		code = .Get()
		if libview is false
			libview = .currentLibView()
		.Send(#LibView_Goto)
		if .gotoObject?(ref, code, libview)
			return true

		if ref.Prefix?('.')
			{
			if not ClassHelp.Class?(code)
				ref = ref[1 ..] // e.g. for rules
			else
				return .gotoMethod?(ref, code)
			}

		.gotoDefinition(ref, libview)
		return true
		}

	currentLibView()
		{
		libview = .Send(#CurrentLibView)
		return libview is 0 ? false : libview
		}

	gotoObject?(object, code, libview)
		{
		if not object.Prefix?('.') or LibRecordType(code) isnt #object
			return false
		.gotoLibView(object, libview)
		.Defer(.SetFocus)
		return true
		}
	gotoLibView(object, libview)
		{
		GotoLibView(.Send(#CurrentName) $ object.Tr(':'),	libs: [.table()], :libview)
		}

	gotoMethod?(method, code)
		{
		method = method[1 ..]
		if false is gotoMethod? = .gotoMethodLine(method, code)
			.create_method(code, method, .GetSelect().cpMin)
		.Defer(.SetFocus)
		return gotoMethod?
		}

	gotoMethodLine(method, code)
		{
		if false is pos = ClassHelp.FindMethod(code, method)
			{
			if method[0].Upper?() and
				false isnt x = ClassHelp.FindBaseMethod(.table(), code, method)
				return .gotoBaseMethod(x.lib, x.name, method)
			if false is pos = ClassHelp.FindDotDeclarations(code, method)
				return false
			}
		.GotoLine(.LineFromPosition(pos))
		return true
		}

	gotoBaseMethod(lib, name, method)
		{
		if 0 isnt goto? = .Send(#GotoBaseMethod, lib, name, method)
			return goto?
		GotoLibView(name $ '.' $ method, libs: [lib])
		return true
		}

	create_method(code, method, pos)
		{
		if not .createMethod?(method)
			return
		pos = ClassHelp.AfterMethod(code, pos)
		code = code[.. pos] $
			'\t' $ method $ '()\r\n' $
			'\t\t{\r\n' $
			'\t\t}\r\n' $
			code[pos ..]
		line = .GetFirstVisibleLine()
		.PasteOverAll(code)
		.SetFirstVisibleLine(line)
		.SetSelect(pos + method.Size() + 2)
		}


	title: 'Go To Definition'
	createMethod?(method)
		{
		msg = 'Method ' $ Display(method) $ ' not found'
		if .GetReadOnly()
			{
			.AlertInfo(.title, msg)
			return false
			}
		return OkCancel(msg $ '\n\nCreate method?', .title, .Hwnd, MB.ICONQUESTION)
		}

	gotoDefinition(ref, libview)
		{
		libs = false
		if ref.Prefix?('_')
			{
			current = .Send(#CurrentName)
			if current isnt pureName = LibraryTags.RemoveTagFromName(current)
				{
				ref = pureName
				libs = [.table()]
				}
			else
				{
				// goto previous definition
				ref = ref[1 ..]
				libs = Libraries()
				libs = libs[.. libs.Find(.table())]
				}
			}
		path = .Send(#Goto_GetPath)
		GotoLibView(ref, .Parent, libs, :libview, path: path is 0 ? '' : path)
		}
	}
