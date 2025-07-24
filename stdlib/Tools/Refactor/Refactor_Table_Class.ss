// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Refactor
	{
	Name: 'Extract Table'
	Desc: 'Convert the selection to a separate TableModel Class'
	DiffPos: 6
	IdleTime: 250 //ms
	Controls: (Vert
		(Pair (Static 'Name') (Field, name: 'classname'))
		(Pair (Static Path) (Field readonly:, name: newpath))
		(Skip 4)
		#(Editor, name: 'newclass')
		name: 'extractVert'
		xmin: 600
		)
	Init(data)
		{
		if .validate(data) is false
			return false
		ob = false
		try ob = ('#' $ .selection).SafeEval()
		if ob is false
			{
			.Info("Could not create table object")
			return false
			}
		data.classname = "Table_" $ ob.GetDefault("table", "???")
		data.newclass = .buildRecordText(ob)
		path = data.libview.Explorer.Getpath(data.libview.Explorer.GetSelected())
		data.newpath = path.BeforeLast(`/`)
		.data = data
		return true
		}
	validate(data)
		{
		if data.select.cpMin >= data.select.cpMax
			{
			.Info(
			'Please select the table object you want to extract to a TableModel class')
			return false
			}
		.selection = data.text[data.select.cpMin :: data.select.cpMax - data.select.cpMin]
		missing = Object()
		for str in #('Tables', 'ensures', 'table', 'schema')
			if not .selection.Has?(str)
			missing.Add(str)

		if not missing.Empty?()
			{
			.Info('Selection missing: ' $ missing.Join(','))
			return false
			}

		return true
		}
	buildRecordText(ob)
		{
		info = .parseSchema(ob.schema)
		foreignKeys = .buildForeignKeyStr(info.foreignKeys)
		str =  "TableModel
	{
	Table: '" $ ob.GetDefault("table", "???") $ "'
	Name: ''

	Columns: " $ '(' $ ob.schema.AfterFirst('(').BeforeFirst(')') $ ')' $ "
	Keys: (" $ info.keys.Join(', ') $ ")
	Indexes: (" $ info.indexes.Join(', ') $ ")
	" $ foreignKeys $ "

	BookLocation: ''
	Permission: ''
	}

"
		str $= '/*' $ .selection $ '*/'
		return str
		}
	parseSchema(text)
		{
		keys = Object()
		foreignKeys = Object()
		indexes = Object()
		trackingInfo = Object(readingKeys: false, readingIndex: false, addingIn: false,
			current: '', last: '', lasttype: '')
		ScannerEach(text)
			{|prev2 /*unused*/, prev /*unused*/, token, next /*unused*/|
			token = token.Tr('\r\n\t')
			if token is 'key'
				trackingInfo.readingKeys = true
			else if token is 'index'
				trackingInfo.readingIndex = true
			else if token is ')'
				.processCloseParen(trackingInfo, keys, indexes)
			else if token is 'in'
				trackingInfo.addingIn = true
			else if token isnt '(' and token isnt '' and token isnt ' '
				.processNextToken(token, trackingInfo, keys, indexes, foreignKeys)
			}
		return Object(:keys, :foreignKeys, :indexes)
		}
	processCloseParen(trackingInfo, keys, indexes)
		{
		if trackingInfo.readingKeys is true and trackingInfo.current isnt ""
			{
			keys.Add(trackingInfo.current)
			trackingInfo.readingKeys = false
			trackingInfo.lasttype = 'key'
			}
		else if trackingInfo.readingIndex is true and trackingInfo.current isnt ""
			{
			indexes.Add(trackingInfo.current)
			trackingInfo.readingIndex = false
			trackingInfo.lasttype = 'index'
			}
		trackingInfo.last = trackingInfo.current
		trackingInfo.current = ''
		}
	processNextToken(token, trackingInfo, keys, indexes, foreignKeys)
		{
		if trackingInfo.addingIn is true
			.processAddingIn(token, trackingInfo, keys, indexes, foreignKeys)
		else if token is 'cascade'
			foreignKeys[trackingInfo.last] = Object(
				foreignKeys[trackingInfo.last], cascade:)
		else
			trackingInfo.current $= token
		}
	processAddingIn(token, trackingInfo, keys, indexes, foreignKeys)
		{
		if trackingInfo.lasttype is 'key'
			{
			keys.Remove(trackingInfo.last)
			foreignKeys.Add(trackingInfo.last, at: token)
			trackingInfo.last = token
			trackingInfo.readingKeys = false
			}
		else if trackingInfo.lasttype is 'index'
			{
			indexes.Remove(trackingInfo.last)
			foreignKeys.Add(trackingInfo.last, at: token)
			trackingInfo.last = token
			trackingInfo.readingIndex = false
			}
		trackingInfo.current = ''
		trackingInfo.addingIn = false
		}

	buildForeignKeyStr(foreignKeys)
		{
		str = ''
		if foreignKeys.Empty?()
			return str
		str = 'ForeignKeys: ('
		for mem in foreignKeys.Members().Sort!()
			{
			str $= mem $ ': ((from: '
			if Object?(foreignKeys[mem])
				{
				str $= foreignKeys[mem][0]
				if foreignKeys[mem].Member?('cascade')
					str $= ', cascade:'
				}
			else
				str $= foreignKeys[mem]
			str $= '))\r\n\t'
			}
		str $= ')'
		return str
		}

	Process(data)
		{
		data.libview.Explorer.NewItem(false, data.classname, data.newclass)
		data.text = data.text.ReplaceSubstr(data.select.cpMin, .selection.Size(), '')
		return true
		}
	}