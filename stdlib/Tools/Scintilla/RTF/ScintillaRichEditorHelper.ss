// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Parse(htmlStr)
		{
		parser = ParseHTMLRichText.GetParsedText(htmlStr)
		if parser is false
			return Object(s: htmlStr, styles: #())

		s = ""
		styles = Object()
		pos = Object(line: 0, ch: 0)
		for child in parser.Children()
			s $= .buildStyled(child, styles, pos)
		return Object(:s, :styles)
		}

	buildStyled(child, styleObject, pos)
		{
		style = .getStyle(child)
		txt = child.Text().Replace('&crlf;', '\r\n')
		if txt is ""
			return txt
		from = pos.Copy()
		.NextPos(txt, from, pos)
		styleObject.Add(Object(:from, to: pos.Copy(), :style, :txt))
		return txt
		}

	DefaultStyle: #(bold: false, italic: false, underline: false, strikeout: false)
	getStyle(child)
		{
		style = .DefaultStyle.Copy()
		// need to handle text that has been appended to the end
		// might not have rtf formatting
		format = child.Attributes().Member?('style')
			? child.Attributes().style
			: ""

		if format =~ (':bold(;|$)')
			style.bold = true
		if format =~ (':italic(;|$)')
			style.italic = true
		if format.Has?('underline')
			style.underline = true
		if format.Has?('line-through')
			style.strikeout = true

		return style
		}

	NextPos(txt, from, to)
		{
		lines = txt.Lines()
		if txt.Suffix?('\r\n')
			lines.Add('')
		to.line = from.line + lines.Size() - 1
		to.ch = lines.Size() > 1 ? 0 : from.ch
		to.ch += lines.Last().Size()
		}

	TrimStyledText(s, styles)
		{
		chars = "^ \t\r\n"
		first = s.Find1of(chars)
		last = s.FindLast1of(chars)
		if last + 1 - first is s.Size()
			return Object(:s, :styles)

		if s isnt "" and styles.Empty?()
			throw 'styles expected but not found'

		styles.RemoveIf({ it.txt.Blank?() })

		if styles.Empty?()
			return Object(s: s.Trim(), :styles)

		styles[0].txt = styles[0].txt.LeftTrim()
		styles[styles.Size()-1].txt = styles[styles.Size()-1].txt.RightTrim()

		// might be able to do this without needing to rebuild and reparse the html
		html = ""
		for style in styles
			{
			html $= "<" $ .setFontStyle(style.style) $ ">"
			html $= XmlEntityEncode(style.txt)
			html $= "</span>"
			}
		ob = .Parse(html)

		return Object(s: ob.s, styles: ob.styles)
		}

	Build(s, styles)
		{
		if styles.Empty?()
			return s

		lines = s.Lines()
		if s.Suffix?('\r\n')
			lines.Add('')
		starts = Object(0)
		for i in lines.Members()
			starts[i + 1] = starts[i] + lines[i].Size() + 2

		str = ''
		for style in styles
			{
			styleFrom = style.from
			styleTo = style.to
			from = starts[styleFrom.line] + styleFrom.ch
			to = starts[styleTo.line] + styleTo.ch
			str $= '<' $ .setFontStyle(style.style) $ '>' $
				XmlEntityEncode(s[from..to]).
					Replace('\r', '').Replace('\n', '<br />') $ '</span>'
			}
		return str
		}

	setFontStyle(style)
		{
		base = 'span style="'
		bold = "font-weight:" $ (style.bold ? "bold" : "normal")
		italic = "font-style:" $ (style.italic ? "italic" : "normal")
		underline = style.underline ? "underline" : ""
		strikethrough = style.strikeout ? " line-through" : ""
		return base $ bold $ ';' $ italic $
			Opt(';text-decoration:', (underline $ strikethrough)) $ '"'
		}
	}
