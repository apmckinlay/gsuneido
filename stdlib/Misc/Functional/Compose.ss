// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
// Compose(f, g) => function (@args) { g(f(@args) }
// Note: the functions are applied in the order listed
// e.g. Compose(Add, Curry(Mul, 2))(4, 1) => 10
function (@fns)
	{
	return {|@args|
		result = (fns[0])(@args)
		for fn in fns[1..]
			result = fn(result)
		result
		}
	}