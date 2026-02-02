// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(md, noIndent? = false, addons = #())
		{
		_mdAddons = .initAddons(addons)
		for addon in _mdAddons
			if addon.Method?(#PreParse)
				md = (addon.PreParse)(md)

		parsed = String?(md) ? MarkdownParser(md, :addons) : md
		html = .Convert(parsed, :noIndent?)

		for addon in _mdAddons
			if addon.Method?(#AfterHtml)
				html = (addon.AfterHtml)(html)
		return html
		}

	initAddons(addons)
		{
		res = Object()
		for addon in addons
			res.Add(Construct(addon))
		return res
		}

	Convert(parsed, noIndent? = false)
		{
		writer = new Md_HtmlWriter(:noIndent?)
		.convertContainer(writer, parsed)
		return writer.Get()
		}

	convertContainer(writer, container, tight? = false, _mdAddons = #())
		{
		for addon in mdAddons
			addon.PreprocessContainer(container)
		container.ForEachBlockItem()
			{ |item|
			.convertItem(writer, item, :tight?)
			}
		}

	convertItem(writer, item, tight? = false)
		{
		if .callAddons(Object(writer, item), #ConvertToHtml)
			return
		switch
			{
		case item.Base?(Md_ATXheadings):
			writer.Add('h' $ item.Level, .ConvertInline(item.ParsedInline))
		case item.Base?(Md_Code):
			writer.Add('pre',
				'<code' $ Opt(' class="language-', item.Info.BeforeFirst(' '), '"') $
					'>' $ Opt(.Encode(item.Codes.Join('\n')), '\n') $ '</code>')
		case item.Base?(Md_Html):
			writer.AddPure('\n' $ item.Html)
		case item.Base?(Md_Paragraph):
			.convertParagraph(writer, item, tight?)
		case item.Base?(Md_ThematicBreak):
			writer.Add('hr', false)
		case item.Base?(Md_BlockQuote):
			writer.AddWithBlock('blockquote')
				{
				.convertContainer(writer, item)
				}
		case item.Base?(Md_List):
			attr = item.Start not in (false, 1) ? Object(start: item.Start) : #()
			writer.AddWithBlock(item.Type, :attr)
				{
				.convertContainer(writer, item, tight?: not item.Loose?)
				}
		case item.Base?(Md_ListItem):
			writer.AddWithBlock('li', noNewline?: .noNewline?(tight?, item))
				{
				.convertContainer(writer, item, :tight?)
				}
			}
		}

	convertParagraph(writer, item, tight?)
		{
		if item.Inline.Blank?()
			Nothing()
		else if item.HeadingLevel isnt false
			writer.Add('h' $ item.HeadingLevel, .ConvertInline(item.ParsedInline))
		else if tight?
			writer.AddPure((writer.Suffix?('<li>') ? '' : '\n') $
				.ConvertInline(item.ParsedInline))
		else
			writer.Add('p', .ConvertInline(item.ParsedInline))
		}

	noNewline?(tight?, container)
		{
		last = false
		container.ForEachBlockItem()
			{
			last = it
			}
		return last is false or tight? is true and last.Base?(Md_Paragraph)
		}

	callAddons(args, call, _mdAddons = #())
		{
		for addon in mdAddons
			if addon.Method?(call) and (addon[call])(@args) is true
				return true
		return false
		}

	ConvertInline(parsedInline)
		{
		s = ''
		for item in parsedInline
			{
			if .callAddons(Object(item), #PreInline)
				continue
			switch (item[0])
				{
			case #text:
				if item.Member?(#codepoints)
					s $= item[1] // entity or numeric characters
				else
					s $= .Encode(item[1])
			case #code:
				s $= '<code>' $ .Encode(item[1]) $ '</code>'
			case #html:
				s $= item[1]
			case #link:
				s $= '<a href="' $ .EncodeURL(item.href) $ '"' $
					Opt(' title="', .Encode(item.GetDefault(#title, '')), '"') $ '>' $
					.ConvertInline(item[1]) $ '</a>'
			case #image:
				s $= '<img src="' $ .EncodeURL(item.href) $
					'" alt="' $ .convertInlineToText(item[1]) $ '"' $
					Opt(' title="', .Encode(item.GetDefault(#title, '')), '"') $ ' />'
			case #linkbreak:
				s $= item.hard? ? '<br />\n' : '\n'
			case #emph:
				tag = item.strong? ? 'strong' : 'em'
				s $= '<' $ tag $ '>' $ .ConvertInline(item[1]) $ '</' $ tag $ '>'
				}
			}
		return s
		}

	convertInlineToText(parsedInline)
		{
		s = ''
		for item in parsedInline
			{
			if item[0] in (#link, #image, #emph)
				s $= .convertInlineToText(item[1])
			else if item[0] isnt #html
				s $= .Encode(item[1])
			}
		return s
		}

	Encode(s)
		{
		// &apos; are not necessary
		return s.
			Replace('&', '\&amp;').
			Replace('<', '\&lt;').
			Replace('>', '\&gt;').
			Replace('"', '\&quot;')
		}

	EncodeURL(s)
		{
		return Url.Encode(s).
			Replace('&', '\&amp;').
			Replace('+', '%20').
			Replace('%28', '(').
			Replace('%29', ')')
		}
	}