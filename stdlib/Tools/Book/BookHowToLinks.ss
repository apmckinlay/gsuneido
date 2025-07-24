// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
class
	{
	hdiSplitOn: `How Do I ...?/`
	CallClass(page, helpbook, skipTitle? = false, opentag = 'li', closetag = 'li',
	skipULTags? = false)
		{
		header = '\n<!-- HowToLinks -->'
		if not TableExists?(helpbook $ "HowToIndex")
			return ''
		QueryApply(helpbook $ "HowToIndex where name is " $ Display(page))
			{ |x|
			name = x.name
			path = ""
			if name =~ '/'
				{
				ob = name.SplitOnLast('/')
				name = ob[1]
				path = ob[0]
				}
			if false is BookPageFind(helpbook, path, name)
				continue
			text = header
			if not skipTitle?
				text $= '<p>See Also:</p>\n'
			if opentag is 'li' and not skipULTags?
				text $= '<ul class="howto">\n'
			for option in x.howtos
				{
				spliton = option.Has?(.hdiSplitOn) ? .hdiSplitOn : '/'
				hdiname = option.AfterLast(spliton)
				hdipath = option.BeforeLast(spliton) $ spliton[..-1] // remove the last /
				y = BookPageFind(helpbook, hdipath, hdiname)
				if y is false or 'hidden' is BookEnabled(helpbook, hdiname)
					continue
				text $= '<' $ opentag $ '><a href=' $ Display('/' $ helpbook $ hdipath $
					'/' $ hdiname) $ '>' $ hdiname $ '</a></' $ closetag $ '>\n'
				}
			if opentag is 'li' and not skipULTags?
				text $= '</ul>\n'
			return text.Trim()
			}
		return ''
		}
	}
