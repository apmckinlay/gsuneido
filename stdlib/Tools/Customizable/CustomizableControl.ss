// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: 'Customizable'
	TabName: false
	New(table = false, name = false, tabName = false)
		{
		super(.layout(table, name, tabName))
		.horz = .FindControl('Horz')
		}

	layout(table, name, tabName)
		{
		.user = CustomizeExpandControl.LayoutName is tabName ? Suneido.User : ''
		ctrl = _parent.Parent
		layout = .buildTabLayout(table, name, tabName)
		return _parent.GetChildren().Size() <= 1 and
			ctrl.Base?(WndPaneControl) and ctrl.Member?('Parent') and
			ctrl.Parent.Base?(TabsControl)
			? Object('Scroll', Object('Horz', layout, xstretch: 1, ystretch: 1))
			: Object('Horz', layout)
		}

	buildTabLayout(table, name, tabName)
		{
		// Needed by .Send, but .Controller is initailized in super which is too late
		.Controller = _parent.Base?(Controller) ? _parent : _parent.Controller
		if table is false
			table = QueryGetTable(.Send('GetQuery'), orview:)
		if name is false and 0 is name = .Send("GetCustomizableName")
			name = false
		if 0 is expandInfo = .Send('Customizable_ExpandInfo')
			expandInfo = Object(availableFields: false, defaultLayout: '')
		if 0 is customKey = .Send('GetAccessCustomKey')
			customKey = ''
		.c = Customizable(table, name, user: .user,
			availableFields: expandInfo.availableFields,
			defaultLayout: expandInfo.defaultLayout, :customKey)
		return .tabLayout(tabName)
		}

	tabLayout(tabName)
		{
		.TabName = tabName
		if tabName is false and 0 is tabName = .Send('TabGetPath')
			tabName = 'Header'
		limitHeight = CustomizeExpandControl.LayoutName is tabName
		return .c.CustomTableTab?(tabName)
			? .c.Table(tabName)
			: .c.Form(tabName, :limitHeight)
		}

	horz: false
	Recalc()
		{
		if .horz isnt false and .horz.GetChildren().NotEmpty?() and
			.horz.GetChildren()[0].GetChildren().NotEmpty?() and
			.Parent.Base?(FormControl)
			{
			ctrl = .horz.GetChildren()[0].GetChildren()[0]
			if .Parent.Left > (ctrl.X + ctrl.Left)
				.X = .Parent.Left - (ctrl.X + ctrl.Left)
			}
		super.Recalc()
		}
	}
