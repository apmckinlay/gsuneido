// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// TODO: handle comments within class:base
class
	{
	New(libs = false)
		{
		if libs is false
			libs = Libraries()
		.libs = libs
		.names = Object()
		.names[0] = 'class'
		.children = Object()
		for i in libs.Members()
			{
			lib = libs[i]
			QueryApply(lib)
				{ |x|
				base = .base(x.text)
				if base is false
					continue
				num = .make_num(i, x.num)
				.names[num] = x.name
				if not .children.Member?(base)
					.children[base] = Object()
				.children[base].Add(
					Object(:num, name: x.name, group:))
				}
			}
		}
	base(text)
		{
		scan = ScannerWithContext(text)
		base = scan.Next()
		if base is 'class'
			{
			token = scan.Next()
			if token is ':'
				base = scan.Next()
			else if token isnt '{'
				base = token
			}
		else if (not String?(base) or base !~ "^[A-Z]")
			return false
		return base
		}
	offset: 100000
	make_num(lib, num)
		{ return lib * .offset + num }
	Children(parent)
		{
		parentName = .name(parent)
		if not .children.Member?(parentName)
			return #()
		list = .children[parentName]
		list.Sort!(function (x, y) { return x.name < y.name })
		return list
		}
	Children?(parent)
		{
		return .group?(.name(parent))
		}
	name(num)
		{
		return .names[num]
		}
	group?(name)
		{
		return .children.Member?(name)
		}
	Get(num)
		{
		if false isnt x = Query1(library = .Get_lib(num), num: .Get_num(num))
			{
			x.group = .group?(x.name)
			x.table = library
			}
		return x
		}

	Get_lib(num)
		{
		return .libs[(num / .offset).Int()]
		}
	Get_num(num)
		{
		return num % .offset
		}
	Nextnum()
		{
		}
	NewItem(x /*unused*/)
		{
		}
	Update(x /*unused*/)
		{
		}
	DeleteItem(num /*unused*/)
		{
		}
	Container?(x)
		{
		return .Children?(x)
		}
	Static?(x /*unused*/)
		{
		return true
		}
	FindNames(pat)
		{
		list = Object()
		for base in .children.Members()
			for child in .children[base]
				if child.name =~ pat
					list.Add(child.name)
		return list
		}
	GetPath(name)
		{
		path = Object(name)
		while ('class' isnt (name = .getbase(name)))
			path.Add(name)
		return path
		}
	getbase(name)
		{
		for base in .children.Members()
			for child in .children[base]
				if child.name is name
					return base
		return false
		}
	}
