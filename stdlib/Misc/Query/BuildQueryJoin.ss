// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
// expects joinob from SelectFields.JoinsOb withDetails? (Object(:str, :fields))
class
	{
	CallClass(query, joinob)
		{
		if joinob.Empty?()
			return query

		if not query.Has?('extend')
			return query $ joinob.Fold("", { |s x| s $= x.str })

		// check if the query already has the exact same join
		newjoinob = Object()
		queryNoWhitespace = query.Tr(" \r\n\t")
		for join in joinob
			if not queryNoWhitespace.Has?(join.str.Tr(" \r\n\t"))
				newjoinob.Add(join)

		// remove any duplicate extends from the query side
		querycols = QueryColumns(query)
		removes = Object()
		for col in querycols
			if newjoinob.Any?({ it.fields.Has?(col) })
				removes.Add(col)

		return '(' $ query $
			(removes.Empty?() ? '' : ' remove ' $ removes.Join(', ')) $ ')' $
			newjoinob.Fold("", { |s x| s $= x.str })
		}
	}