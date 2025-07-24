// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Export Data"
	CallClass(query, fields = '', selectFieldsWithJoins = false,
		excludeSelectFields = '', fileName = '', includeInternal = false,
		columns = false)
		{
		OkCancelModal(Object(this, query, fields, selectFieldsWithJoins,
			excludeSelectFields, fileName, includeInternal, columns), title: .Title)
		}

	New(.query, fields = '', .selectFieldsWithJoins = false, excludeSelectFields = '',
		fileName = '', includeInternal = false, .columns = false)
		{
		super(.controls(excludeSelectFields, includeInternal, fields))
		.Data.SetField('header', "Prompts")
		.Data.SetField('fields', fields)
		.Data.AddObserver(.RecordChanged)
		.setFileName(fileName)
		}
	controls(excludeSelectFields, includeInternal, fields)
		{
		.sf = false
		selectFieldsNoJoins = SelectFields(
			.columns is false or includeInternal is true
				? QueryColumns(.query) : .columns
			excludeSelectFields, joins: false, :includeInternal)
		prompts = selectFieldsNoJoins.Prompts()
		if .selectFieldsWithJoins is false
			.sf = selectFieldsNoJoins
		else
			{
			.sf = .selectFieldsWithJoins
			prompts.MergeUnion(fields)
			}
		if .selectFieldsWithJoins isnt false
			{
			// ensure that the user can't select prompts that are not in the query
			availablePrompts = QueryAvailableColumns(.query).Map!(SelectPrompt)
			prompts = prompts.Intersect(availablePrompts)
			}
		return .layout(prompts)
		}

	layout(prompts)
		{
		saveFileName = .SaveFileNameField()
		return Object('Record'
			Object('Vert'
				saveFileName,
				'Skip'
				#(Pair
					(Static 'Format')
					(ChooseList listField: export_formats, selectFirst:,
						width: 20, name: format)),
				'Skip',
				Object('Pair',
					#(Static 'Fields'),
					Object('ChooseTwoList', prompts, title: 'Fields',
						width: 20, name: 'fields')),
				'Skip',
				#(Pair (Static 'Header') (ChooseList #('None' 'Prompts' 'Fields')
					name: header))
				)
			)
		}

	SaveFileNameField()
		{
		return Object('Pair'
			#(Static "File name")
			Object('SaveFile', title: .Title, name: 'file'))
		}

	RecordChanged(member)
		{
		if member is 'format'
			.setFileName(.Data.Get().file.BeforeLast('.'))
		}

	setFileName(fileName)
		{
		if fileName is ''
			return
		data = .Data.Get()
		if false isnt (fn = .formatFunc(data))
			.Data.SetField('file', fileName $ '.' $ fn.Ext)
		}

	OK()
		{
		if .Data.Valid() isnt true
			 return false

		data = .Data.Get()
		if SaveFileControl.Method?('FileNameOnly?') and
			not SaveFileControl.FileNameOnly?(data.file)
			return false

		fn = .formatFunc(data)
		try
			{
			if fn is false
				throw 'unrecognized format: ' $ data.format
			fields = data.fields is '' ? false :
				.sf.PromptsToFields(data.fields)
			.Export(fn, data, fields)
			return true
			}
		catch (x)
			Alert("Error during Export: " $ x, 'Global Export',
				flags: MB.ICONERROR)
		return false
		}

	Export(fn, data, fields)
		{
		DoWithSelectedFileName(data.file)
			{ |fileName|
			fn(.query, fileName, fields, data.header, .sf.Fields)
			}
		}

	formatFunc(data)
		{
		fn = false
		for c in Plugins().Contributions("ImportExport", "export_formats")
			if c.name is data.format
				fn = Global(c.impl)

		return fn
		}
	}
