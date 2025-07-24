// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(itemVal = '', pathVal = '')
		{ return ToolDialog(0, Object(this, itemVal, pathVal)) }

	New(.itemVal, pathVal, ctrls = false)
		{
		if ctrls isnt false
			super(ctrls)
		.Tree = .FindControl(#TreeView)

		.Path = .FindControl(#path)
		.Path.Set(pathVal)

		.Item = .FindControl(#item)
		.Item.Set(itemVal)
		}

	Startup()
		{
		if '' is path = .Path.Get()
			return
		pathOb = path.Split('/')
		pathOb[0] = .Model.DisplayName(pathOb[0], pathOb[0])
		.Tree.GotoPath(pathOb)
		}

	New2()
		{
		super.New2()
		.Model = Construct(LibViewNewItemModel)
		}

	Getter_CurItem()
		{ return .Tree.GetSelectedItem() }

	Children?(item)
		{ return .Model.Children?(.Tree.GetParam(item)) }

	RootSelected?()
		{ return false }

	SelectTreeItem(olditem/*unused*/, newitem)	// From TreeView control
		{
		path = .Tree.Path(newitem)
		prefix = path.AfterLast('/').Trim()

		if .itemVal !~ '^Rule_|^Field_|^Table_'
			{
			item = .Item.Get().Replace('^Rule_|^Field_|^Table_', '')
			if prefix is 'Fields' or prefix is 'Rules' or prefix is 'Tables'
				item = prefix.Replace('s$', '') $ '_' $ item
			.Item.Set(item)
			}
		.Path.Set(path)
		}

	Valid?()
		{
		missing = Object()
		if .Item.Get().Blank?()
			missing.Add('Item')
		if .Path.Get().Blank?()
			missing.Add('Path')
		if not valid? = missing.Empty?()
			.AlertError(.Title, 'Required: ' $ missing.Join(', '))
		return valid?
		}
	}
