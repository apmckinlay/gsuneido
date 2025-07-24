// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_checkSettings()
		{
		cl = Clscan
			{
			Clscan_execute(cmd)
				{
				Assert(cmd is: _expectedCheckScanCmd)
				return true
				}
			Clscan_getScannerSettingsOb(unused)
				{ return _getScannerSettingsReturn }
			Clscan_compareSettings(@unused)
				{ return _comareSettingsReturn }
			}
		.SpyOn(UserSettings.Get).Return(#())
		.SpyOn(UserSettings.Put).Return(#())
		_comareSettingsReturn = ''
		_getScannerSettingsReturn = true
		func = cl.Clscan_checkSettings
		Assert(func('', '', #()) is:  '')
		Assert(func('', 'Scanner', #()) is:  '')
		// No settings just basic scan call
		Assert(func('Scan', 'Scanner', #()) is:  'Scan')
		// run check against scanner to get final command
		scanSettings = #(#(check: 'GetColors' scan: 'SetColors' value: 'black'))
		_expectedCheckScanCmd = 'Scanner GetColors'
		Assert(func('ScanToFile', 'Scanner', scanSettings)
			is: 'ScanToFile SetColors "black"')
		// run check with all settings to get final command
		scanSettings = #(#(check: 'GetColors' scan: 'SetColors' value: 'black'),
			#(check: 'GetPageSize', scan: 'SetPageSize' value: 'LETTER'),
			#(check: 'GetResolution', scan: 'SetResolution' value: '200'))
		_expectedCheckScanCmd = 'Scanner GetColors GetPageSize GetResolution'
		Assert(func('ScanToFile', 'Scanner', scanSettings)
			is: 'ScanToFile SetColors "black" SetPageSize "LETTER" SetResolution "200"')

		//fetching settings failed
		_getScannerSettingsReturn = false
		Assert(func('ScanToFile', 'Scanner', scanSettings)
			is: 'Problem detecting scanner')

		_getScannerSettingsReturn = true
		_comareSettingsReturn = 'Resolution doesnt match scanner'
		spy = .SpyOn(YesNo)
		spy.Return(false)
		Assert(func('ScanToFile', 'Scanner', scanSettings)	is: 'User Cancelled')
		}

	Test_getScannerSettingsOb()
		{
		scanSettingsCheckString = ''
		getSettingsOb = Clscan.Clscan_getScannerSettingsOb
		Assert(getSettingsOb(scanSettingsCheckString) is: #())

		scanSettingsCheckString = 'Selected scanner HP Photosmart 5520 TWAIN\r\n' $
			'Supported page sizes:\r\n' $
			'NONE\r\n' $
			'USLETTER\r\n' $
			'USEXECUTIVE\r\n' $
			'USSTATEMENT\r\n' $
			'BUSINESSCARD'
		Assert(getSettingsOb(scanSettingsCheckString)
			is: #('page sizes': #("NONE", "USLETTER", "USEXECUTIVE", "USSTATEMENT",
					"BUSINESSCARD")))

		scanSettingsCheckString = 'Selected scanner HP Photosmart 5520 TWAIN\r\n' $
			'Supported resolutions:\r\n' $
			'75\r\n' $
			'100\r\n' $
			'200\r\n' $
			'300\r\n' $
			'600\r\n' $
			'1200\r\n' $
			'2400\r\n' $
			'Supported page sizes:\r\n' $
			'NONE\r\n' $
			'USLETTER\r\n' $
			'USEXECUTIVE\r\n' $
			'USSTATEMENT\r\n' $
			'BUSINESSCARD\r\n' $
			'Supported color types:\r\n' $
			'RGB\r\n' $
			'GRAY\r\n' $
			'BW'
		getSettingsOb = Clscan.Clscan_getScannerSettingsOb
		Assert(getSettingsOb(scanSettingsCheckString)
			is: #('resolutions': #("75", "100", "200", "300", "600", "1200", "2400"),
				'page sizes': #("NONE", "USLETTER", "USEXECUTIVE", "USSTATEMENT",
					"BUSINESSCARD"),
				'color types':#('RGB', 'GRAY', 'BW')))

		scanSettingsCheckString = 'Selected scanner HP Photosmart 5520 TWAIN\r\n' $
			'Supported duplex:\r\n' $
			'False'
		Assert(getSettingsOb(scanSettingsCheckString) is: #('duplex': false))

		scanSettingsCheckString = 'Selected scanner HP Photosmart 5520 TWAIN\r\n' $
			'Supported resolutions:\r\n' $
			'min value=50, max value=1200, step=1'
		Assert(getSettingsOb(scanSettingsCheckString)
			is: #(resolutions: #("min value=50, max value=1200, step=1")))

		untestedVersionString = 'Qt: Untested Windows version 6.2 detected!\r\n' $
			'Selected scanner HP Photosmart 5520 TWAIN\r\n' $
			'Supported resolutions:\r\n' $
			'75\r\n' $
			'200\r\n' $
			'300\r\n' $
			'Supported color types:\r\n' $
			'RGB\r\n' $
			'GRAY\r\n' $
			'BW\r\n' $
			'Supported duplex:\r\n' $
			'False'
		Assert(getSettingsOb(untestedVersionString) is:
			#(resolutions: #('75', '200', '300'), 'color types': #('RGB', 'GRAY', 'BW'),
				duplex: false))
		}


	expectedFormatedOb: #("ColorType = RGB",
		"Filename = C:/scanning/aScannedDocument.pdf",
		"Resolution = 75 DPI per inches",
		"Duplex enabled = False",
		"Contrast = 0",
		"Brightness = 0",
		"JpegQuality = 85",
		"Threshold = 128",
		"Rotation = 0",
		"Orientation = 0",
		"LogToFile =",
		"MultiFile = False",
		"Deskew = False",
		"Crop = False",
		"Invert = False",
		"",
		"Scan Complete")
	Test_formatResults_makeOptionOb()
		{
		scanResultString = 'Selected scanner HP Photosmart 5520 TWAIN\r\n' $
			'Scanner settings:\r\n' $
			'ColorType = RGB\r\n' $
			'Filename = C:/scanning/aScannedDocument.pdf\r\n' $
			'Resolution = 75 DPI per inches\r\n' $
			'Duplex enabled = False\r\n' $
			'Contrast = 0\r\n' $
			'Brightness = 0\r\n' $
			'JpegQuality = 85\r\n' $
			'Threshold = 128\r\n' $
			'Rotation = 0\r\n' $
			'Orientation = 0\r\n' $
			'LogToFile =\r\n' $
			'MultiFile = False\r\n' $
			'Deskew = False\r\n' $
			'Crop = False\r\n' $
			'Invert = False\r\n\r\n' $
			'Scan Complete'
		formatted = Clscan.Clscan_formatResults(scanResultString)
		Assert(formatted is: .expectedFormatedOb)

		Assert(Clscan.Clscan_makeOptionOb(formatted)
			is: #("ColorType": 'RGB',
				'Filename': 'C:/scanning/aScannedDocument.pdf',
				'Resolution': '75 DPI per inches',
				'Duplex enabled': 'False',
				'Contrast': '0',
				'Brightness': '0',
				'JpegQuality': '85',
				'Threshold': '128',
				'Rotation': '0',
				'Orientation': '0',
				'LogToFile': '',
				'MultiFile': 'False',
				'Deskew': 'False',
				'Crop': 'False',
				'Invert': 'False'
				'results': 'Scan Complete'))
		}

	Test_compareSettings()
		{
		compareSettings = Clscan.Clscan_compareSettings
		Assert(compareSettings(#(), #()) is: '')

		companyScanSettings = #(resolutions: #(value: 200))
		scannerSettings = #(resolutions: #('100', '200'))
		Assert(compareSettings(scannerSettings, companyScanSettings) is: '')

		scannerSettings = #(resolutions: #('100'))
		Assert(compareSettings(scannerSettings, companyScanSettings)
			is: 'The resolution specified is not supported by your scanner.\r\n')

		companyScanSettings = #(resolutions: #(value: 200), duplex: #(value: 'Y'))
		scannerSettings = #(resolutions: #('100'), duplex: false)
		Assert(compareSettings(scannerSettings, companyScanSettings)
			is: 'Scanner does not support changing double-sided options.\r\n' $
				'Depending on scanner settings, both or one side may be scanned.\r\n' $
				'The resolution specified is not supported by your scanner.\r\n')

		// not using duplex so wont show up in scannerSettings
		companyScanSettings = #(resolutions: #(value: 100), duplex: #(value: 'N'))
		scannerSettings = #(resolutions: #('100'))
		Assert(compareSettings(scannerSettings, companyScanSettings) is: '')

		companyScanSettings = #(resolutions: #(value: 100), duplex: #(value: 'Y'))
		scannerSettings = #(resolutions: #('100'), duplex: true)
		Assert(compareSettings(scannerSettings, companyScanSettings) is: '')

		companyScanSettings = #(resolutions: #(value: 25))
		scannerSettings = #(resolutions: #("min value=50, max value=1200, step=5"))
		Assert(compareSettings(scannerSettings, companyScanSettings)
			is: 'The resolution specified is not supported by your scanner.\r\n')

		companyScanSettings = #(resolutions: #(value: 25))
		scannerSettings = #(resolutions: #("Min value=50, max value=1200, step=5"))
		Assert(compareSettings(scannerSettings, companyScanSettings)
			is: 'The resolution specified is not supported by your scanner.\r\n')

		companyScanSettings = #(resolutions: #(value: 1250))
		Assert(compareSettings(scannerSettings, companyScanSettings)
			is: 'The resolution specified is not supported by your scanner.\r\n')

		companyScanSettings = #(resolutions: #(value: 375))
		Assert(compareSettings(scannerSettings, companyScanSettings) is: '')
		}

	Test_handleScanResults()
		{
		spy = .SpyOn(BookLog)
		fn = Clscan.Clscan_handleScanResults
		filename = 'testfile'
		scanResults = 'Selected scanner fred'
		Assert(fn(filename, scanResults), msg: 'scanner fred')
		Assert(spy.CallLogs()[0].s is: "Scan Attachment - Clscan Scanner: fred")
		Assert(spy.CallLogs()[1].s is: "Scan Attachment End - Clscan")

		scanResults = 'testfile is saved.'
		Assert(fn(filename, scanResults), msg: 'testfile is saved')
		Assert(spy.CallLogs()[2].s is: "Scan Attachment End - Clscan")

		scanResults = 'this is sometext\nsome more text\nUnable to scan.'
		Assert(fn(filename, scanResults)
			is: "Please ensure the scanner is turned on and connected.\r\n" $
				"Unable to scan.")
		Assert(spy.CallLogs() isSize: 3)

		scanResults = 'this is line\nthis is more line\n' $
			'Problem while opening the scanner fred, please try again'
		Assert(fn(filename, scanResults)
			is: "Please ensure the scanner is turned on and connected.\r\n" $
				"Problem while opening the scanner fred, please try again")
		Assert(spy.CallLogs() isSize: 3)

		scanResults = ''
		Assert(fn(filename, scanResults), msg: 'scan end')
		Assert(spy.CallLogs()[3].s is: "Scan Attachment End - Clscan")

		scanResults = '\nProblem while opening the scanner: ' $
			'The network location cannot be reached.'
		Assert(fn(filename, scanResults)
			is: "Please ensure the scanner is turned on and connected.\r\n" $
				"Problem while opening the scanner: " $
				"The network location cannot be reached.")
		}

	}
