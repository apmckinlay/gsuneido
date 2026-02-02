// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	FieldControl: 'KeyField'
	New(query, .field, columns = false, .fillin = false, from = false,
		mandatory = false, allowOther = false, .whereField = false,
		width = 10, .readonly = false, noClear = false,
		.restrictions = false, .nameField = false , abbrevField = false,
		.access = false, .prefixColumn = false, keys = false,
		restrictionsValid = true, .filterOnEmpty = true, style = 0,
		.excludeSelect = #(), status = '', .invalidRestrictions = false,
		title = '', whereFieldKey = false, font = "", size = "", weight = "",
		tabover = false, .saveInfoName = "", hidden = false, .customizeQueryCols = false,
		upper = false, lower = false, .optionalRestrictions = #(), .startLast = false)
		{
		super(Object(.FieldControl, name: "Value",
				query: query = .queryFromFunction(query),
				:field, :fillin, :from, :mandatory, :allowOther, :abbrevField, :width,
				:readonly, :noClear, :restrictions, :nameField, :restrictionsValid,
				:style, :status, :invalidRestrictions, :whereFieldKey, :font, :size,
				:weight, :tabover, addAccessOption: access isnt false, :hidden,
				:upper, :lower, :optionalRestrictions))
		.Title = title
		.original_query = query
		.query = .Key_BuildQuery(query, .restrictions, .invalidRestrictions)
		.columns = columns is false ? columns : columns.Copy()
		if String?(.excludeSelect)
			.excludeSelect = Global(.excludeSelect)()
		.locatekeys = keys
		super.SetReadOnly(readonly)
		.Setup_KeyControl()
		}
	queryFromFunction(query)
		{
		return Function?(query) ? query() : query
		}
	Key_BuildQuery(query, restrictions, invalidRestrictions, noSend = false,
		optRestrictions = #())
		{
		query = .queryFromFunction(query)
		if .noRestrictions?(restrictions, invalidRestrictions, optRestrictions)
			return query
		// build  and add restrictions
		restrictions = .build_where_clause(restrictions, noSend)
		invalidRestrictions = .build_where_clause(invalidRestrictions, noSend)
		if restrictions isnt ''
			query = QueryAddWhere(query, "where " $ restrictions)
		if invalidRestrictions isnt ''
			query = QueryAddWhere(query, "where " $ invalidRestrictions)

		query = .handleOptionalRestrictions(query, optRestrictions)

		return query
		}
	noRestrictions?(restrictions, invalidRestrictions, optRestrictions)
		{
		return restrictions is false and invalidRestrictions is false and
			optRestrictions.Empty?()
		}
	handleOptionalRestrictions(query, optRestrictions)
		{
		if optRestrictions.Empty?()
			return query

		for field in optRestrictions
			{
			dd = Datadict(field)
			if KeyListViewBase.ChooseDateCtrl?(dd)
				query = QueryAddWhere(query, 'where ' $
					KeyListViewBase.BuildChooseDateQueryWhere(field))
			else if KeyListViewBase.ChooseListActiveInactive?(dd)
				query = QueryAddWhere(query, 'where ' $
					KeyListViewBase.BuildChooseListActiveInactiveWhere(field))
			}

		return query
		}
	build_where_clause(restrictions, noSend = false)
		{
		if restrictions is false
			return ""

		// detect if just a field name was specified and send GetField.
		// Otherwise restrictions are assumed to be valid restriction expressions
		// for the query
		if restrictions =~ "^[a-zA-Z0-9_]+[?!]?$"
			{
			if noSend
				return ''
			result = .Send('GetField', restrictions)
			if result is 0
				{
				SuneidoLog(
					"ERROR: KeyControl couldn't get restriction field value for: " $
					restrictions)
				result = ''
				}
			return result
			}
		return restrictions
		}
	Setup_KeyControl()
		{
		.keys = .Send("GetKeys")
		if (Object?(.keys) and .keys.Find(.field) isnt false)
			{
			.RemoveButton()
			.Field.Setmode("unique")
			}
		}
	SetTitle(title)
		{
		.Title = title
		}

	Getter_DialogControl()
		{
		.Field_SetFocus()  // refreshes the whereField
		.valid? = .Field.Valid?()
		.prefix = GetWindowText(.Field.Hwnd)
		saveInfoName = .GetDropDownKeepSizeName()
		if saveInfoName is ''
			SuneidoLog("INFO: KeyControl doesn't have Name for KeyListView",
				calls:, params: Object(field: .field, query: .query))
		return Object(KeyListView, .query, .columns, saveInfoName, .prefix,
			access: .access, prefixColumn: .prefixColumn, keys: .locatekeys,
			field: .field, value: .Get(), excludeSelect: .excludeSelect,
			customizeQueryCols: .customizeQueryCols,
			optionalRestrictions: .optionalRestrictions, closeButton?:,
			startLast: .startLast)
		}

	GetDropDownKeepSizeName()
		{
		return .saveInfoName isnt '' ? .saveInfoName : .Name
		}

	ProcessResults(result)
		{
		val = .Convert(result[0][.field])
		.Set(val)
		.Field.Process_newvalue()	// 021024 apm,khd moved this before newvalue
									// so fillins are done before newvalue
			// needed this for ETA Orders shipper/consignee creating picks/drops
		.NewValue(.Get())
		.Field.Last_ellipsis = .Get()
		.Send("ListRecordSelected", result[0], result[1])
		}

	ReprocessValue()
		{
		.reProcessValue(.valid?, .prefix)
		}

	reProcessValue(valid?, fieldText)
		{
		// user may have entered previously invalid value into master table while in list.
		// Need to make sure value is processed again
		// and the timestamp value is sent as the new value.
		if not valid? and fieldText isnt ''
			{
			.Field.Process_newvalue()
			.NewValue(.Get())
			}
		if valid? and fieldText isnt ''
			.Send('KeyControl_ReprocessValue')
		}

	Convert(x) { return x }

	ChangeKey(newkey)
		{
		.prefixColumn = newkey
		.field = newkey
		.Field.ChangeKey(.field)

		// make sure key is first column
		if .columns isnt false
			.columns.Remove(newkey).Add(newkey, at: 0)
		}
	// methods for handling whereField
	GetWhereField()
		{
		if (.whereField isnt false and
			.filterOnEmpty isnt true and '' is .Send('GetField', .whereField))
			return false

		return .whereField
		}
	Field_SetFocus()
		{
		.Send('Field_SetFocus')
		// optRestrictions should not be passed in when building query here because
		// the KeyListView from drop-down automatically takes that into consideration
		.query = .Key_BuildQuery(.original_query, .restrictions, .invalidRestrictions)
		if (.whereField isnt false)
			{
			wherevalue = .Send('GetField', .whereField)
			where = ' where ' $ .whereField $ ' is ' $ Display(wherevalue)
			if (.filterOnEmpty isnt true and wherevalue is "")
				where = ""

			.query = QueryAddWhere(.query, where)
			if (not Object?(.fillin) or not .fillin.Has?(.whereField))
				.Field.Setquery(.query)
			}
		}
	FieldReturn()
		{
		dirty? = .Dirty?()
		.Field.KillFocus()
		if not dirty?
			return

		if valid? = .Valid?()
			.NewValue(.Get())
		if .Destroyed?()
			return
		.Field.SetValid(valid?, force:)
		}
	SetReadOnly(readOnly)
		{
		if .readonly
			return
		super.SetReadOnly(readOnly)
		}

	SetDefaultReadOnly(readonly, controllerReadOnly)
		{
		.readonly = readonly
		.Field.SetDefaultReadOnly(readonly, controllerReadOnly)
		.Button.SetReadOnly(controllerReadOnly or readonly)
		}

	// Public method for getting control info
	// specifically used by multi-select option for 'In List'
	Key_GetControlInfo()
		{
		if .Field.Method?("Key_DisplayField")
			displayField = .Field.Key_DisplayField()
		else
			displayField = .nameField isnt false ? .nameField : .field
		return Object(query: .query,
			columns: Object?(.columns)
				? .columns.Copy() : QueryColumns(.query),
			:displayField, field: .field)
		}

	SelectAll()
		{
		.Field.SelectAll()
		}
	KeyIdField_Access()
		{
		valid? = .Field.Valid?() // have to check valid before access
		AccessGoTo(.access, .field, .Field.Get(), .Window.Hwnd,
			onDestroy: {
				if not .GetReadOnly()
					.reProcessValue(valid?, GetWindowText(.Field.Hwnd))
				SetFocus(.Field.Hwnd)
				})
		}
	SetListRestrictions(restrictions)
		{
		.restrictions = restrictions
		.Field.SetFieldRestrictions(restrictions)
		}
	ValidData?(@args)
		{
		return GetControlClass.FromControlName(.FieldControl).ValidData?(@args)
		}
	// Helper method so customizable can use fillin logic
	FillinRecord(control, field, rec, fillincustom = false)
		{
		// using number method in case the KeyControl options aren't named
		query = .queryFromFunction(
			control.GetDefault('query', control.GetDefault(1, false)))
		valfield = control.GetDefault('field', control.GetDefault(2, false))
		fillin = control.GetDefault('fillin', control.GetDefault(4, false)) /*= arg pos */
		from = control.GetDefault('from', control.GetDefault(5, fillin))    /*= arg pos */

		if query is false or field is false
			return

		KeyFieldControl.FillinRecord(
			query, valfield, fillin, from, rec, rec[field], fillincustom)
		}

	CustomizableSetDefaultValue(x, custom, dd)
		{
		custom[x.custfield_field] = x.custfield_default_value
		.FillinRecord(dd.Control, x.custfield_field, custom)
		return true
		}

	IsKeyControl?(field)
		{
		control = Datadict(field).Control
		return control[0] is 'Key' or
			GetControlClass.FromControl(control).Base?(KeyControl)
		}
	}
