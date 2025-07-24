// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Summarize"
	CallClass(ctrl, option = false)
		{
		ModalWindow(Object(this, ctrl, option), closeButton?: false)
		}

	New(.ctrl, option)
		{
		super(.layout(option))
		.data = .FindControl('Data')
		.grid = .FindControl('Grid')
		.countField = .FindControl('count')
		if .option isnt '' and false isnt userData = UserSettings.Get(.option)
			.data.Set(userData)
		}

	layout(option)
		{
		.initializeOption(option)
		.query = .ctrl.GetQuery()
		.numrows = 8
		exclude = .ctrl.Method?('GetExcludeSelectFields')
			? .ctrl.GetExcludeSelectFields()
			: #()
		columns = .ctrl.GetFields()
		.sf = SelectFields(columns, ExcludeLowerFields(columns, exclude))
		ctrls = Object(#((Static Field) (Static Function) (Static Result)))
		for (i = 0; i < .numrows; ++i)
			{
			fields = Object('ChooseList', .sf.Prompts().Sort!(), name: "fieldlist" $ i)
			ops = Object('ChooseList', #("max", "min", "total", "average"),
				name: "function" $ i)
			vals = Object('Field', readonly:, name: "val" $ i)
			ctrls.Add(Object(fields, ops, vals))
			}
		return Object('Vert',
			Object('Record',
				Object('Vert'
					Object('Horz', 'Fill', .presetsControl())
					Object("Grid", ctrls)
					)
				)
			'Skip',
			Object('Horz', 'Fill', #(Static Count), 'Skip',
				#(Field readonly:, name: "count"), name: 'Horz2'),
			'Skip',
			#(HorzEven
				Skip
				(Button, "Summarize")
				Skip
				(Button, "Clear")
				Skip
				(Button "Close")
				Skip),
			xstretch: 0)
		}

	initializeOption(option)
		{
		if option is false
			.option = .ctrl.GetDefault('Option', .ctrl.GetDefault('Title', ''))
		else
			.option = option
		.option = Opt(.option, ' Summarize')
		}

	presetsControl()
		{
		if .option is '' // nothing to save the preset under
			return 'Skip'
		return Object('Presets', .option)
		}

	On_Summarize()
		{
		if .data.Valid(forceCheck:) isnt true
			return

		data = .data.Get()
		data.Set_default("")
		sumFields = Object()
		sums = .buildSumFields(data, sumFields)
		query = QueryHelper.ExtendColumns(QueryStripSort(.query), .sf, sumFields)
		results = Query1("(" $ query $ ") summarize count" $ sums)
		if results is false
			results = Record(count:0)
		for i in .. .numrows
			{
			func = data['function' $ i]
			field_prompt = data['fieldlist' $ i]
			val = func is "" or field_prompt is ""
				? ""
				: results[func $ "_" $ .sf.PromptToField(field_prompt)]
			.grid["val" $ i].Set(val)
			}
		.countField.Set(results.count)
		}
	buildSumFields(data, sumFields)
		{
		sums = ""
		for i in .. .numrows
			{
			func = data['function' $ i]
			field_prompt = data['fieldlist' $ i]
			if func is "" or field_prompt is ""
				continue
			sumFields.Add(.sf.PromptToField(field_prompt))
			sums $= ", " $ func $ " " $ .sf.PromptToField(field_prompt)
			}
		return sums
		}
	On_Clear()
		{
		.data.Set(Record())
		}
	On_Close()
		{
		if .option isnt ''
			{
			data = .data.Get()
			UserSettings.Put(.option, data)
			}
		// should this be calling Result or just destroying the window?
		.Window.Result("")
		}
	}