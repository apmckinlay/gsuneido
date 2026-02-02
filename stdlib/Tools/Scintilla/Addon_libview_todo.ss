// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
CodeViewAddon
	{
	Name: 	LibViewTodo
	Inject: bottomLeft
	Controls(ob)
		{ ob.Add(#LibViewTodo) }

	Addon_RedirMethods()
		{ return #(CheckCode_QualityChanged) }

	Init_qctext: ''
	InitialSet()
		{ .set(.Init_qctext) }

	prevQcText: ''
	set(qcText)
		{
		data = Object(
			table: .Controller.Table,
			name: .Controller.RecName,
			num: .Controller.Num,
			group: .Controller.Group,
			text: .Get())
		if data.group is true
			qcText = ''
		.prevQcText = data.qcText = qcText
		.AddonControl.Set(data)
		}

	CheckCode_QualityChanged(checks)
		{ .setTodo(checks.warningText) }

	Invalidate()
		{ .setTodo(.prevQcText) }

	AfterSave()
		{ .setTodo(.prevQcText) }

	setTodo(qcText)
		{
		.set(qcText)
		if .showWarnings?(qcText)
			.SubSplit.SetSplit(.SubSplit.GetSplit())
		}

	prevWarningOb: #()
	showWarnings?(warnings)
		{
		if '' is warnings = warnings.Replace('[0-9]', '')
			return false

		warningOb = warnings.Lines().Map!(#Trim).Filter({ it isnt '' })
		show? = .todoHeight() <= 0 and warningOb.Any?({ not .prevWarningOb.Has?(it)})
		.prevWarningOb = warningOb
		return show?
		}

	todoHeight()
		{ return .AddonControl.GetChild().GetClientRect().GetHeight() }
	}
