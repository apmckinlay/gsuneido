// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function (fn)
	{
	try
		return fn.Source()
	if Type(fn) not in (#Function, #Method, #Block, #Class)
		return false
	name = Name(fn).BeforeFirst('.').BeforeFirst(' ')
	s = Display(fn)
	lib = s.AfterFirst('/* ').BeforeFirst(' ')
	if lib isnt "function"
		try
			return Query1(lib, group: -1, :name).text
	return false
	}