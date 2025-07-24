// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
/*
For inheriting classes:
Table: 'table_name'
Name: 'Readable table name'
Columns: (<all the fields in the table>)
Keys: ('single_key', 'combo, key') // if no keys, use ('')
// if key is also foreign key
Keys: (('otherTableKey', in: 'another_table', cascade: t/f))
UniqueIndexes: ('unique')
Indexes: ('single_index', 'combo, index')
ForeignKeys: (table_name: ((from: 'field', to: 'optional_renamed_field', cascade: t/f)),
	bizpartners: ((from: 'bizpartner_num'), // index(bizpartner_num) in biz_partners
		(from: 'arfc_num',
			to: 'bizpartner_num') // index(arfc_num) in biz_partners(bizpartner_num)
		),
	ap_checks: ((from: 'apchk_num', cascade:)) // index(apchk_num) in ap_checks cascade
	)
BookLocation: Plugin_<>Reporter bookLocation
Permission: Plugin_<>Reporter auth
DuplicateFieldInfo: function()
	{
	return Object(
		fields: Customizable.CustomFields('biz_partners'),
		configTable: 'company',
		configField: 'company_checkdup_bizpartnerfields')
	}
*/
class
	{
	Table: false
	Name: false
	contribs: false

	// only used if table is shared between screens, so customization is a different name
	RenamedForCustomize?: false
	// GetContributions uses Contributions which is Memoized
	New()
		{
		if false is .contribs = GetContributions('Table_' $ .Table)
			.contribs = #()
		}
	// these are called in SchemaChecker
	Columns: ()
	GetColumns()
		{
		return .Columns.Union(.getContribs('Columns'))
		}
	getContribs(type)
		{
		return .contribs.GetDefault(type, #())
		}

	Keys: ()
	GetKeys()
		{
		return .Keys.Copy().Append(.getContribs('Keys'))
		}

	UniqueIndexes: ()
	GetUniqueIndexes()
		{
		return .UniqueIndexes.Copy().Append(.getContribs('UniqueIndexes'))
		}

	Indexes: ()
	GetIndexes()
		{
		indexes = .Indexes.Copy().Append(.getContribs('Indexes'))
		ob = .GetDuplicateFieldInfo()
		if ob isnt false and TableExists?(ob.configTable)
			{
			config = Query1(ob.configTable)
			if config isnt false
				{
				indexes.MergeUnion(config[ob.configField])
				// does not need to be indexes if it is already a key or unique index
				indexes = indexes.Difference(.GetUniqueIndexes()).Difference(.GetKeys())
				}
			}
		return indexes
		}

	ForeignKeys: ()
	GetForeignKeys()
		{
		// Can't use .Append here as there might be multiple foreign keys to the same
		//	table, and .Append will just overwrite any in the Base list.
		bContrib = .ForeignKeys.Copy().Set_default(Object())
		fkContrib = .getContribs('ForeignKeys')
		for m in fkContrib.Members()
			bContrib[m] = bContrib[m].Copy().Append(fkContrib[m])
		return bContrib
		}

	// first pass focus on what is excluded on the data entry screen
	ExcludeFields: ()
	GetExcludeFields()
		{
		return .ExcludeFields.Copy().Append(.getContribs('ExcludeFields'))
		}

	Schema(withoutForeignKeys = false)
		{
		// make sure to use the immutable version of map
		schema = .Table $ ' ' $
			'(' $ .GetColumns().Join(',') $ ') ' $
			.GetKeys().Map(.schemaKeys).Join(' ') $ ' ' $
			.GetUniqueIndexes().Map(.schemaUniqueIndexes).Join(' ') $ ' ' $
			.GetIndexes().Map(.schemaIndexes).Join(' ')
		if not withoutForeignKeys
			schema $= ' ' $ .GetForeignKeys().Map2(.schemaForeignKeys).Join(' ')
		return schema
		}
	schemaKeys(key)
		{
		if Object?(key)
			{
			return 'key (' $ key[0] $ ')' $
				(key.Member?('in') ? ' in ' $ key.in : "") $
				(key.Member?('cascade') ? ' cascade' : "")
			}
		return 'key (' $ key $ ')'
		}
	schemaUniqueIndexes(index)
		{
		return 'index unique(' $ index $ ')'
		}
	schemaIndexes(index)
		{
		return 'index (' $ index $ ')'
		}
	schemaForeignKeys(table, fields)
		{
		fks = Object()
		for key in fields
			{
			fk = 'index(' $ key.from $ ') in ' $ table
			if key.Member?('to')
				fk $= ' (' $ key.to $ ')'
			if key.Member?('cascade')
				fk $= ' cascade'
			if key.Member?('cascade_update')
				fk $= ' cascade update'
			fks.Add(fk)
			}
		return fks.Join(' ')
		}
	Ensure(withoutForeignKeys = false)
		{
		Database('ensure ' $ .Schema(:withoutForeignKeys))
		}

	DuplicateFieldInfo: function() { return false }
	GetDuplicateFieldInfo()
		{
		info = (.DuplicateFieldInfo)()
		if info is false
			info = Object(fields: Object())
		else
			{
			info = info.Copy()
			info.fields = info.fields.Copy().MergeUnion(.getCustomFields())
			}

		for func in .getContribs('DuplicateFieldInfo')
			{
			x = (func)()
			for m in x.Members().Remove('fields')
				{
				if info.Member?(m)
					throw 'Cannot define ' $ m $ ' more than once'
				info[m] = x[m]
				}
			if x.Member?('fields')
				{
				info.fields.Append(x.fields)
				}
			}
		if '' isnt excludeFields = info.GetDefault('excludeFields', '')
			info.fields.Remove(@excludeFields)
		return info is #(fields: ()) ? false : info
		}

	getCustomFields()
		{
		return TableExists?(.Table)
			? Customizable.CustomFields(.Table)
			: #()
		}

	GetDepTables()
		{
		tables = Object()
		tables.Add(@.GetKeys().
			Filter({ Object?(it) and it.Member?('in') }).
			Map({ it.in }))
		tables.Add(@.GetForeignKeys().Members())
		return tables
		}

	UserInterfaceColumns()
		{
		cols = QueryAvailableColumns(.Table).Difference(.UserInterfaceExclude())
		contribs = GetContributions('Table_' $ .Table)
		removeOb = Customizable.GetNonPermissableFields(.Table)

		usrExcludes = Object()
		if contribs.Member?('UserInterfaceExclude')
			{
			for excl in contribs['UserInterfaceExclude']
				if String?(excl)
					usrExcludes.AddUnique(excl)
				else
					usrExcludes.MergeUnion(excl())
			}

		cols = cols.Difference(usrExcludes.MergeUnion(removeOb))
		return cols.RemoveIf(Internal?)
		}

	UserInterfaceExclude() { return #() }

	BaseQueryWithExcludedCols()
		{
		return .Table $ ' project ' $ .UserInterfaceColumns().Join(',')
		}
	}
