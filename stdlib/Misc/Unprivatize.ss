// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// Name.Name_priv => Name.priv
function (s)
	{
	if false isnt name = s.Extract(`\A([A-Z]\w*)\.`)
		return s.Replace(`\A` $ name $ '.' $ name $ '_', name $ '.')
	return s
	}