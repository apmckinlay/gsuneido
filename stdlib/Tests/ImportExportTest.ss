// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ExportCSV()
		{
		tableName = .MakeTable('(one, two, three, four, five) key(three)')
		Transaction(update:)
			{ |trans|
			query = trans.Query(tableName)
			query.Output(
				#(one: 'one', two: 'two', three: 3, four: 'f,o,u,r', five: '"five"'))
			}
		ExportCSV(tableName, tableName, header:)
		File(tableName, 'r')
			{ |file|
			line = file.Readline()
			Assert(line is: '"one","two","three","four","five"')
			line = file.Readline()
			Assert(line is: '"one","two",3,"f,o,u,r","""five"""')
			}
		}
	Test_ExportXML()
		{
		date = #20011021.192837
		tableName = .MakeTable('(export_fname, export_bdate, export_age,
			export_status, export_code)	key( export_fname )')
		Transaction(update:)
			{ |trans|
			query = trans.Query(tableName)
			query.Output(#(export_fname: "Bob", export_bdate:, export_age: 35,
				export_status:, export_code: (font)))
			query.Output(Object(export_fname: "Adam", export_bdate: date.ShortDateTime(),
				export_status: false, export_code: #(font)))
			}

		fileName = .MakeFile()
		numRec = ExportXML(tableName, fileName, header:)

		Assert(numRec is: 2)
		expected_lines = Object('<?xml version="1.0"?>',
			'<!--  suneido xml export  -->'
			'<table>'
			'<record>'
			'<export_fname type="String">Adam</export_fname>'
			'<export_bdate type="String">' $ date.ShortDateTime() $ '</export_bdate>'
			'<export_status type="Boolean">false</export_status>'
			'<export_code type="Object">#(&quot;font&quot;)</export_code>'
			'</record>'
			'<record>'
			'<export_fname type="String">Bob</export_fname>'
			'<export_bdate type="Boolean">true</export_bdate>'
			'<export_age type="Number">35</export_age>'
			'<export_status type="Boolean">true</export_status>'
			'<export_code type="Object">#(&quot;font&quot;)</export_code>'
			'</record>'
			'</table>'
			)
		lines = GetFile(fileName).Lines()
		Assert(lines isSize: expected_lines.Size())
		for expected_line in expected_lines
			Assert(lines has: expected_line)
		}

	Test_mapToXmlName()
		{
		x = ExportXML.ExportXML_mapToXmlName('')
		Assert(x is: ':A:')
		x = ExportXML.ExportXML_mapToXmlName('hello')
		Assert(x is: 'hello')
		x = ExportXML.ExportXML_mapToXmlName(':_azAZ09.-')
		Assert(x is: ':_azAZ09.-')
		x = ExportXML.ExportXML_mapToXmlName('?  \\= what')
		Assert(x is: ':A:0x3f0x200x200x5c0x3d0x20what')
		x = ExportXML.ExportXML_mapToXmlName('Sample prompt looks like this?')
		Assert(x is: 'Sample0x20prompt0x20looks0x20like0x20this0x3f')
		x = ExportXML.ExportXML_mapToXmlName('123')
		Assert(x is:':A:123')
		}

	Test_ImportCSV()
		{
		txt = '"import_1","import_2","import_3","import_4","import_5","import_6"\r\n' $
			'"one","two",3,"f,o,u,r","""five""","03"'
		fileName = .MakeFile(txt)
		tableName = .MakeTable('(import_1, import_2, import_3,
			import_4, import_5, import_6) key(import_3)')

		ImportCSV(fileName, tableName, header:)

		Transaction(read:)
			{ |trans|
			query = trans.Query(tableName)
			record = query.Next()
			Assert(record isnt: false)
			Assert(record.import_1 is: 'one')
			Assert(record.import_2 is: 'two')
			Assert(record.import_3 is: 3)
			Assert(record.import_4 is: 'f,o,u,r')
			Assert(record.import_5 is: '"five"')
			Assert(record.import_6 is: "03")
			record = query.Next()
			Assert(record is: false)
			}
		}
	Test_ImportXML()
		{
		fileText = Object()
		fileText.Add('<?xml version="1.0"?>')
		fileText.Add('<!--  suneido xml export  -->')
		fileText.Add('<table>')
		fileText.Add('<record>')
		fileText.Add('<import_fname type="String">Charles</import_fname>')
		fileText.Add('<import_bdate type="Date">#20011021.192837</import_bdate>')
		fileText.Add('<import_age type="Number">29</import_age>')
		fileText.Add('<import_married type="Boolean">true</import_married>')
		fileText.Add('</record>')
		fileText.Add('<record>')
		fileText.Add('<import_fname type="String">Maria</import_fname>')
		fileText.Add('<import_bdate type="Date">#20010605</import_bdate>')
		fileText.Add('<import_code type="Object">#(some, stuff)</import_code>')
		fileText.Add('<import_married type="Boolean">false</import_married>')
		fileText.Add('</record>')
		fileText.Add('<record>')
		fileText.Add('<import_fname type="String">Stan</import_fname>')
		fileText.Add('<import_bdate type="Date">#20011010</import_bdate>')
		fileText.Add('<import_code type="Object">#()</import_code>')
		fileText.Add('<import_married type="Boolean">true</import_married>')
		fileText.Add('</record>')
		fileText.Add('<record>')
		fileText.Add('<import_fname type="String">George</import_fname>')
		fileText.Add('<import_married type="Boolean">no</import_married>')
		fileText.Add('</record>')
		fileText.Add('<record>')
		fileText.Add('<import_fname type="String">Karen</import_fname>')
		fileText.Add('<import_married type="Boolean"></import_married>')
		fileText.Add('</record>')
		fileText.Add('</table>')
		fileName = .MakeFile(fileText.Join('\r\n'))
		tableName = .MakeTable('(import_fname, import_bdate, import_age,
			import_code, import_married) key(import_fname)')

		numRec = ImportXML(fileName, tableName, header:)

		Assert(numRec is: 5)
		Transaction(read:)
			{ |trans|
			query = trans.Query(tableName $ " sort import_fname")
			record = query.Next()
			Assert(record isnt: false)
			Assert(record.import_fname is: "Charles")
			Assert(record.import_bdate is: #20011021.192837)
			Assert(record.import_age is: 29)
			Assert(record.import_married)
			record = query.Next()
			Assert(record isnt: false)
			Assert(record.import_fname is: "George")
			Assert(record.import_married is: false)
			record = query.Next()
			Assert(record isnt: false)
			Assert(record.import_fname is: "Karen")
			Assert(record.import_married is: false)
			record = query.Next()
			Assert(record isnt: false)
			Assert(record.import_fname is: "Maria")
			Assert(record.import_bdate is: #20010605)
			Assert(record.import_code is: #(some, stuff))
			Assert(record.import_married is: false)
			record = query.Next()
			Assert(record isnt: false)
			Assert(record.import_fname is: 'Stan')
			Assert(record.import_code is: #())
			Assert(record.import_married)
			}
		}
	Test_duplicate()
		{
		fileText = Object()
		fileText.Add('"dup_name","dup_abbrev","dup_other"')
		fileText.Add('"one","o",1')
		fileText.Add('"one","t",2')
		fileText.Add('"two","t",3')
		fileText.Add('"two","t",4')
		fileName = .MakeFile(fileText.Join('\r\n'))
		tableName = .MakeTable('(dup_name, dup_abbrev, dup_other)
			key(dup_name) index unique(dup_abbrev)')
		ImportCSV(fileName, tableName, header:)
		i = 0
		data = #((one, o), ('one*1', t), (two, 't*1'), ('two*1', 't*2'))
		QueryApply(tableName $ ' sort dup_name, dup_abbrev')
			{ |x|
			Assert(x.dup_name is: data[i][0])
			Assert(x.dup_abbrev is: data[i][1])
			++i
			}
		Assert(i is: data.Size())
		}
	}
