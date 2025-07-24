// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.

// TODO multiple author, category, contributor, link
// TODO author/contributor uri, email, etc.
// TODO link - other types
// TODO content - other types

class
	{
	New()
		{
		.entries = []
		}

	feed: false
	Feed(@args)
		{
		if .feed isnt false
			throw "Atom: only one feed allowed"
		.feed = new .feed_element
		.elements(.feed, args)
		}
	feed_element: class
		{
		Name: 'feed'
		Required: (title, link)
		Elements: (title, link, author, id, updated,
			subtitle, rights, contributor, generator, post, category,
			icon, logo)
		}

	AddEntry(@args)
		{
		.entries.Add(entry = new .entry_element)
		.elements(entry, args)
		if not .feed.Member?(#author) and not entry.Member?(#author)
			throw "Atom: author must be specified on feed or entry"
		}
	entry_element: class
		{
		Name: 'entry'
		Required: (title, link)
		Elements: (title, link, author, id, updated,
			summary, published, contributor, content, edit, category, rights,
			source)
		}

	ToString()
		{
		return '<?xml version="1.0" encoding="us-ascii"?>\n' $
			.xml(.feed, .entries_tostring())

		}
	DateFormat: 'yyyy-MM-ddTHH:mm:ssZ'
	entries_tostring()
		{
		s = ''
		for entry in .entries
			s $= .xml(entry)
		return s
		}

	elements(ob, args)
		{
		for m in args.Members()
			if ob.Elements.Has?(m)
				ob[m] = args[m]
			else
				throw 'Atom: invalid element: ' $ ob.Name $ '/' $ m
		for m in ob.Required
			if not ob.Member?(m)
				throw 'Atom: missing element: ' $ ob.Name $ '/' $ m
		}

	xml(ob, body = '')
		{
		if ob is false
			return ''
		s = ''
		if not ob.Member?(#updated)
			ob.updated = Date().GMTime().Format(.DateFormat)
		if not ob.Member?(#id)
			ob.id = 'urn:uuid:' $ UuidString()
		for m in ob.Elements
			if ob.Member?(m)
				{
				s $= '\t'
				switch (m)
					{
				case 'link' :
					s $= Xml(m, href: ob[m], rel: 'alternate', type: 'text/html')
				case 'author', 'contributor' :
					s $= Xml(m, Xml('name', XmlEntityEncode(ob[m])))
				case 'updated' :
					s $= Xml(m, ob[m])
				case 'content' :
					content = ob[m]
					type = 'text'
					if content.Prefix?('<div>')
						type = 'xhtml'
					else if content.Prefix?('<')
						{
						type = 'html'
						content = XmlEntityEncode(content)
						}
					s $= Xml(m, content, :type)
				default :
					s $= Xml(m, XmlEntityEncode(ob[m]))
					}
				s $= '\n'
				}
		return Xml(ob.Name, '\n' $ s $ body) $ '\n'
		}
	}