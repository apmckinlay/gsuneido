// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function (record, field, dd = false, value = #(0))
	{
	value = Same?(value, #(0)) ? record[field] : value
	if value is ""
		return true

	if dd is false
		dd = Datadict(field)
	args = dd.Control.Copy()
	args[0] = value
	args.record = record.Copy()
	args.fieldToValidate = field
	try
		ctrl = Global(dd.Control[0] $ 'Control')
	catch (err /*unused*/, "can't find")
		return true

	return ctrl.Method?('ValidData?')
		? ctrl.ValidData?(@args) : true
	}
