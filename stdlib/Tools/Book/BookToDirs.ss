// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(book, dest = '.')
		{
		if DirExists?(path = Paths.Combine(dest, book))
			throw "already exists: " $ Display(path)
		_table = book
		.export(book, "", Paths.Combine(dest, book))
		}
	export(book, path, dest)
		{
		output = 0
		filepath = ""
		out = {|fp,text|
			if output++ is 0
				EnsureDirectories(filepath)
			text = .links(path, text)
			PutFile(fp, text)
			}
		n = 0
		QueryApply(book, :path)
			{|x|
			name = x.name
			if x.name.Prefix?('.') or x.name is "Cover"
				continue
			filepath = Paths.Combine(dest, .filename(name))
			text = x.text
			if x.path =~ `^/res\>` and x.text isnt ""
				out(filepath, x.text)
			else
				{
				if not text.Prefix?("<")
					text = text.Eval()
				if text.Prefix?("<")
					{
					_path = path
					_name = x.name
					out(filepath $ ".html", Asup(text))
					}
				}
			.export(book, Paths.Combine(x.path, x.name), filepath)
			}
		}
	filename(name)
		{
		return name.Tr('?!', 'QX').Tr('^a-zA-Z0-9. -', '_')
		}
	links(from, text)
		{
		return text.Replace(`[("](suneido:)?/suneidoc/.*?[)"]`, {
			path = it[1..-1].RemovePrefix("suneido:").RemovePrefix("/suneidoc")
			if path !~ `[.](gif|png|jpg)`
				path $= '.html'
			it[0] $ Paths.AbsToRel(from, path) $ it[-1]
			})
		}
	}