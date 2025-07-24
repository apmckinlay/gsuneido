// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
#(
ExtensionPoints:
	(
	('checks')
	)
Contributions:
	(
	('Formulas', 'checks',
		Fn: function (code) { return FormulaTonumber.CheckEmptyPlaceHolder(code) })
	)
)
