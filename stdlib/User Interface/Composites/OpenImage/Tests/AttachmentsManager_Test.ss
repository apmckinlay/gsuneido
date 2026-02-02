// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ProcessQueue()
		{
		// deleting old files
		mock = .setupMock()
		mock.ProcessQueue()
		Assert(.getAttachments(mock) is: #())
		mock.Verify.Times(5).AttachmentsManager_deleteFile([anyArgs:])
		mock.Verify.AttachmentsManager_deleteFile('123', .rec1, 'attach')
		mock.Verify.AttachmentsManager_deleteFile('abc', .rec1, 'attach')
		mock.Verify.AttachmentsManager_deleteFile('222', .rec1, 'attach')
		mock.Verify.AttachmentsManager_deleteFile('456', .rec2, 'custom1')
		mock.Verify.AttachmentsManager_deleteFile('789', .rec3, 'custom2')

		mock.Verify.Times(5).AttachmentsManager_logAction([anyArgs:])
		mock.Verify.AttachmentsManager_logAction('123', [fakeNum: '111'],
			'attach', 'rename')
		mock.Verify.AttachmentsManager_logAction('abc', [fakeNum: '111'],
			'attach', 'rename')
		mock.Verify.AttachmentsManager_logAction('222', [fakeNum: '111'],
			'attach', 'rename')
		mock.Verify.AttachmentsManager_logAction('456', [fakeNum: 'key'],
			'custom1', 'replace')
		mock.Verify.AttachmentsManager_logAction('789', [], 'custom2', 'replace')

		mock = .setupMock()
		mock.ProcessQueue(restore?:)
		Assert(.getAttachments(mock) is: #())
		mock.Verify.Times(5).AttachmentsManager_deleteFile([anyArgs:])
		mock.Verify.AttachmentsManager_deleteFile('333', .rec1, 'attach')
		mock.Verify.AttachmentsManager_deleteFile('222', .rec1, 'attach')
		mock.Verify.AttachmentsManager_deleteFile('abc', .rec1, 'attach')
		mock.Verify.AttachmentsManager_deleteFile('def', .rec2, 'custom1')
		mock.Verify.AttachmentsManager_deleteFile('ghi', .rec3, 'custom2')

		mock.Verify.Never().AttachmentsManager_logAction([anyArgs:])
		}

	Test_RestoreOneByKey()
		{
		// restoring 1 file
		mock = .setupMock()
		mock.RestoreOneByKey(.rec2)
		Assert(.getAttachments(mock) is: Object(
			Object(new_file: 'abc', old_file: '123', rec: .rec1,
				fieldName: 'attach', action: 'rename'),
			Object(new_file: '222', old_file: 'abc', rec: .rec1,
				fieldName: 'attach', action: 'rename'),
			Object(new_file: '333', old_file: '222', rec: .rec1,
				fieldName: 'attach', action: 'rename'),
			Object(new_file: 'ghi', old_file: '789', rec: .rec3,
				fieldName: 'custom2', action: 'replace')))
		mock.Verify.Times(1).AttachmentsManager_deleteFile([anyArgs:])
		mock.Verify.AttachmentsManager_deleteFile('def', .rec2, 'custom1')
		mock.Verify.Never().AttachmentsManager_logAction([anyArgs:])

		// restoring 1 file afer several changes
		mock = .setupMock()
		mock.RestoreOneByKey(.rec1)
		Assert(.getAttachments(mock) is: Object(
			Object(new_file: 'def', old_file: '456', rec: .rec2,
				fieldName: 'custom1', action: 'replace'),
			Object(new_file: 'ghi', old_file: '789', rec: .rec3,
				fieldName: 'custom2', action: 'replace')))
		mock.Verify.Times(3).AttachmentsManager_deleteFile([anyArgs:])
		mock.Verify.AttachmentsManager_deleteFile('333', .rec1, 'attach')
		mock.Verify.AttachmentsManager_deleteFile('222', .rec1, 'attach')
		mock.Verify.AttachmentsManager_deleteFile('abc', .rec1, 'attach')
		mock.Verify.Never().AttachmentsManager_logAction([anyArgs:])

		// restoring using composite key
		mock = .setupMock(skipFill?:)
		mock.AttachmentsManager_keyFields = #(fakeName, fakeDate)
		mock.QueueDeleteFile('abc', '123',
			rec4 = [fakeName: 'test', fakeDate: '202404', attach: ''], 'attach', 'rename')
		mock.QueueDeleteFile('456', 'abc', rec4, 'attach', 'rename')
		mock.QueueDeleteFile('aaa', 'bbb',
			rec5 = [fakeName: 'test', fakeDate: '202405', attach: ''], 'attach', 'rename')
		Assert(.getAttachments(mock) is: Object(
			Object(new_file: 'abc', old_file: '123',  action: 'rename',
				rec: rec4, fieldName: 'attach'),
			Object(new_file: '456', old_file: 'abc', action: 'rename',
				rec: rec4, fieldName: 'attach'),
			Object(new_file: 'aaa', old_file: 'bbb', action: 'rename',
				rec: rec5, fieldName: 'attach')))
		mock.RestoreOneByKey(rec4)
		Assert(.getAttachments(mock) is: Object(
			Object(new_file: 'aaa', old_file: 'bbb', action: 'rename',
				rec: rec5, fieldName: 'attach')))
		mock.Verify.Times(2).AttachmentsManager_deleteFile([anyArgs:])
		mock.Verify.AttachmentsManager_deleteFile('abc', rec4, 'attach')
		mock.Verify.AttachmentsManager_deleteFile('456', rec4, 'attach')
		mock.Verify.Never().AttachmentsManager_logAction([anyArgs:])
		}

	Test_handleStdAttachments_windows()
		{
		if not Sys.Windows?()
			return

		mock = .setupMock(skipFill?:)
		.SpyOn(OpenImageWithLabelsControl.OpenImageWithLabelsControl_getCopyTo).
			Return(`\\PC\Attachments\`)
		mock.handleStdAttachments([])
		Assert(.getAttachments(mock) is: #())

		rec = [fake_attachments: #([attachment0: `202401\file.txt`]), fakeNum: '123']
		mock.handleStdAttachments(rec)
		Assert(.getAttachments(mock) equalsSet: Object(
			Object(new_file: '', old_file: `\\PC\Attachments/202401/file.txt`, :rec,
				fieldName: 'fake_attachments', action: 'record delete')))
		mock.ProcessQueue()
		Assert(.getAttachments(mock) is: #())

		delimiter = OpenImageWithLabelsControl.LabelDelimiter
		rec = [fakeNum: '456', fake_attachments: Object(
			[attachment1: `202401\foo.txt`, attachment3: `202401\bar.txt`],
			[attachment0: `202402\helloworld.txt`],
			[attachment2: `202402\file.txt` $ delimiter $ 'label1, label2',
				attachment3: delimiter $ 'label3'])]
		mock.handleStdAttachments(rec)
		Assert(.getAttachments(mock) equalsSet: Object(
			Object(new_file: '',
				old_file: `\\PC\Attachments/202401/foo.txt`, :rec,
				fieldName: 'fake_attachments', action: 'record delete'),
			Object(new_file: '',
				old_file: `\\PC\Attachments/202401/bar.txt`, :rec,
				fieldName: 'fake_attachments', action: 'record delete'),
			Object(new_file: '',
				old_file: `\\PC\Attachments/202402/helloworld.txt`, :rec,
				fieldName: 'fake_attachments', action: 'record delete'),
			Object(new_file: '',
				old_file: `\\PC\Attachments/202402/file.txt`, :rec,
				fieldName: 'fake_attachments', action: 'record delete')))
		mock.ProcessQueue()
		Assert(.getAttachments(mock) is: #())

		rec = [fakeNum: '111', fake_attachments_display: #(`\\PC\Attachments\foo.txt`)]
		mock.handleStdAttachments(rec)
		Assert(.getAttachments(mock) is: #())
		}

	Test_handleStdAttachments_linux()
		{
		if Sys.Windows?()
			return

		mock = .setupMock(skipFill?:)
		.SpyOn(OpenImageWithLabelsControl.OpenImageWithLabelsControl_getCopyTo).
			Return(`/PC/Attachments/`)
		mock.handleStdAttachments([])
		Assert(.getAttachments(mock) is: #())

		rec = [fake_attachments: #([attachment0: `202401/file.txt`]), fakeNum: '123']
		mock.handleStdAttachments(rec)
		Assert(.getAttachments(mock) equalsSet: Object(
			Object(new_file: '',
				old_file: `/PC/Attachments/202401/file.txt`, :rec,
				fieldName: 'fake_attachments', action: 'record delete')))
		mock.ProcessQueue()
		Assert(.getAttachments(mock) is: #())

		delimiter = OpenImageWithLabelsControl.LabelDelimiter
		rec = [fakeNum: '456', fake_attachments: Object(
			[attachment1: `202401/foo.txt`, attachment3: `202401/bar.txt`],
			[attachment0: `202402/helloworld.txt`],
			[attachment2: `202402/file.txt` $ delimiter $ 'label1, label2',
				attachment3: delimiter $ 'label3'])]
		mock.handleStdAttachments(rec)
		Assert(.getAttachments(mock) equalsSet: Object(
			Object(new_file: '',
				old_file: `/PC/Attachments/202401/foo.txt`, :rec,
				fieldName: 'fake_attachments', action: 'record delete'),
			Object(new_file: '',
				old_file: `/PC/Attachments/202401/bar.txt`, :rec,
				fieldName: 'fake_attachments', action: 'record delete'),
			Object(new_file: '',
				old_file: `/PC/Attachments/202402/helloworld.txt`, :rec,
				fieldName: 'fake_attachments', action: 'record delete'),
			Object(new_file: '',
				old_file: `/PC/Attachments/202402/file.txt`, :rec,
				fieldName: 'fake_attachments', action: 'record delete')))
		mock.ProcessQueue()
		Assert(.getAttachments(mock) is: #())

		rec = [fakeNum: '111', fake_attachments_display: #(`/PC/Attachments/foo.txt`)]
		mock.handleStdAttachments(rec)
		Assert(.getAttachments(mock) is: #())
		}

	Test_handleCustomAttachments_windows()
		{
		if not Sys.Windows?()
			return

		mock = .setupMock(skipFill?:)

		.SpyOn(OpenImageWithLabelsControl.OpenImageWithLabelsControl_getCopyTo).
			Return(`\\PC\Attachments\`)
		mock.handleCustomAttachments([])
		mock.When.customAttachmentField?([anyArgs:]).Return(false)
		mock.When.customAttachmentField?('custom_999999').Return(true)
		Assert(.getAttachments(mock) is: #())

		rec = [custom_999999: `202401\file.txt`, fakeNum: '123',
			custom_999995: `202402\file2.txt`]
		mock.handleCustomAttachments(rec)
		Assert(.getAttachments(mock) equalsSet: Object(
			Object(new_file: '',
				old_file: `\\PC\Attachments/202401/file.txt`, :rec,
				fieldName: 'custom_999999', action: 'record delete')))
		mock.ProcessQueue()
		Assert(.getAttachments(mock) is: #())

		rec = [custom_999999: `202401\helloWorld.txt` $
			OpenImageWithLabelsControl.LabelDelimiter $ 'label1, label2', fakeNum: '123']
		mock.handleCustomAttachments(rec)
		Assert(.getAttachments(mock) equalsSet: Object(
			Object(new_file: '',
				old_file: `\\PC\Attachments/202401/helloWorld.txt`, :rec,
				fieldName: 'custom_999999', action: 'record delete')))
		mock.ProcessQueue()
		Assert(.getAttachments(mock) is: #())
		}

	Test_handleCustomAttachments_linux()
		{
		if Sys.Windows?()
			return

		mock = .setupMock(skipFill?:)

		.SpyOn(OpenImageWithLabelsControl.OpenImageWithLabelsControl_getCopyTo).
			Return(`/PC/Attachments/`)
		mock.handleCustomAttachments([])
		mock.When.customAttachmentField?([anyArgs:]).Return(false)
		mock.When.customAttachmentField?('custom_999999').Return(true)
		Assert(.getAttachments(mock) is: #())

		rec = [custom_999999: `202401/file.txt`, fakeNum: '123',
			custom_999995: `202402/file2.txt`]
		mock.handleCustomAttachments(rec)
		Assert(.getAttachments(mock) equalsSet: Object(
			Object(new_file: '', old_file: `/PC/Attachments/202401/file.txt`, :rec,
			fieldName: 'custom_999999', action: 'record delete')))
		mock.ProcessQueue()
		Assert(.getAttachments(mock) is: #())

		rec = [custom_999999: `202401/helloWorld.txt` $
			OpenImageWithLabelsControl.LabelDelimiter $ 'label1, label2', fakeNum: '123']
		mock.handleCustomAttachments(rec)
		Assert(.getAttachments(mock) equalsSet: Object(
			Object(new_file: '', old_file: `/PC/Attachments/202401/helloWorld.txt`,
			:rec, fieldName: 'custom_999999', action: 'record delete')))
		mock.ProcessQueue()
		Assert(.getAttachments(mock) is: #())
		}

	setupMock(skipFill? = false)
		{
		mock = Mock(AttachmentsManager)
		mock.AttachmentsManager_query = 'fakeTable'
		mock.AttachmentsManager_keyFields = #(fakeNum)
		mock.AttachmentsManager_oldAttachments = Object()
		mock.When.QueueDeleteFile([anyArgs:]).CallThrough()
		mock.When.ProcessQueue([anyArgs:]).CallThrough()
		mock.When.RestoreOneByKey([anyArgs:]).CallThrough()
		mock.When.handleStdAttachments([anyArgs:]).CallThrough()
		mock.When.handleCustomAttachments([anyArgs:]).CallThrough()
		mock.When.deleteFile([anyArgs:]).Return(true)
		mock.When.logAction([anyArgs:]).Return(true)
		mock.When.normallyLinkCopy?([anyArgs:]).Return(true)
		mock.When.skipToDelete?([anyArgs:]).Return(false)
		if not skipFill?
			.fillMock(mock)
		return mock
		}

	fillMock(mock)
		{
		.rec1 = [fakeNum: '111', attach: '']
		.rec2 = [fakeNum: 'key']
		.rec3 = Record()
		mock.QueueDeleteFile('abc', '123', .rec1, 'attach', 'rename')
		mock.QueueDeleteFile('222', 'abc', .rec1, 'attach', 'rename')
		mock.QueueDeleteFile('333', '222', .rec1, 'attach', 'rename')
		mock.QueueDeleteFile('def', '456', .rec2, 'custom1', 'replace')
		mock.QueueDeleteFile('ghi', '789', .rec3, 'custom2', 'replace')
		Assert(.getAttachments(mock) is: Object(
			Object(new_file: 'abc', old_file: '123', rec: .rec1,
				fieldName: 'attach', action: 'rename'),
			Object(new_file: '222', old_file: 'abc', rec: .rec1,
				fieldName: 'attach', action: 'rename'),
			Object(new_file: '333', old_file: '222', rec: .rec1,
				fieldName: 'attach', action: 'rename'),
			Object(new_file: 'def', old_file: '456', rec: .rec2,
				fieldName: 'custom1', action: 'replace'),
			Object(new_file: 'ghi', old_file: '789', rec: .rec3,
				fieldName: 'custom2', action: 'replace')))
		}

	getAttachments(mock)
		{
		return mock.AttachmentsManager_oldAttachments
		}

	Test_skipToDelete?()
		{
		c = AttachmentsManager
			{
			AttachmentsManager_normallyLinkCopy?(@unused)
				{
				return true
				}
			AttachmentsManager_copyTo() { return '' }
			AttachmentsManager_findCreationNumField() { return false }
			}
		c = new c('query', #('key_field'))
		fieldName = 'test_attachments'
		action = 'action'
		new_file = `202404\new_file`
		old_file = `202404\old_file`
		Assert(c.QueueDeleteFile(new_file, '', rec = [], fieldName, action))
		Assert(c.AttachmentsManager_oldAttachments
			is: Object([:new_file, old_file: '', :rec, :fieldName, :action]))

		c = AttachmentsManager
			{
			AttachmentsManager_normallyLinkCopy?(@unused)
				{
				return true
				}
			AttachmentsManager_fileExist?(@unused)
				{
				return false
				}
			AttachmentsManager_copyTo() { return '' }
			AttachmentsManager_protectFolders() { return Object() }
			AttachmentsManager_findCreationNumField() { return false }
			}
		c = new c('query', #('key_field'))
		Assert(c.QueueDeleteFile(new_file, old_file, rec = [],
			'test_attachments', action))
		Assert(c.AttachmentsManager_oldAttachments
			is: Object([:new_file, old_file: '', :rec, :fieldName, :action]))

		cl = AttachmentsManager
			{
			AttachmentsManager_normallyLinkCopy?(@unused)
				{
				return true
				}
			AttachmentsManager_fileExist?(@unused)
				{
				return true
				}
			AttachmentsManager_windows?()
				{
				return true
				}
			AttachmentsManager_copyTo() { return '' }
			AttachmentsManager_protectFolders() { return Object() }
			AttachmentsManager_findCreationNumField() { return false }
			}
		c = new cl('query', #('key_field'))
		// no attachment in the record
		Assert(c.QueueDeleteFile(new_file, old_file, rec = [], fieldName, action))
		Assert(c.AttachmentsManager_oldAttachments
			is: Object([:new_file, :old_file, :rec, :fieldName, :action]))

		c = new cl('query', #('key_field'))
		// custom attachment should not have this issue
		Assert(c.QueueDeleteFile(new_file, old_file, rec = [test_attachments: 'new_file'],
			fieldName, action))
		Assert(c.AttachmentsManager_oldAttachments
			is: Object([:new_file, :old_file, :rec, :fieldName, :action]))

		c = new cl('query', #('key_field'))
		// no attachment in the record
		old_file = 'c:/sub_folder/202404/old_file'
		Assert(c.QueueDeleteFile(new_file, old_file,
			rec = [test_attachments: [[attachment0: 'sub_folder/new_file']]],
			'test_attachments', action))
		Assert(c.AttachmentsManager_oldAttachments
			is: Object([:new_file, :old_file, :rec, :fieldName, :action]))

		c = new cl('query', #('key_field'))
		old_file = 'c:/sub_folder/old_file'
		Assert(c.QueueDeleteFile(new_file, old_file,
			rec = [test_attachments: [[attachment0: 'sub_folder/new_file']]],
			'test_attachments', action))
		Assert(c.AttachmentsManager_oldAttachments
			is: Object([:new_file, old_file: '', :rec, :fieldName, :action]))

		c = new cl('query', #('key_field'))
		old_file = OpenImageWithLabelsControl.SplitFullPath('202404/old_file')
		Assert(c.QueueDeleteFile(new_file, old_file,
			rec = [test_attachments: [[attachment0: '202404/old_file']]],
			'test_attachments', action))
		Assert(c.AttachmentsManager_oldAttachments
			is: Object([:new_file, old_file: '', :rec, :fieldName, :action]))

		c = new cl('query', #('key_field'))
		old_file = 'c:/sub_folder/202404/old_file'
		Assert(c.QueueDeleteFile('', old_file,
			rec = [test_attachments: [[attachment0: '202404/old_file']]],
			'test_attachments', AttachmentsManager.RecordDeleteAction))
		Assert(c.AttachmentsManager_oldAttachments
			is: Object([new_file: '', :old_file, :rec, :fieldName,
				action: AttachmentsManager.RecordDeleteAction]))

		// case insensitive on windows
		c = new cl('query', #('key_field'))
		old_file = OpenImageWithLabelsControl.SplitFullPath(`202404\OLD_FILE`)
		Assert(c.QueueDeleteFile(new_file, old_file,
			rec = [test_attachments: [[attachment0: '202404/old_file']]],
			'test_attachments', action))
		Assert(c.AttachmentsManager_oldAttachments
			is: Object([:new_file, old_file: '', :rec, :fieldName, :action]))

		cl = AttachmentsManager
			{
			AttachmentsManager_normallyLinkCopy?(@unused)
				{
				return true
				}
			AttachmentsManager_fileExist?(@unused)
				{
				return true
				}
			AttachmentsManager_windows?()
				{
				return false
				}
			AttachmentsManager_copyTo() { return '' }
			AttachmentsManager_protectFolders() { return Object() }
			AttachmentsManager_findCreationNumField() { return false }
			}
		c = new cl('query', #('key_field'))
		old_file = OpenImageWithLabelsControl.SplitFullPath('202404/old_file')
		Assert(c.QueueDeleteFile(new_file, old_file,
			rec = [test_attachments: [[attachment0: '202404/old_file']]],
			'test_attachments', action))
		Assert(c.AttachmentsManager_oldAttachments
			is: Object([:new_file, old_file: '', :rec, :fieldName, :action]))

		c = new cl('query', #('key_field'))
		// case sensitive on non-windows
		old_file = `c:\sub_folder\202404\OLD_FILE`
		Assert(c.QueueDeleteFile(new_file, old_file,
			rec = [test_attachments: [[attachment0: '202404/old_file']]],
			'test_attachments', action))
		Assert(c.AttachmentsManager_oldAttachments
			is: Object([:new_file, :old_file, :rec, :fieldName, :action]))
		}

	Test_oldFileNotLinkCopy()
		{
		new_file = ''
		action = 'remove'
		fieldName = 'test_attachments'
		cl = AttachmentsManager
			{
			AttachmentsManager_normallyLinkCopy?() { return true }
			AttachmentsManager_copyTo() { return `\\share\att` }
			AttachmentsManager_protectFolders() { Object('OldAttachments') }
			AttachmentsManager_fileExist?(unused) { return true }
			AttachmentsManager_findCreationNumField() { return false }
			}

		c = new cl('query', #('key_field'))
		// file prefix indicates file was not attached with "normally copy & link"
		old_file = 'myPdf.pdf'
		Assert(c.QueueDeleteFile(new_file, old_file, [test_attachments: [
			[attachment0: 'myPdf.pdf']
			]], 'test_attachments', action))
		Assert(c.AttachmentsManager_oldAttachments is: Object())

		// file prefix indicates file was not attached with "normally copy & link"
		c.AttachmentsManager_oldAttachments = Object()
		old_file = `C:\business\axon\attachments\myPdf.pdf`
		Assert(c.QueueDeleteFile(new_file, old_file, [test_attachments: [
			[attachment0: `C:\business\axon\attachments\myPdf.pdf`]
			]], 'test_attachments', action))
		Assert(c.AttachmentsManager_oldAttachments is: Object())

		// file is in protected "old attachments" folder - carried over from
		// pre normally copy and link activity
		c.AttachmentsManager_oldAttachments = Object()
		old_file = `\\share\att\OldAttachments\myPdf.pdf`
		Assert(c.QueueDeleteFile(new_file, old_file, [test_attachments: [
			[attachment0: 'OldAttachments\myPdf.pdf']
			]], 'test_attachments', action))
		Assert(c.AttachmentsManager_oldAttachments is: Object())

		// file is in protected "old attachments" folder - carried over from
		// pre normally copy and link activity - file is replaced - new in std location
		c.AttachmentsManager_oldAttachments = Object()
		old_file = `\\share\att\OldAttachments\myPdf.pdf`
		new_file = `\\share\att\202403\myPdf2.pdf`
		Assert(c.QueueDeleteFile(new_file, old_file,
			rec = [test_attachments: [[attachment0: '202403\myPdf2.pdf']]],
			'test_attachments', 'replace'))
		Assert(c.AttachmentsManager_oldAttachments is:
			Object([:new_file, old_file: '', :rec, :fieldName, action: 'replace']))

		// file is in protected "old attachments" folder - carried over from
		// pre normally copy and link activity - file renamed, can still new file for
		// cleanup
		c.AttachmentsManager_oldAttachments = Object()
		old_file = `\\share\att\OldAttachments\myPdf.pdf`
		new_file = `\\share\att\OldAttachments\renamed.pdf`
		Assert(c.QueueDeleteFile(new_file, old_file,
			rec = [test_attachments: [[attachment0: 'OldAttachments\renamed.pdf']]],
			'test_attachments', 'rename'))
		Assert(c.AttachmentsManager_oldAttachments is:
			Object([:new_file, old_file: '', :rec, :fieldName, action: 'rename']))

		// file is managed by normally link copy settings, in standard location
		c.AttachmentsManager_oldAttachments = Object()
		old_file = `\\share\att\202403\myPdf.pdf`
		Assert(c.QueueDeleteFile(new_file, old_file, rec = [test_attachments: []],
			'test_attachments', 'replace'))
		Assert(c.AttachmentsManager_oldAttachments is:
			Object([:new_file, :old_file, :rec, :fieldName, action: 'replace']))

		// file is NOT managed by normally link copy settings,
		// reflects sample companies demo data
		c.AttachmentsManager_oldAttachments = Object()
		old_file = `\\share\att\myPdf.pdf`
		Assert(c.QueueDeleteFile(new_file, old_file, rec = [test_attachments: []],
			'test_attachments', 'replace'))
		Assert(c.AttachmentsManager_oldAttachments is:
			Object([:new_file, old_file: '', :rec, :fieldName, action: 'replace']))

		// file is managed by normally link copy settings, in standard location
		c.AttachmentsManager_oldAttachments = Object()
		old_file = `\\share\att\202403\myPdf.pdf`
		Assert(c.QueueDeleteFile(new_file, old_file, rec = [],
			'test_attachments', 'replace'))
		Assert(c.AttachmentsManager_oldAttachments is:
			Object([:new_file, :old_file, :rec, :fieldName, action: 'replace']))
		}

	Test_findCreationNumField()
		{
		table = .MakeTable('(a, b, c) key (a,b)')
		mgr = AttachmentsManager(table, #(a,b))
		Assert(mgr.AttachmentsManager_creationNumField is: false)

		table = .MakeTable('(test_num, test_name, test_abbrev) key (test_num)')
		mgr = AttachmentsManager(table, #(test_num))
		Assert(mgr.AttachmentsManager_creationNumField is: 'test_num')

		mgr = AttachmentsManager(table, #(test_num, test_name))
		Assert(mgr.AttachmentsManager_creationNumField is: 'test_num')

		mgr = AttachmentsManager(table, #(test_name))
		Assert(mgr.AttachmentsManager_creationNumField is: 'test_num')
		}

	Test_checkCreationDate()
		{
		cl = AttachmentsManager
			{
			AttachmentsManager_query: 'test_query'
			AttachmentsManager_keyFields: (test_num, test_name)
			AttachmentsManager_creationNumField: 'test_num'
			}
		Assert(cl.AttachmentsManager_checkCreationDate([test_num: #20240416]))
		Assert(cl.AttachmentsManager_checkCreationDate([test_num: #20240417]) is: false)
		Assert(cl.AttachmentsManager_checkCreationDate([test_num: false]) is: false)

		cl = AttachmentsManager
			{
			AttachmentsManager_query: 'test_query'
			AttachmentsManager_keyFields: (test_num)
			}
		Assert(cl.AttachmentsManager_checkCreationDate([test_num: #20240416]) is: false)
		Assert(cl.AttachmentsManager_checkCreationDate([test_num: #20240417]) is: false)
		Assert(cl.AttachmentsManager_checkCreationDate([test_num: false]) is: false)
		}
	}
