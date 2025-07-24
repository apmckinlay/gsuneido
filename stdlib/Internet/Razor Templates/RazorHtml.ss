// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
class
	{
	opening_tag: "^<([a-zA-Z_][a-zA-Z0-9_-]*)"
	CallClass(dest, s)
		{
		Dbg(HTML: s)
		s = s.LeftTrim()
		outer_tag = s.Extract(.opening_tag)
		.assertOpeningTag(outer_tag)
		opening_tag = '<' $ outer_tag
		closing_tag = '</' $ outer_tag $ '>'
		at_start = true
		forever
			{
			len = s.Size()
			at_pos = s.Find('@')
			closing_pos = s.Find(closing_tag)
			.assertClosingTag(closing_pos, len, closing_tag)

			opening_pos = s.Find(opening_tag, at_start ? 1 : 0)
			if closing_pos < at_pos and opening_pos > closing_pos
				break
			s = .output(dest, s, Min(at_pos, opening_pos))
			if at_pos < opening_pos
				s = RazorCode(dest, s)
			else if opening_pos < len
				s = RazorHtml(dest, s)
			at_start = false
			Dbg(MORE_HTML: s)
			}
		if outer_tag is 'text'
			s = .output(dest, s[6 ..], closing_pos - 6) /*= size of <text>*/
		else
			s = .output(dest, s, closing_pos + closing_tag.Size())
		return s
		}
	assertOpeningTag(outer_tag)
		{
		if outer_tag is false
			throw "Razor html section must start with opening tag"
		}
	assertClosingTag(closing_pos, len, closing_tag)
		{
		if closing_pos is len
			throw "Razor html section missing closing tag " $ closing_tag
		}
	output(dest, s, n)
		{
		dest.Html(s[.. n])
		return s[n ..]
		}
	}