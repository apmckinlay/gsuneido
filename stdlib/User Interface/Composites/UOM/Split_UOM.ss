// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (field, strips = false)
	{
	if not String?(field)
		return Object(value: field, uom: "")

	value = field.BeforeFirst(" ")
	uom = field.AfterFirst(" ")
	if not value.Tr(",").Number?()
		{
		value = ''
		uom = field
		}
	uom = uom.Trim()

	if strips
		uom = uom.Replace("[s|S]$", "")
	return Object(value: value.Tr(",").Number?() ? Number(value) : value, :uom)
	}