// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	New(.columns, mandatory = false, width = false)
		{
		super(field: #('Field', readonly: true), :mandatory, :width)
		.ob = Object()
		}

	listcontrol: Controller
		{
		Name: 'EditManyListControl'
		Xmin: 400
		Ymin: 200
		list: false
		New(data, .columns, readonly = false)
			{
			super(.layout(data, .columns))
			.list = .Vert.List
			if readonly
				.list.SetReadOnly(readonly)
			}

		layout(data, columns)
			{
			return Object('Vert',
				Object('List', :data, :columns, noDragDrop:, noHeaderButtons:,
					defWidth: false, columnsSaveName: .Name))
			}
		List_WantNewRow()
			{
			return true
			}

		List_AfterEdit(col, row, data /*unused*/, valid?)
			{
			if not valid?
				.list.AddInvalidCell(col, row)
			else
				.list.RemoveInvalidCell(col, row)
			}

		OK()
			{
			data = .list.Get()
			for row in .list.Get()
				if Object?(row.list_invalid_cells) and
					not row.list_invalid_cells.Empty?()
					return false
			result = data.Map({ Record().Merge(it.Project(.columns)) }).Instantiate()
			return result
			}
		}

	Getter_DialogControl()
		{
		return Object(.listcontrol, .Get(), .columns, .GetReadOnly())
		}

	ProcessResults(result)
		{
		super.ProcessResults(.ob = result)
		}

	ob: #()
	Get()
		{
		newob = Object()
		for ob in .ob
			newob.Add(ob.Copy())
		return newob
		}

	Set(val)
		{
		if val is ''
			{
			.Field.Set('')
			.ob = Object()
			return
			}
		if not Object?(val)
			return

		.ob = val.Copy()
		setval = Object()
		for rec in .ob
			for col in .columns
				if rec.Member?(col)
					setval.Add(Prompt(col) $ ': ' $ rec[col])

		.Field.Set(setval.Join(', '))
		}
	}
