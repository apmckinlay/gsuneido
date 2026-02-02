// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Default(@args)
		{
		(.getBuilder()[args[0]])(@+1args)
		}

	getBuilder(_outputType = #html)
		{
		return outputType is #html
			? .htmlBuilder
			: .mdBuilder
		}

	htmlBuilder: class
		{
		BuildHeading(text, level, style = false)
			{
			args = Object('h' $ level, text)
			if style isnt false
				args.style = style
			return Xml(@args) $ '\n'
			}

		BuildLink(text, href)
			{
			return Xml('a', text, :href)
			}

		BuildImage(book, name)
			{
			return '<img src="suneido:/' $ book $ '/res/' $ name $ '" align="middle"> '
			}

		BuildTable(table)
			{
			s = '\n<table width="100%">\n'
			for row in table
				{
				s $= '\n<tr>\n'
				for col in row
					s $= '\t' $ Xml('td', col) $ '\n'
				s $= '\n</tr>\n'
				}
			s $= '\n</table>\n'
			return s
			}

		BuildParagraph(text)
			{
			return Xml('p', text)
			}
		}

	mdBuilder: class
		{
		BuildLink(text, href)
			{
			return '[' $ text $ '](<' $ href $ '>)'
			}

		BuildHeading(text, level)
			{
			return '#'.Repeat(level) $ ' ' $ text $ '\n\n'
			}

		BuildImage(book, name)
			{
			return '![](<suneido:/' $ book $ '/res/' $ name $ '>)'
			}

		BuildTable(table, cols)
			{
			s = '|' $ '     |'.Repeat(cols) $ '\n'
			s $= '|' $ ' --- |'.Repeat(cols) $ '\n'
			for row in table
				{
				s $= '| ' $ row.Join(' | ') $ ' |\n'
				}
			return s $ '\n'
			}

		BuildParagraph(text)
			{
			return '\n' $ text $ '\n'
			}
		}
	}