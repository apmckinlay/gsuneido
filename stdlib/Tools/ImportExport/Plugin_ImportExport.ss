// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
#(
ExtensionPoints:
	(
	('import_formats'),
	('export_formats')
	)
Contributions:
	(
	(ImportExport, export_formats,
		name: "Comma separated", impl: ExportCSV)
	(ImportExport, export_formats,
		name: "Tab delimited", impl: ExportTab)
	(ImportExport, export_formats,
		name: "XML", impl: ExportXML)

	(ImportExport, import_formats,
		name: "Comma separated", impl: ImportCSV)
	(ImportExport, import_formats,
		name: "Tab delimited", impl: ImportTab)
	(ImportExport, import_formats,
		name: "XML", impl: ImportXML)
	)
)