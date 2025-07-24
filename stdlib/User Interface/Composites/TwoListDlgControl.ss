// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xmin: 450
	Ymin: 250
	Xstretch: 1
	Ystretch: 1
	Name: "ChooseTwoList"
	New(list, initial_list, title = "", mandatory_list = #(), extra_buttons = #(),
		.control = 'TwoList', displayMemberName = '', noSort = false, delimiter = ',')
		{
		super(.layout(list, initial_list, mandatory_list, extra_buttons,
			displayMemberName, noSort, delimiter))
		.Title = title
		.TwoList = .FindControl('TwoList')
		}

	layout(list, initial_list, mandatory_list, extra_buttons, displayMemberName,
		noSort = false, delimiter = ',')
		{
		buttons = Object('Horz', #(LinkButton, 'Help'), #(Fill min: 10))
		buttons.Append(extra_buttons)
		return Object('Vert'
			Object(.control, :list, :initial_list, :mandatory_list, :displayMemberName,
				name: 'TwoList', :noSort, :delimiter),
			#(Skip 5),
			buttons)
		}

	On_Help()
		{
		return OpenBook(Suneido.CurrentBook $ 'Help',
			Object(path: "/General/Reference/", name: "Field Chooser"))
		}

	OK()
		{
		if .TwoList.Member?('Valid') and false is .TwoList.Valid()
			return false
		return .TwoList.Get()
		}
	}
