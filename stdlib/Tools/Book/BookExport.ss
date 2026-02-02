// Copyright (C) 2000 Suneido Software Corp.
// e.g. BookExport('book', 'c:/book')
class
	{
	New(.book, .dir)
		{
		if false isnt x = Query1(book, name: "Copyright")
			.copyright = x.text
		.list = Object()
		.pages()
		.export()
		.res()
		}
	copyright: false
	pages(path = "")
		{
		QueryApply(.book $ " where path = " $ Display(path) $
			" sort order, name")
			{|x|
			if .html(x.text) is false
				continue
			name = x.path $ "/" $ x.name
			.list.Add(name)
			.pages(name) // do children (if any)
			}
		}
	html(text)
		{
		if not BookContent.Match(.book, text)
			text = text.Eval() // needs to use Eval
		else
			text = BookContent.ToHtml(.book, text)
		if text.Prefix?('<')
			return text
		return false
		}
	sanitize(name)
		{
		return name.Tr('?', 'q').Tr('^a-zA-Z0-9/', '_')
		}
	export(path = "") // recursive
		{
		first = true
		QueryApply(.book $ " where path = " $ Display(path) $
			" sort order, name")
			{|x|
			if false is text = .html(x.text)
				continue
			if first
				{
				EnsureDir(.dir $ .sanitize(path))
				first = false
				}
			name = Paths.Combine(x.path, x.name)
			filename = .dir $ .sanitize(name) $ ".htm"

			nav = .nav(x, name)

			text = nav $ text $ '<br />' $ nav
			if .copyright isnt false
				text $= '<div align="center" style="color="gray"; font-size: "8pt";
					margin-top: 6pt text-align: center;">' $ .copyright $ '</div>'
			_table = .book // needed by Asup
			_path = x.path
			_name = x.name
			text = HtmlWrap(text, .book)
			text = .links(text, path)
			text = text.Replace('<head>',
				'<head>\n<title>Suneido ' $ name $ '</title>', 1)
			PutFile(filename, text)
			.export(name) // do children (if any)
			}
		}
	nav(x, name)
		{
		nav = '<table width="100%" bgcolor="lightgrey" border="0"
			cellspacing="0" style="margin-bottom: 0;">
			<tr style="font-size: 75%"><td style="width: 12em;">'
		// left
		i = .list.Find(name)
		if i > 0
			{
			prev = .list[i - 1]
			nav $= .ahref(prev, '&lt;&lt Previous')
			}
		nav $= '&nbsp;</td><td align="center">'
		// center
		nav $= '<a href="http://www.suneido.com" target="_top">Suneido</a> &gt '
		nav $= .ahref("/Contents", "Contents")
		paths = x.path[1..].Split('/')
		for j in paths.Members()
			{
			link = '/' $ paths[.. j + 1].Join('/')
			nav $= ' &gt ' $ .ahref(link, paths[j])
			}
		nav $= '</td><td align="right" style="width: 12em;">'
		// right
		if i + 1 < .list.Size()
			{
			next = .list[i + 1]
			nav $= .ahref(next, 'Next &gt;&gt;')
			}
		nav $= '&nbsp;</td></tr></table>\n'
		return nav
		}
	ahref(link, label)
		{
		return '<a href="' $ link $ '">' $ label $ '</a>'
		}
	links(src, path) // convert links
		{
		dst = ''
		forever // <a href=
			{
			i = src.Find('<a href=')
			dst $= src[.. i]
			src = src[i..]
			if src is ''
				break
			len = src.Find('>') + 1
			ahref = src[.. len]
			link = .link(path, ahref[9 .. -2])
			dst $= '<a href="' $ link $ '">'
			src = src[len..]
			}
		src = dst
		dst = ""
		forever // <img
			{
			i = src.Find('<img ')
			dst $= src[.. i]
			src = src[i..]
			if src is ''
				break
			len = src.Find('>') + 1
			img = src[.. len]
			link = .link(path, img.Extract('src="([^"]*)'))
			dst $= img.Replace('src="([^"]*)', 'src="' $ link)
			src = src[len..]
			}
		dst = dst.Replace("\<url\([^)]+\)")
			{|url|
			'url(' $ .link(path, url[4 .. -1]) $ ')'
			}
		return dst
		}
	link(path, link) // convert one link
		{
		if link.Prefix?("http://")
			return link
		if link.Prefix?("suneido:") // images
			link = link[8..]
		else if not link.Prefix?('#')
			link = .sanitize(link) $ '.htm'
		if link !~ '/'
			return link
		link = link.RemovePrefix('/' $ .book)
		if path is ''
			path = '/'
		xpath = path.Split('/').Set_default(true)
		xlink = link.Split('/').Set_default(false)
		xrel = Object()
		// skip common prefix
		for (i = 0; xpath[i] is xlink[i]; ++i)
			{ }
		// go up to common parent
		for (j = i; j < xpath.Size(); ++j)
			xrel.Add("..")
		// add rest of link
		for (; i < xlink.Size(); ++i)
			xrel.Add(xlink[i])
		return xrel.Join('/')
		}
	res()
		{
		EnsureDir(.dir $ "/res")
		QueryApply(.book $ " where path =~ '^/res\>' and path !~ '/[.]'")
			{|x|
			path = Paths.Combine(.dir $ x.path, x.name)
			EnsureDirectories(path)
			PutFile(path, x.text)
			}
		}
	}