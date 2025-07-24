// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: ImageViewer
	CallClass(text, title = '', table = false, name = false)
		{
		Window(Object(this, text, table, name), :title, keep_placement: .Title)
		}

	New(.text, .table = false, .name = false)
		{
		}

	Controls()
		{
		return Object('Vert', .imageControl())
		}

	url: false
	imageControl()
		{
		return .name =~ '[.](png|jpg|jpeg|gif|svg|bmp|cur)$'
			? Object('Mshtml', style: .Style(.url = InMemory.Add(.text)))
			: .name =~ '[.](emf|wmf|ico)$'
				? Object('Image', .text)
				: ''
		}

	Style(url)
		{
		return 'body {
			background-image: url(' $ url.Replace(' ', '%20') $ ');
			background-repeat: no-repeat;
			}'
		}

	Destroy()
		{
		if .url isnt false
			InMemory.Remove(.url)
		super.Destroy()
		}
	}
