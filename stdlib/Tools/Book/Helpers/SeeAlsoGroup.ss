// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// for use in help Asup
// e.g. <$ SeeAlsoGroup('GroupObjectFind') $>
// where /res/GroupObjectFind is a list of names one per line
// TODO take a comma separated list of groups (sort the names and remove duplicates)
BookHelperBase
	{
	CallClass(group, extra = "", _table = false, _path = false, _name = false)
		{
		if _table is false or _path is false or _name is false
			return 'SeeAlsoGroup requires _table, _path, _name'
		if false is x = Query1(_table, path: '/res/.groups', name: group)
			return 'ERROR: SeeAlsoGroup not found: ' $ group
		names = x.text.Lines()
		names.Remove(name)
		list = Object()
		for name in names
			{
			page = Query1(_table, :name, :path)
			if page is false
				page = QueryFirst(table $
					' where path.Has?("/Reference") and not path.Prefix?("<deleted>")'$
					' and name is ' $ Display(name) $ ' sort path')
			if page is false
				list.Add(name)
			else
				list.Add(.BuildLink(name,
					'/' $ Paths.Combine(_table, page.path) $ '/' $ name))
			}
		return .BuildParagraph('See also:\n' $ list.Join(',\n') $ Opt(', ', extra))
		}
	}