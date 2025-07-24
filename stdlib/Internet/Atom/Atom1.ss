// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.Id, .Title, .Updated)
		{
		.Entries = Object()
		}

	Entry(id, title, updated, content, author)
		{
		.Entries.Add([:id, :title, :updated, :content, :author])
		}

	DateFormat: 'yyyy-MM-ddTHH:mm:ssZ'

	ToString()
		{
		a = this
		return XmlBuilder(indent: 4).
			Instruct().
			feed(xmlns: "http://www.w3.org/2005/Atom")
				{
				.id(a.Id)
				.title(a.Title)
				.updated(a.Updated.GMTime().Format(a.DateFormat))
				for e in a.Entries
					.entry()
						{
						.id(e.id)
						.title(e.title)
						.updated(e.updated.GMTime().Format(a.DateFormat))
						.content(e.content, type: 'html')
						.author { .name(e.author) }
						}
				}.
			ToString()
		}
	}