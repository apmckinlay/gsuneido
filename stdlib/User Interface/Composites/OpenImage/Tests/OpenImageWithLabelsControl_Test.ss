// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_CleanupLabels()
		{
		f = OpenImageWithLabelsControl.CleanupLabels
		Assert(f('') is: '')
		Assert(f('  , ,') is: '')
		Assert(f('web') is: 'web')
		Assert(f('web, web') is: 'web')
		Assert(f('one,two,  three,four , five') is: 'five, four, one, three, two')
		Assert(f('a, ,b') is: 'a, b')
		}
	Test_SplitLabel()
		{
		imageLabel = OpenImageWithLabelsControl
			{
			LabelDelimiter: 'Label:'
			}

		Assert(imageLabel.SplitLabel("")
			is: #(file: "", labels: '', subfolder: ''))

		Assert(imageLabel.SplitLabel("FileName")
			is: #(file: 'FileName', labels: '', subfolder: ''))

		Assert(imageLabel.SplitLabel('201009\\FileName Label:LabelValue')
			is: #(file: 'FileName', labels: 'LabelValue', subfolder: '201009'))

		Assert(imageLabel.SplitLabel('Steve\\FileName')
			is: #(file: 'FileName', labels: '', subfolder: 'Steve'))

		Assert(imageLabel.SplitLabel('\\FileName')
			is: #(file: 'FileName', labels: '', subfolder: ''))

		Assert(imageLabel.SplitLabel('C:\\FileName Label:LabelValue')
			is: #(file: 'C:\\FileName', labels: 'LabelValue', subfolder: ''))

		Assert(imageLabel.SplitLabel('c:\\Folder Space\\FileName Label:LabelValue')
			is: #(file: 'c:\\Folder Space\\FileName', labels: 'LabelValue', subfolder: '')
			)

		Assert(imageLabel.SplitLabel('201009\\NotSub\\FileName Label:LabelValue')
			is: #(file: 'FileName', labels: 'LabelValue', subfolder: '201009\\NotSub'))

		Assert(imageLabel.SplitLabel('\\\\server\\FileName Label:LabelValue')
			is: #(file: '\\\\server\\FileName', labels: 'LabelValue', subfolder: ''))

		Assert(imageLabel.SplitLabel('\\\\server\\path\\Staff\\FileName Label:LabelValue')
			is: #(file: '\\\\server\\path\\Staff\\FileName', labels: 'LabelValue',
				subfolder: ''))

		Assert(imageLabel.SplitLabel('Staff/FileName Label:LabelValue')
			is: #(file: 'FileName', labels: 'LabelValue', subfolder: 'Staff'))
		}
	Test_SplitFile()
		{
		imageLabel = OpenImageWithLabelsControl
			{
			LabelDelimiter: 'Label:'
			}
		fn = imageLabel.SplitFile
		Assert(fn(`\\server\path\Staff\FileName Label:LabelValue`, 'LabelValue')
			is: 'FileName')
		Assert(fn(
			`\\server\path\Staff\FileName Label:LabelValue, LabelValue2`, 'LabelValue')
			is: 'FileName')
		Assert(fn(
			`\\server\path\Staff\FileName Label:LabelValue, LabelValue2`, 'LabelValue2')
			is: 'FileName')
		Assert(fn(
			`\\server\path\Staff\FileName Label:LabelValue,LabelValue2`, 'LabelValue')
			is: 'FileName')
		Assert(fn(
			`\\server\path\Staff\FileName Label:LabelValue,LabelValue2`, 'LabelValue2')
			is: 'FileName')
		Assert(fn(
			`\\server\path\Staff\FileName Label:LabelValue,LabelValue2`, 'LabelValue3')
			is: '')
		Assert(fn(`\\server\path\Staff\FileName Label:LabelValue,LabelValue2`, false)
			is: 'FileName')
		}

	Test_FullPath()
		{
		mock = Mock(OpenImageWithLabelsControl)
		mock.When.FullPath([anyArgs:]).CallThrough()

		// Test: copyTo is "", file is returned
		mock.When.getCopyTo([anyArgs:]).Return('')
		Assert(mock.FullPath('file.txt', 'subfolder') is: 'file.txt')

		// Test: copyTo has a value, built path is returned
		mock.When.getCopyTo([anyArgs:]).Return('copyTo/')
		Assert(mock.FullPath('file.txt', '/subfolder/') is: 'copyTo/subfolder/file.txt')

		// Test: subfolder is "", built path is returned
		mock.When.getCopyTo([anyArgs:]).Return('copyTo')
		Assert(mock.FullPath('file.txt', '') is: 'copyTo/file.txt')

		// Test: file has path components, file is returned
		mock.When.getCopyTo([anyArgs:]).Return('copyTo/')
		Assert(mock.FullPath('folder/file.txt', 'subfolder/') is: 'folder/file.txt')

		// Test: file is "", file is returned
		mock.When.getCopyTo([anyArgs:]).Return('copyTo/')
		Assert(mock.FullPath('', 'subfolder') is: '')

		// Test: default arguments, built path is returned
		mock.File = 'file.txt'
		mock.OpenImageWithLabelsControl_subfolder = 'subfolder'
		Assert(mock.FullPath() is: 'copyTo/subfolder/file.txt')
		}
	}
