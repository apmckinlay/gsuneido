// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(helpbook)
		{
		index = 'index_' $ helpbook
		builder = Ftsearch.Create()
		.Foreach_record(helpbook)
			{|x, text|
			builder.Add(x.num, x.path $ '/' $ x.name, text)
			}
		idx = builder.Pack()
		PutFile(index, idx)
		Suneido[index] = idx
		return ''
		}

	Foreach_record(helpbook, block, extraWhere = '')
		{
		_table = helpbook // needed by suneidoc Asup
		xr = new XmlReader
		handler = new .handler
		handler.SkipClasses = .skipClasses(.htmlWrapPrefix(helpbook))
		xr.SetContentHandler(handler)
		.queryApply(helpbook, extraWhere)
			{|x|
			text = .process1(x, handler, xr)
			block(x, text)
			}
		}
	htmlWrapPrefix(helpbook) // overridden by test
		{
		return HtmlWrapPrefix(helpbook)
		}
	queryApply(helpbook, extraWhere, block) // overridden by test
		{
		QueryApply(helpbook $ ' where path !~ "^/res\>" ' $ extraWhere, block)
		}
	skipClasses(styles)
		{
		// This attempts to figure out which classes are skipped (display: none)
		// Primarily it looks for .myclass { ... display: none
		// It doesn't handle if the class and display: none are on different lines.
		// It doesn't handle if a class is display: none but later changed.
		// It tries to ignore media="print" styles.
		skip = false
		ob = Object()
		for line in styles.Lines()
			{
			if line =~ `^<style.*media="print"`
				skip = true
			else if line.Has?('</style>')
				skip = false
			else if not skip and line.Has?("display: none")
				{
				if line !~ "^[.]\w+ {"
					throw "ERROR: IndexHelp: unhandled: " $ line
				cls = line[1..].Extract("^\w+")
				ob[cls] = true
				}
			}
		ob.Delete(#showhide)
		return ob
		}

	process1(x, handler, xr)
		{
		_path = x.path // needed by suneidoc Asup
		_name = x.name // needed by suneidoc Asup
		if not x.text.Prefix?('<')
			x.text = x.name // just index name
		handler.S = ''
		try
			{
			xr.Parse(Asup(x.text,
				Object(GetHelpPage: Name(OptContribution('GetHelpPage', GetHelpPage)))))
			handler.Check()
			}
		catch (e)
			throw 'ERROR: IndexHelp:' $ x.path $ "/" $ x.name $ " - " $ e
		return handler.S
		}

	handler: XmlContentHandler
		{
		// This attempts to skip hidden content
		// that is inside head, script, style, or SkipClasses.
		// It does not handle visible parts nested inside hidden parts.
		// It assumes tags are ended and nested correctly.
		skip: 0
		skipTag: false
		StartElement(qname, atts)
			{
			if .skip is 0
				{
				if qname in (#head, #script, #style)
					{
					.skip++
					.skipTag = qname
					}
				if false isnt cls = atts.GetDefault("class", false)
					if .SkipClasses.Member?(cls)
						{
						.skip++
						.skipTag = qname
						}
				}
			else if qname is .skipTag
				.skip++
			}

		EndElement(qname)
			{
			if .skipTag is qname
				{
				--.skip
				if .skip is 0
					.skipTag = false
				else if .skip < 0
					throw "ERROR: IndexHelp: extra closing tag: " $ qname
				}

			}

		Characters(s)
			{
			if .skip > 0
				return
			s = s.Trim()
			if s isnt '' and not .entity?(s)
				.S $= s $ '\n'
			}

		entity?(s)
			{
			return s =~ '^&\w+;$'
			}

		Check()
			{
			if .skip isnt 0
				throw "ERROR: IndexHelp: mismatched: " $ .skipTag $ ' ' $ .skip
			}
		}
	}
