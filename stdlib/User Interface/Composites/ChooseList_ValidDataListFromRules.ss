// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(args)
		{
		split = args.GetDefault('splitValue', ',')
		list = .SplitListFromString(args.record[args.listField], split)
		if not Object?(list)
			return list // just in case it's not an object, caller checks for this
		if args.Member?('allowOtherField') and
			Object?(othr = .SplitListFromString(args.record[args.allowOtherField], split))
			list.MergeUnion(othr)
		return list
		}

	SplitListFromString(list, splitValue)
		{
		if not String?(list)
			return list
		return list.Split(splitValue).Map!(#Trim)
		}
	}