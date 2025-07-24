// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New()
		{
		.items = []
		}

	channel: false
	Channel(@args)
		{
		if .channel isnt false
			throw "Rss2: only one channel allowed"
		.channel = new .channel_element
		.elements(.channel, args)
		}
	channel_element: class
		{
		Name: 'channel'
		Elements: (title, link, description,
			language, copyright, managingEditor, webMaster, pubDate,
			lastBuildDate, category, generator, docs, cloud, ttl, image,
			rating, textinput, skipDays, skipHours)
		Required: (title, link, description)
		}

	image: false
	Image(@args)
		{
		if .image isnt false
			throw "Rss2: only one image allowed"
		.image = new .image_element
		if not args.Member?(#link) and .channel isnt false
			args = args.Copy().Add(.channel.link, at: #link)
		.elements(.image, args)
		}
	image_element: class
		{
		Name: 'image'
		Elements: (url, title, link, width, height)
		Required: (url, title, link)
		}

	textInput: false
	TextInput(@args)
		{
		if .textInput isnt false
			throw "Rss2: only one textInput allowed"
		.textInput = new .textInput_element
		.elements(.textInput, args)
		}
	textInput_element: class
		{
		Name: 'textInput'
		Elements: (title, description, name, link)
		Required: (title, description, name, link)
		}

	AddItem(@args)
		{
		.items.Add(item = new .item_element)
		.elements(item, args)
		if not item.Member?(#title) and not item.Member?(#description)
			throw "Rss2: item must have either title or description"
		}
	item_element: class
		{
		Name: 'item'
		Elements: (title, link, description, author, category,
			comments, enclosure, guid, pubDate, source)
		Required: ()
		}

	ToString()
		{
		return '<?xml version="1.0" encoding="us-ascii"?>\n' $
			Xml('rss', version: '2.0')
				{
				'\n' $ .xml(.channel, .xml(.image) $ .xml(.textInput) $ .items_tostring())
				}
		}
	items_tostring()
		{
		s = ''
		for item in .items
			{
			if not item.Member?(#guid)
				item.guid = .guid()
			s $= .xml(item)
			}
		return s
		}
	guid()
		{
		return 'http://suneido.com/' $
			Display(Timestamp())[1 ..].Tr('.', '_')
		}

	elements(ob, args)
		{
		for m in args.Members()
			if ob.Elements.Has?(m)
				ob[m] = args[m]
			else
				throw 'Rss2: invalid element: ' $ ob.Name $ '/' $ m
		for m in ob.Required
			if not ob.Member?(m)
				throw 'Rss2: missing element: ' $ ob.Name $ '/' $ m
		}

	xml(ob, body = '')
		{
		if ob is false
			return ''
		s = ''
		for m in ob.Elements
			if ob.Member?(m)
				s $= Xml(m, XmlEntityEncode(ob[m])) $ '\n'
		return Xml(ob.Name, '\n' $ s $ body) $ '\n'
		}
	}