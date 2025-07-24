// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Add Select'
	New(.queryState)
		{
		.indexList = this.FindControl('indexList')
		.presetList = this.FindControl('presetList')
		}
	Controls()
		{
		.indexes = .queryState.indexes
		.presets = .queryState.presets
		presetHorz = .presets.Empty?()
			? ''
			: Object('Horz'
				Object('ChooseList', .presets.Members(), selectFirst:,
					mandatory:, name: 'presetList')
				'Skip'
				#(Button 'Add Preset')
				)
		return Object('Vert'
			Object(#Static, .queryState.msg)
			'Skip'
			Object('Horz'
				Object('ChooseList', .indexes.Map(SelectPrompt), selectFirst:,
					mandatory:, name: 'indexList')
				'Skip'
				#(Button 'Add Select')
				),
			presetHorz
			'Skip'
			Object('HorzEqual', 'Fill',
				 #Skip #(Button 'Close')))
		}
	forceClose: false
	On_Add_Select()
		{
		if not .indexList.Valid?()
			return
		.queryState.filter = .getFilterFromChoice()
		.forceClose = true
		.Window.Result(.queryState)
		}

	On_Add_Preset()
		{
		if not .presetList.Valid?()
			return
		.queryState.filter = .presets[.presetList.Get()]
		.forceClose = true
		.Window.Result(.queryState)
		}

	getFilterFromChoice()
		{
		choice = .indexList.Get()
		field = .indexes.FindOne({ SelectPrompt(it) is choice })
		filter = [condition_field: field, check:]
		operation = .filterOperation(field)
		filter[field] = Object(:operation, value: '', value2: "")
		return filter
		}

	filterOperation(field)
		{
		dd = Datadict(field)
		operation = (dd.Base?(Field_date) and not dd.Base?(Field_num)) or
			dd.Base?(Field_number)
			? "greater than"
			: "equals"
		return operation
		}

	On_Close()
		{
		.forceClose = true
		.Window.Result(.queryState)
		}
	ConfirmDestroy() // stop alt+f4 to close window
		{
		return .forceClose
		}
	}