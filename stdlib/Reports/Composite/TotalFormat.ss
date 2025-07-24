// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
VertFormat
	{
	New(item, skip = .16, noline = false)
		{
		super(.LineLayout(noline),
			.getItem(item),
			Object('Vskip' skip))
		.skip = .GetItems().Last().GetSize().h
		}
	LineLayout(noline)
		{
		return Object('Hline', before: 50, after: 50, :noline)
		}

	getItem(item)
		{
		if String?(item)
			return Object(item, xstretch: 1)
		item = item.Copy()
		item.xstretch = 1
		return item
		}
	GetSize(data = #{})
		{
		size = super.GetSize(data).Copy()
		item = .GetItems()[1]
		value = item.Member?('Field') ? data[item.Field] : data
		size.d = item.GetSize(value).d + .skip
		return size
		}
	}
