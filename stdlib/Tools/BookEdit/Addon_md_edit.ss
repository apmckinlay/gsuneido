// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Addon_html_edit
	{
	AddLink(table, result, text)
		{
		pos = .GetSelect()
		url = '/' $ table $ result.path $ "/" $ result.name
		namepos = result.name.Find(text)
		text = '[' $ result.name $ '](<' $ url $ '>)'
		.Paste(text)
		.SetSelect(pos.cpMin + 1 + namepos, pos.cpMax - pos.cpMin)
		}

	On_Add_Image_Tag()
		{
		text = .GetSelText()
		table = .Send('CurrentTable')
		.Paste('![](<suneido:/' $ table $ '/res/' $ text $ '>)')
		}

	On_Bold()
		{
		.insertFormat('**')
		}
	On_Italic()
		{
		.insertFormat('*')
		}
	On_Code()
		{
		.insertFormat('`')
		}
	insertFormat(symbol)
		{
		text = .GetSelText()
		trailing_whitespace = text.Extract('[ \t\r\n]*$')
		text = text.RightTrim()
		size = symbol.Size()
		offset = 0
		pos = .GetSelect()
		str = .GetRange(pos.cpMin - size, pos.cpMax + size)

		if (not str.Prefix?(symbol) or not str.Suffix?(symbol))
			{
			offset = size
			size = 0
			text = symbol $ text $ symbol
			}

		.SetSelect(pos.cpMin - size, (pos.cpMax - pos.cpMin) + size * 2)
		.Paste(text $ trailing_whitespace)
		.SetSelect(pos.cpMin + offset - size, pos.cpMax - pos.cpMin)
		}

	On_Add_Paragraph_Tags() { }
	On_H1() {}
	On_H2() {}
	On_H3() {}
	On_H4() {}
	On_P() {}
	On_LI() {}
	On_DT() {}
	On_DD() {}
	On_PRE() {}
	}