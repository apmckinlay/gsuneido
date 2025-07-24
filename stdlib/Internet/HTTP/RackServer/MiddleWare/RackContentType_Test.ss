// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		result = ['200 OK', #(), ""]
		app = {|env/*unused*/| result }
		mw = new RackContentType(app)
		Assert(mw(#(path: 'fred'))
			is: result)
		Assert(mw(#(path: 'favicon.ico'))
			is: ['200 OK', #(Content_Type: "image/x-icon"), ""])
		Assert(mw(#(path: 'script.js'))
			is: ['200 OK', #(Content_Type: "text/javascript"), ""])

		result = ['200 OK', #(Content_Type: 'text/plain'), ""]
		Assert(mw(#(path: 'fred'))
			is: result)
		Assert(mw(#(path: 'script.js'))
			is: result)
		}
	}