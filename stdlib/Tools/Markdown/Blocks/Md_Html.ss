// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Base
	{
	New(.Html, .end)
		{
		if .Html =~ end
			.Close()
		}

	Match(line, container)
		{
		if false is indentN = .IgnoreLeadingSpaces(line)
			return false

		for i, condition in .conditions
			{
			// rule 1 - 6 may interrupt a paragraph
			if container is false and i is 6 /*=rule #7*/
				continue
			match? = String?(condition.start)
				? line[indentN..] =~ condition.start
				: (condition.start)(line[indentN..])
			if match?
				return new this(line, condition.end)
			}

		return false
		}

	conditions: (
		(start: '(?i)^<(pre|script|style|textarea)( |\t|>|$)',
			end: '(?i)</(pre|script|style|textarea)>'),
		(start: `^<!--`,
			end: `-->`),
		(start: `^(?q)<?`,
			end: `(?q)?>`),
		(start: `^<![a-zA-Z]`,
			end: `>`),
		(start: `^(?q)<![CDATA[`,
			end: `(?q)]]>`),
		(start: `^(?i)(<|</)(address|article|aside|base|basefont|blockquote|body|` $
			`caption|center|col|colgroup|dd|details|dialog|dir|div|dl|dt|fieldset|` $
			`figcaption|figure|footer|form|frame|frameset|h1|h2|h3|h4|h5|h6|head|` $
			`header|hr|html|iframe|legend|li|link|main|menu|menuitem|nav|noframes|` $
			`ol|optgroup|option|p|param|search|section|summary|table|tbody|td|tfoot|` $
			`th|thead|title|tr|track|ul)( |\t|>|/>|$)`,
			end: `^\s*$`),
		(start: function (line)
			{
			if false is length = Md_Helper.MatchHTMLTag(line)
				return false
			return line[length..].Blank?()
			},
			end: `^\s*$`))

	Continue(line)
		{
		if line =~ .end
			{
			.Close()
			return .BlankLine?(line) ? false : line
			}
		return line
		}

	Add(line)
		{
		.Html $= '\n' $ line
		}

	Close()
		{
		.Html = .Html
		super.Close()
		}
	}