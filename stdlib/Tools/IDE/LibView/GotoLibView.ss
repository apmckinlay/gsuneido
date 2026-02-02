// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// handles library:Name and library:Name:line
class
	{
	// name can be name or name.method or lib:name or lib:name:line# (overrides libs arg)
	CallClass(name, scintilla = false, libs = false, line = false,
		libview = false, path = '', list = false)
		{
		path = path.Tr('()') // for libraries that are not in use
		name = name.Trim()
		if false is result = .format_variables(name, libs, line, path, list)
			return
		name = result.name
		line = result.line
		list = result.list
		method = result.method
		if false is result = .process(name, line, path, list, scintilla)
			return

		new_item = result.new_item
		path = result.path
		if libview is false
			libview = .get_libview()
		if false isnt new_item and
			false is (path = .output_new_item(new_item, libview))
			return

		libview.GotoPathLine(path.Tr('()'), line, skipFolder?:)
		if method isnt ''
			libview.GotoMethodLine(method)
		}
	format_variables(name, libs, line, path, list)
		{
		fromFindOutput = false
		if name.Has?(':')
			{
			parts = name.Split(':')
			if parts.Size() is 1
				return false
			fromFindOutput = true
			libs = Object(parts[0])
			name = parts[1]
			if parts.Size() > 2
				line = Number(parts[2])
			}
		if false is result = .split(name)
			return false
		result.list = .generateList(list, result.name, libs, line, fromFindOutput)
		if ((libs is false or libs.Has?('Builtin')) and false isnt .hasBuiltIn(name))
			result.list.Add('Go To Documentation')
		.extraAddRemove(result.list, path, fromFindOutput, result.name)
		result.line = line is false ? 1 : line
		return result
		}
	split(name)
		{
		hasJsOrCss = name =~ '\.(js|css)$'
		method = hasJsOrCss ? '' : name.AfterFirst('.')
		name = hasJsOrCss ? name : name.BeforeFirst('.')
		if not .identifier?(name) and not hasJsOrCss
			return false
		method = method.RemovePrefix(name $ '_')
		return Object(:method, :name)
		}
	identifier?(s)
		{
		// same as Strings.Identifier? except it also allows '?' within the name
		// (for __protect rules on names ending in '?')
		return s =~ `\A[[:alpha:]][a-zA-Z0-9_?]*!?\Z`
		}
	generateList(list, name, libs, line, fromFindOutput)
		{
		exact = fromFindOutput or line isnt false
		if list is false
			list = Gotofind(name, libs, :exact)
		return list.Copy()
		}
	extraAddRemove(list, path, fromFindOutput, name)
		{
		// don't include the record you are on
		if path isnt ''
			list.Remove(path)
		if not fromFindOutput and (.isTableName?(name) or .isViewName?(name))
			list.AddUnique('schema:' $ name)
		}
	isTableName?(table)
		{
		return not QueryEmpty?('tables', :table)
		}
	isViewName?(view_name)
		{
		return not QueryEmpty?('views', :view_name)
		}
	process(name, line, path, list, scintilla)
		{
		new_item = false
		if list.Empty?()
			{
			if .goto_built_in?(name) is true
				return false
			if false is (new_item = .new_item(name, line, path))
				{ Beep(); return false }
			}
		else if list.Size() is 1 or scintilla is false
			path = list[0]
		else // multiple matches
			{
			pt = scintilla.PointFromPosition(scintilla.GetSelectionStart())
			ClientToScreen(scintilla.Hwnd, pt)
			pt.y += 15/*=height*/
			i = ContextMenu(list).Show(scintilla.Window.Hwnd, pt.x, pt.y, left:)
			if i <= 0
				return false
			path = list[i - 1]
			}
		if path.Prefix?('schema:')
			{
			SchemaView.Goto(path.AfterFirst(':'))
			return false
			}
		if path is 'Go To Documentation'
			if .goto_built_in?(name) is true
				return false
		return Object(:new_item, :path)
		}
	goto_built_in?(name)
		{
		if false isnt x = .hasBuiltIn(name)
			{
			GotoUserManual(x.path $ '/' $ x.name)
			return true
			}
		return false
		}
	hasBuiltIn(name)
		{
		if BuiltinNames().Has?(name) and
			(false isnt x = QueryFirst('suneidoc where name is ' $ Display(name) $
				' sort path'))
			return x
		return false
		}
	new_item(name, line, path)
		{
		// check to see if definition exists in unused libraries
		other_list = Gotofind(name, LibraryTables().Difference(Libraries()),
			exact: line isnt false)
		msg = 'Definition for ' $ name
		msg $= other_list.Size() > 0
			? ' already exists in ' $ other_list[0].Split('/')[0]
			:  ' does not exist'
		return LibViewNewItemControl(name, msg, path)
		}
	get_libview()
		{
		return GotoPersistentWindow('LibViewControl', LibViewControl)
		}
	output_new_item(new_item, libview)
		{
		name = new_item.item
		if false is libview.Explorer.GotoPath(new_item.path)
			{
			Alert('Cannot find path specified, "' $ name $ '" not created',
				title: "Goto LibView", flags: MB.ICONINFORMATION)
			return false
			}
		return KeyException.TryCatch()
			{
			libview.Explorer.NewItem(false, :name, text: .text(name))
			return new_item.path $ '/' $ new_item.item
			}
		}
	text(name) // TODO use LibView New plugin templates
		{
		if name.Suffix?("Test")
			return .Test_text()
		else if name.Prefix?("Rule")
			return .Rule_text()
		else
			return ""
		}
	Test_text()
		{
		return "Test
	{
	Test_one()
		{
		}
	}"
		}
	Rule_text()
		{
		return "function ()
	{
	}"
		}
	}
