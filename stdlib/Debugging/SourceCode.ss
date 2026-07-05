// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function (fn)
	{
	if Type(fn) not in (#Function, #Method, #Block, #Class)
		return false
	name = Name(fn).BeforeFirst('.').BeforeFirst(' ')
	s = Display(fn)
	if not s.Has?('/* ')
		return false
	lib = s.AfterFirst('/* ').BeforeFirst(' ')
	if lib is ""
		return false
	name = lib.Has?('__') ? name $ '__' $ lib.AfterFirst('__') : name
	lib = lib.BeforeFirst('__')
	if lib in ('function', 'block', 'builtin') or lib !~ '^[A-Za-z][_a-zA-Z0-9]*$'
		return false
	try
		return Query1(lib, group: -1, :name).text
	return false
	}
