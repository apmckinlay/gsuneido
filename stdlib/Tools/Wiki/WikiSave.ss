// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	maxSize: 64000
	CallClass(page, entity_body, table = 'wiki')
		{
		err = ""
		values = Url.SplitQuery(entity_body)
		if values.Member?("Save")
			{
			text = String(values.text)
			if values.editmode is 'append'
				if false isnt x = Query1(table, name: page)
					text = x.text $ "\n----\n" $ text

			err = .Save(text, table, page)
			}
		WikiUnlock(page)
		if err isnt ""
			return '<html><body>' $ err $ '</body></html>'
		return '<html><head>
			<meta http-equiv="Refresh" content="0;URL=Wiki?' $ page $ '">
			</head><body>Thank you</body></html>'
		}
	Save(text, table, page)
		{
		if text.Size() > .maxSize
			return "ERROR: page too large (" $ text.Size() $ " > " $ .maxSize $ ")"
		if '' isnt valid = WikiTitleValid?(page)
			return valid
		text = text.Tr('\r').Trim()
		if text is "" and not .wikiNotesPage?(page)
			QueryDelete(table, [name: page])
		else
			{
			Transaction(update:)
				{ |t|
				if false is x = t.Query1(table, name: page)
					{
					num = NextTableNum(table, t)
					t.QueryOutput(table,
						[:num, name: page, :text, created: Date()])
					}
				else
					{
					x.text = text
					x.edited = Date()
					x.Update()
					}
				}
			}
		.removeOrphans(text)
		ServerSuneido.Set('wiki_regenerate', true)
		return ""
		}

	wikiNotesPage?(page)
		{
		for book in BookTables()
			{
			if false is notes = OptContribution(book $ 'WikiNotes', false)
				continue
			if page.Prefix?(notes.PagePrefix)
				return true
			}
		return false
		}

	removeOrphans(text)
		{
		updateOrphanList = false
		for name in WikiOrphans.ListOrphans()
			{
			if text.FindRx('\<' $ name $ '\>') < text.Size()
				{
				QueryApply1('wiki', :name)
					{
					updateOrphanList =true
					it.orphaned = false
					it.Update()
					}
				}
			}
		if updateOrphanList is true
			WikiOrphans.ResetList()
		}
	}
