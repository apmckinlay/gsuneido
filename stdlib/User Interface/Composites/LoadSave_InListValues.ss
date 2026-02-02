// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Load(hwnd, fieldSaveName)
		{
		if false is name = ParamsChooseListOpenDialog(hwnd, fieldSaveName)
			return false
		if false is valRec = Query1('params', report: name)
			return false
		return valRec.params
		}

	Save(hwnd, fieldSaveName, values)
		{
		if false is name =  (.save_dialog)(hwnd, fieldSaveName)
			return
		if not QueryEmpty?('params', report: "InlistValues - " $ name)
			QueryDo('update params
				where report is ' $ Display("InlistValues - " $ name) $
				' set params = ' $ Display(values) $
				', report_options = ' $ Display(fieldSaveName))
		else
			QueryOutput('params',
				Record(user: Suneido.User, report: "InlistValues - " $ name,
				params: values, report_options: fieldSaveName))
		}

	save_dialog: Controller
		{
		Title: 'Save List'
		CallClass(hwnd, fieldSaveName)
			{
			return OkCancel(Object(this, fieldSaveName), .Title, hwnd)
			}
		New(field)
			{
			.saveUnderField = field
			.field = .Vert.Pair.Field
			.list = .Vert.ListBox
			prefix = 'InlistValues - '
			QueryApply("params where report =~ " $ Display(prefix) $
				" and report_options is " $ Display(field) $ " sort report")
				{|x|
				.list.AddItem(x.report[prefix.Size()..])
				}
			}
		Controls()
			{
			return #(Vert
				(Pair (Static 'Name')(Field))
				(ListBox xstretch: 1))
			}
		ListBoxSelect(sel)
			{
			.field.Set(.list.GetText(sel))
			}
		OK()
			{
			name = .field.Get()
			nameLimit = 100
			if name is "" or name.Size() > nameLimit
				{
				.AlertInfo('Save List',
					'Please enter or select a name (max: ' $ nameLimit $ ' characters)')
				return false
				}
			overwriteMsg = 'Saved List Name Already Exists. Do you want to overwrite it?'
			if not QueryEmpty?("params
				where report is " $ Display('InlistValues - ' $ name) $ " and " $
				"report_options isnt " $ Display(.saveUnderField)) and
				not OkCancel(overwriteMsg, 'Save List')
					return false
			return name
			}
		}
	}
