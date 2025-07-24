// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	On_Link()
		{
		if "" is text = .GetSelText().Trim()
			return
		table = .Send('CurrentTable')
		matches = QueryAll(table $
			' where name is ' $ Display(text) $ ' and path !~ "^/res\>"' $
			' project name, path')

		if matches.Size() is 0
			{
			if not YesNo(
				'No exact matches found.\nWould you like to seach for partial matches?',
				'Alert', .Window.Hwnd, MB.ICONQUESTION)
				return
			matches = QueryAll(table $
				' where name =~ ' $ Display('(?i)(?q)' $ text) $
					' and path !~ "^/res\>"' $
				' project name, path')
			if matches.Size() is 0
				{
				Alert('No matches found', title: 'Link', flags: MB.ICONERROR)
				return
				}
			}

		result = matches.Size() is 1
			? matches[0]
			: BookEditLinkListControl(.Window.Hwnd, matches)

		if (result is false)
			return

		pos = .GetSelect()
		tag = '<a href="' $ '/' $ table $ result.path $ "/" $ result.name $ '">'
		namepos = result.name.Find(text)
		text =  tag $ result.name $ "</a>"
		.Paste(text)
		.SetSelect(pos.cpMin + tag.Size() + namepos, pos.cpMax - pos.cpMin)
		}

	On_Add_Paragraph_Tags()
		{
		src = .Get()
		data = Object(newText: '', paragraph: '')
		pre? = false
		for line in src.Lines()
			{
			line = line
			if pre?
				{
				data.newText $= line $ "\n"
				if line is "</pre>"
					{
					data.newText $= "\n"
					pre? = false
					}
				}
			else if line is "<pre>"
				{
				.addParagraph(data)
				pre? = true
				data.newText $= line $ "\n"
				}
			else if line is ''
				.addParagraph(data)
			else
				{
				line = line.Trim()
				data.paragraph $= line $ (line =~ "<br( ?/)?>" ? "\n" : " ")
				}
			}
		.addParagraph(data)
		.PasteOverAll(data.newText)
		.On_Refresh()
		}

	addParagraph(data)
		{
		paragraph = data.paragraph.Trim()
		if paragraph isnt ''
			{
			if not paragraph.Prefix?('<')
				paragraph = Xml('p', paragraph)
			data.newText $= paragraph $ "\n\n"
			}
		data.paragraph = ''
		}


	On_H1()
		{
		.insertFormat('h1')
		}
	On_H2()
		{
		.insertFormat('h2')
		}
	On_H3()
		{
		.insertFormat('h3')
		}
	On_H4()
		{
		.insertFormat('h4')
		}
	On_P()
		{
		.insertFormat('p')
		}
	On_LI()
		{
		.insertFormat('li')
		}
	On_DT()
		{
		.insertFormat('dt')
		}
	On_DD()
		{
		.insertFormat('dd')
		}
	On_PRE()
		{
		.insertFormat('pre')
		}
	On_Bold()
		{
		.insertFormat('b')
		}
	On_Italic()
		{
		.insertFormat('i')
		}
	On_Underline()
		{
		.insertFormat('u')
		}
	On_Code()
		{
		.insertFormat('code')
		}
	insertFormat(tag)
		{
		text = .GetSelText()
		trailing_whitespace = text.Extract('[ \t\r\n]*$')
		text = text.RightTrim()
		tag1 = 2 + tag.Size()
		tag2 = 3 + tag.Size()
		offset = 0
		pos = .GetSelect()
		str = .GetRange(pos.cpMin - 2 - tag.Size(), pos.cpMax + tag2)

		if (not str.Prefix?('<' $ tag $ '>') or not str.Suffix?('</' $ tag $ '>'))
			{
			offset = tag1
			tag1 = tag2 = 0
			text = "<" $ tag $ ">" $ text $ "</" $ tag $ ">"
			}

		.SetSelect(pos.cpMin - tag1, (pos.cpMax - pos.cpMin) + tag2 + tag1)
		.Paste(text $ trailing_whitespace)
		.SetSelect(pos.cpMin + offset - tag1, pos.cpMax - pos.cpMin)
		}
	}