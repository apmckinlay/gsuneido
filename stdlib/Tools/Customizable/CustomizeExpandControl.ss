// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xmin: 500
	Ymin: 400
	Title: 'Customize Expand'

	LayoutName: 'CustomizableExpand'

	CallClass(query, fields, defaultExpandLayout, customKey)
		{
		return OkCancel(Object(this, query, fields, defaultExpandLayout, customKey),
			.Title)
		}

	New(query, .fields, defaultExpandLayout, customKey)
		{
		.list = .FindControl('ListBox')
		.editor = .FindControl("Editor")
		.c = Customizable(QueryGetTable(query, orview:), user: Suneido.User,
			defaultLayout: defaultExpandLayout, :customKey)
		.sf = SelectFields(fields, headerSelectPrompt:, joins: false)
		for prompt in .sf.Fields.Members()
			.list.AddItem(prompt)
		.editor.Set(.c.Layout(.sf, .LayoutName))
		.editor.Dirty?(false)
		}

	Controls()
		{
		ctrls = Object('Vert',
			Object('TitleNotes', .Title, name: 'title')
			Object('HorzSplit',
				#(Vert
					#(Border (Static 'Available') 3)
					#(ListBox sort:))
				Object(#Vert
					#(Border (Static 'Expand') 3)
					Object('ScintillaAddons', xstretch: 10, margin: 5, Addon_detab:,
						Addon_highlight_customFields: Object(field_dicts: .fields))))
			#Skip)
		buttons = Object(#Horz
			#(Button 'Add Field to Layout') #Skip #(Button 'Restore Default Layout'))
		if AccessPermissions(Customizable.PermissionOption()) is true
			buttons.Add(
				#Skip #Skip #(Button 'Save As Default') #(Static '(for all users)'))
		buttons.Add(#Fill, #Skip, #OkCancel)
		ctrls.Add(buttons)
		return ctrls
		}

	HelpButton_HelpPage()
		{
		return "/General/Reference/Customization Options"
		}

	ListBoxDoubleClick(i)
		{
		.addField(i)
		}

	On_Add_Field_to_Layout()
		{
		if false is i = .getCurrentlySelectedField()
			return
		.addField(i)
		}

	getCurrentlySelectedField()
		{
		if .list.GetCurSel() is -1
			{
			.AlertInfo(.Title, "Please select a field")
			return false
			}
		return .list.GetCurSel()
		}

	addField(i)
		{
		prompt = .list.GetText(i)
		if .checkFieldUsed(prompt)
			{
			.AlertError(.Title, "This field is already selected on the Expand")
			return
			}
		.editor.Paste(prompt $ ' ')
		}

	checkFieldUsed(prompt)
		{
		sf = SelectFields(.fields, headerSelectPrompt:, joins: false, includeInternal:)
		field = sf.PromptToField(prompt)
		fields = sf.FormulaFields(.editor.Get())
		return fields.Has?(field)
		}

	On_Restore_Default_Layout()
		{
		layout = .c.DefaultLayout(.sf, .LayoutName)
		.editor.Set(layout)
		.editor.Dirty?(true)
		}

	On_Save_As_Default()
		{
		if true is .saveLayout(asDefault:)
			.AlertInfo(.Title, 'Layout saved as default for all users')
		}

	On_OK()
		{
		.Send('On_OK')
		}

	maxLines: 15
	OK()
		{
		return .saveLayout()
		}

	saveLayout(asDefault = false)
		{
		curLayout = .editor.Get().Trim()
		if .hasDuplicateField?(curLayout)
			return false
		if curLayout.Lines().Size() >= .maxLines
			{
			.AlertInfo(.Title, 'You cannot have more than ' $ .maxLines $ ' lines')
			return false
			}
		if .editor.Dirty?() or asDefault is true
			.c.SaveLayout(curLayout, .sf, .LayoutName, :asDefault, limitHeight:)
		return true
		}

	hasDuplicateField?(layout)
		{
		dups = .sf.FormulaFields(layout).DuplicateValues()
		if not dups.Empty?()
			{
			.AlertInfo(.Title,
				'You cannot add the same field (' $
					dups.Map(SelectFields.GetFieldPrompt).Join(', ') $
					') to the layout more than once.')
			return true
			}
		return false
		}
	}
