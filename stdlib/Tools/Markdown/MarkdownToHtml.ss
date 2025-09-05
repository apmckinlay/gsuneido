// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(md, noIndent? = false)
		{
		parsed = String?(md) ? MarkdownParser(md) : md
		return .Convert(parsed, :noIndent?)
		}

	writer: class
		{
		New(.noIndent? = false)
			{
			}
		indent: 0
		s: ''
		DoWithIndent(block)
			{
			if .noIndent? is false
				.indent++
			block()
			if .noIndent? is false
				.indent--
			}

		Add(tag, s)
			{
			if s is false
				{
				.s $= '\n' $ '\t'.Repeat(.indent) $ '<' $ tag $ ' />'
				return
				}

			.s $= '\n' $ '\t'.Repeat(.indent) $ '<' $ tag $ '>' $ s $ '</' $ tag $ '>'
			}

		AddPure(s)
			{
			.s $= s
			}

		AddWithBlock(tag, block, attr = #(), noNewline? = false)
			{
			attr = attr.Map2({ |m, v| m $ '="' $ v $ '"' }).Join(' ')
			newline = noNewline? ? '' : '\n'
			.s $= '\n' $ '\t'.Repeat(.indent) $ '<' $ tag $ Opt(' ', attr) $ '>'
			.DoWithIndent(block)
			.s $= newline $ '\t'.Repeat(.indent) $ '</' $ tag $ '>'
			}

		Suffix?(suffix)
			{
			return .s.Suffix?(suffix)
			}

		Get()
			{
			return .s.RemovePrefix('\n') $ '\n'
			}
		}

	Convert(parsed, noIndent? = false)
		{
		writer = new .writer(:noIndent?)
		.convertContainer(writer, parsed)
		return writer.Get()
		}

	convertContainer(writer, container, tight? = false)
		{
		container.ForEachBlockItem()
			{ |item|
			.convertItem(writer, item, :tight?)
			}
		}

	convertItem(writer, item, tight? = false)
		{
		switch
			{
		case item.Base?(Md_ATXheadings):
			writer.Add('h' $ item.Level, .convertInline(item.ParsedInline))
		case item.Base?(Md_Code):
			writer.Add('pre',
				'<code' $ Opt(' class="language-', item.Info.BeforeFirst(' '), '"') $
					'>' $ Opt(.encode(item.Codes.Join('\n')), '\n') $ '</code>')
		case item.Base?(Md_Html):
			writer.AddPure('\n' $ item.Html)
		case item.Base?(Md_Paragraph):
			if item.HeadingLevel isnt false
				writer.Add('h' $ item.HeadingLevel, .convertInline(item.ParsedInline))
			else if tight?
				writer.AddPure((writer.Suffix?('<li>') ? '' : '\n') $
					.convertInline(item.ParsedInline))
			else
				writer.Add('p', .convertInline(item.ParsedInline))
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

	noNewline?(tight?, container)
		{
		last = false
		container.ForEachBlockItem()
			{
			last = it
			}
		return last is false or tight? is true and last.Base?(Md_Paragraph)
		}

	convertInline(parsedInline)
		{
		s = ''
		for item in parsedInline
			{
			switch (item[0])
				{
			case #text:
				s $= .encode(item[1])
			case #code:
				s $= '<code>' $ .encode(item[1]) $ '</code>'
			case #html:
				s $= item[1]
			case #link:
				s $= '<a href="' $ Url.Encode(item.href).Replace('&', '\&amp;') $ '">' $
					.encode(item[1]) $ '</a>'
				}
			}
		return s
		}

	encode(s)
		{
		// &apos; are not necessary
		return s.
			Replace('&', '\&amp;').
			Replace('<', '\&lt;').
			Replace('>', '\&gt;').
			Replace('"', '\&quot;')
		}
	}