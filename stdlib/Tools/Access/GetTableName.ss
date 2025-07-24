// Copyright (C) 2014 Axon Development Corporation All rights reserved worldwide.
function (table)
	{
	if table.Prefix?('custom_') and table.Suffix?('_table')
		return SelectPrompt(table.RemoveSuffix('_table'))

	if table =~ 'custom_table_\d\d\d\d\d\d\d\d\d\d\d\d'
		return Customizable.CustomTableName(table)

	if false isnt result = Tables.GetTable(table, 'Name')
		return result

	SuneidoLog("ERROR: unable to find name for table: " $ table, calls:)
	return table
	}