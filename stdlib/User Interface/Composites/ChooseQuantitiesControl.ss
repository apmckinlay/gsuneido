// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
// contributed by Ajith.R
ChooseField
	{
	Name: 'ChooseQuantities'
	New(list, listField = false, cols_head = #("List", "Number"), saveEmpty = true)
		{
		.field = listField
		.saveEmpty = saveEmpty
		.list = .build_list(list,.field)
		.cols_head = cols_head
		}
	build_list(list, field)
		{
		if field isnt false
			list = .Send("GetField", field)
		val_ob = Object()
		list_ob = Object()
		for i in list.Split(",")
			list_ob[i.BeforeFirst("(").Trim()] =
				Number(i.AfterFirst("(").BeforeLast(")").Trim())
		for i in .Get().Split(",")
			val_ob[i.BeforeFirst("(").Trim()] =
				Number(i.AfterFirst("(").BeforeLast(")").Trim())
		for ob in list_ob.Members().Sort!()
			{
			if val_ob.Member?(ob)
				list_ob[ob] = val_ob[ob]
			else
				if not .saveEmpty
					list_ob[ob] = 0
			}
		str = ""
		for i in list_ob.Members()
			str $=  i $ "(" $ list_ob[i] $ "),"
		return str[.. -1]
		}
	Getter_DialogControl()
		{
		.list = .build_list(.list, .field)
		return Object('Quantities',
			list: .list, cols_head: .cols_head, saveEmpty: .saveEmpty)
		}
	}
