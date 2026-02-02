// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: "ListFormView"
	New(form, .query = false, columns = false, title = '', columnsSaveName = '')
		{
		super(.makeControls(form, query, columns, title, columnsSaveName))
		.form = .FindControl('Form')
		.list = .FindControl('listFormVirtualList')
		.form.SetReadOnly(true)
		.Defer(.load_last_entry)
		}

	GetForm()
		{
		return .form
		}

	makeControls(form, query, columns, title, columnsSaveName)
		{
		f = Object('Record', form.Copy())
		f.name = "Form"
		return	Object('Vert',
				Object('VertSplit',
					Object('VirtualList', :query, :columns, startLast:,	:columnsSaveName,
						preventCustomExpand?:, filtersOnTop:, :title,
						name: 'listFormVirtualList'),
					Object('Scroll', f),
					splitSaveName: columnsSaveName))
		}

	Get()
		{
		return .form.Get()
		}

	GetQuery()
		{
		return .query
		}

	VirtualList_ItemSelected(x, source)
		{
		if source.Name is 'listFormVirtualList'
			.form.Set(x)
		}

	VirtualList_SetWhere()
		{
		.load_last_entry()
		}

	load_last_entry()
		{
		.list.On_VirtualListThumb_ArrowEnd()
		SetFocus(.list.GetGridHwnd())

		selectedRec = .list.GetSelectedRecord()
		if selectedRec is false
			selectedRec = Record()
		.form.Set(selectedRec)
		}
	}
