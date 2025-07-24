// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(name, source, text = '', t = false, reporterMode = 'simple')
		{
		c = LastContribution('Reporter')
		if '' is book = c.ReporterBook
			return

		DoWithTran(t, update:)
			{ |t|
			if not Object?(source) or false is location = .getBookLocation(source)
				return
			loc = '/' $ location $ c.GetPath(reporterMode)
			if false isnt t.Query1(book, path: loc, name: name.Tr('~'))
				return

			rec = Record(name: name.Tr('~'),
				order: name.Prefix?('~') ? 10 : 20, /*= 10 standard, 20 user default*/
				path: loc,
				text: Display(text)
				plugin: false,
				num: NextTableNum(book, t))

			t.QueryOutput(book, rec)
			}
		}

	getBookLocation(source)
		{
		return Reporter_table.BookLocation(source)
		}
	}
