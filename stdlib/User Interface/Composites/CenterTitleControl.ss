// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
TitleNotesControl
	{
	Name: 'CenterTitle'
	New(@args)
		{
		super(@(args.Add(true, at: 'center')))
		}
	HelpButton_HelpPage()
		{
		book_option = 0
		if false is currentbook = Suneido.GetDefault('CurrentBook', false)
			return book_option

		option = String(.Controller).RemoveSuffix(`()`)
		if option is 'RecordControl' and .Controller.Member?('Controller')
			option = String(.Controller.Controller).RemoveSuffix(`()`)
		if false isnt rec = QueryFirstBookOption(currentbook, option)
			book_option = rec.path $ '/' $ rec.name
		return book_option
		}
	}
