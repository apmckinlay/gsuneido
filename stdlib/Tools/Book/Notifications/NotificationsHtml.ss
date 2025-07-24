// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(from, to)
		{
		data = ServerEval("BookNotification.GetNewEntries", Suneido.User, from, to)
		if Suneido.Member?('ShowNewEntries')
			Suneido.ShowNewEntries.status = true
		return .BuildHtml(data)
		}
	BuildHtml(data)
		{
		page = data.Join('\n')
		font = Sys.GUI?()
			? StdFonts.GetCSSFont(sizeFactor: 0.9)
			: "font-size: 14px;"
		html = Xml('html',
			Xml('head',
				Xml('style',
					'\nh3 { margin-bottom: 0em; margin-bottom: 0em; }' $
					'\nbody { ' $ font $ ' }\n',
					type: "text/css")
				) $ '\n' $
			Xml('body', page))
		return html
		}
	}
