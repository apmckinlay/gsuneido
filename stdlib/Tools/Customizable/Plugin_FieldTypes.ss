// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
#(
ExtensionPoints:
	(
	("type"),
	)
Contributions:
	(
		(FieldTypes, type, name: 'Checkmark', base: boolean_yesno, customize:, reporter:)
		(FieldTypes, type, name: 'Number', base: number_custom, customize:)
		(FieldTypes, type, name: 'Dollar', base: dollar, customize:, reporter:)
		(FieldTypes, type, name: "Text, single line", base: string_custom,
			customize:, reporter:,
			compatible: function () { Suneido.User is 'default' ? 'string' : false })
		(FieldTypes, type, name: "Text, from custom table", customize:,
			base: string_custom_table, customOptions: CustomFieldControl_CustomTable,
			compatible: 'list', oneWay:)
		(FieldTypes, type, name: 'Text, multi line', base: string_multiline_custom,
			customize:, reporter:,
			compatible: function () { Suneido.User is 'default' ? 'string' : false })
		(FieldTypes, type, name: 'Info', base: info, customize:)
		(FieldTypes, type, name: 'Zip/Postal', base: zip_postal, customize:)
		(FieldTypes, type, name: 'State/Province', base: state_prov, customize:)
		(FieldTypes, type, name: 'Date', base: date_custom, customize:)
		(FieldTypes, type, name: 'Date and Time', base: date_time, customize:, reporter:)
		(FieldTypes, type, name: 'Choose Several Dates', base: multi_dates, customize:)
		(FieldTypes, type, name: 'Year and Month', base: year_month, customize:)
		(FieldTypes, type, name: 'Time', base: time, customize:)
		(FieldTypes, type, name: 'Choose List',
			base: string_chooselist_custom, customize:, compatible: 'list')
		(FieldTypes, type, name: 'Choose Several from List',
			base: string_choosemany_custom,	customize:)
		(FieldTypes, type, name: 'Attachment', base: attachment, customize:)
		(FieldTypes, type, name: 'Number, no decimals',
			base: number_no_decimals_custom, reporter:)
		(FieldTypes, type, name: 'Number, 1 decimal',
			base: number_one_decimals_custom, reporter:)
		(FieldTypes, type, name: 'Number, 2 decimals',
			base: number_two_decimals_custom, reporter:)
		(FieldTypes, type, name: 'Number, 3 decimals',
			base: number_three_decimals_custom, reporter:)
		(FieldTypes, type, name: 'Number, 4 decimals',
			base: number_four_decimals_custom, reporter:)
		(FieldTypes, type, name: 'Short Date', base: short_date_custom, reporter:)
		(FieldTypes, type, name: 'Long Date', base: long_date_custom, reporter:)
	)
)
