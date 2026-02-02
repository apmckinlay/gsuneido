// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(book = 'suneidoc', dest = '.')
		{
		if DirExists?(path = Paths.Combine(dest, book))
			throw "already exists: " $ Display(path)
		_table = book
		_outputType = #md
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
				if text !~ `\A[<#]`
					text = text.Eval()
				if text =~ `\A[<#]`
					{
					_path = path
					_name = x.name
					out(filepath $ ".md", Asup(text))
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
		return text.Replace(`<(suneido:)?/suneidoc/.*?>`, {
			path = it[1..-1].RemovePrefix("suneido:").RemovePrefix("/suneidoc")
			if path !~ `[.](gif|png|jpg)`
				path $= '.md'
			it[0] $ Paths.AbsToRel(from, path) $ it[-1]
			})
		}
	}