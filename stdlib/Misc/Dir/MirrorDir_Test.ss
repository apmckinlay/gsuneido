// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		_keepExtra? = false
		_failedList = Object()
		_skipDirs = #('a/ab', 'b/af/')
		mock = Mock(MirrorDir)
		mock.When.ensureSlash('a').Return('a/')
		mock.When.ensureSlash('b').Return('b/')
		mock.When.slash().Return('/')
		mock.When.ensureDir('b').Return(true)
		mock.When.ensureDir('b/aa').Return(true)
		mock.When.copyFile('a/ab', 'b/ab').Return(true)
		mock.When.copyFile('a/ad', 'b/ad').Return(false)
		mock.When.copyFile('a/ag', 'b/ag').Return(true)
		mock.When.different?('a/ab', 'b/ab').Return(true)
		mock.When.different?('a/ac', 'b/ac').Return(false)
		mock.When.different?('a/ad', 'b/ad').Return(true)
		mock.When.dirList('a/').Return(
			#('aa/': 0, 'ab': 1, 'ac': 2, 'ad': 3, 'ag': 111))
		mock.When.dirList('b/').Return(
			#('aa/': 0, 'af/': 0, 'ab': 1, 'ac': 2, 'ad': 3, 'ae': 4, 'ag': 222))
		mock.When.skip?([anyArgs:]).CallThrough()
		mock.When.deleteExtra([anyArgs:]).CallThrough()
		mock.Eval(MirrorDir.MirrorDir_mirror, 'a', 'b')
		Assert(_failedList is: #('a/ad'))
		mock.Verify.deleteItem('b/', 'ae')
		mock.Verify.Never().deleteItem('b/', 'af/') // skipped
		mock.Verify.mirror('a/aa/', 'b/aa/')
		mock.Verify.Never().different?('a/ag', 'b/ag')
		mock.Verify.Times(1).copyFile('a/ag', 'b/ag')
		mock.Verify.Never().copyFile('a/ab', 'b/ab') // skipped
		mock.Verify.Times(1).copyFile('a/ad', 'b/ad')
		mock.Verify.Never().copyFile('a/ac', 'b/ac') // Same
		}

	Test_keepExtra()
		{
		_keepExtra? = true
		_failedList = Object()
		mock = Mock(MirrorDir)
		mock.When.ensureSlash('a').Return('a/')
		mock.When.ensureSlash('b').Return('b/')
		mock.When.slash().Return('/')
		mock.When.ensureDir('b').Return(true)
		mock.When.ensureDir('b/aa').Return(true)
		mock.When.copyFile('a/ab', 'b/ab').Return(true)
		mock.When.copyFile('a/ad', 'b/ad').Return(false)
		mock.When.different?('a/ab', 'b/ab').Return(true)
		mock.When.different?('a/ac', 'b/ac').Return(false)
		mock.When.different?('a/ad', 'b/ad').Return(true)
		mock.When.dirList('a/').Return(
			#('aa/': 0, 'ab': 0, 'ac' : 0, 'ad': 0))
		mock.When.dirList('b/').Return(
			#('aa/': 0, 'af/': 0, ab: 0, ac: 0, ad: 0, ae: 0))
		mock.When.skip?([anyArgs:]).Return(false)
		mock.When.deleteExtra([anyArgs:]).CallThrough()
		mock.Eval(MirrorDir.MirrorDir_mirror, 'a', 'b')
		Assert(_failedList is: #('a/ad'))
		mock.Verify.mirror('a/aa/', 'b/aa/')
		}

	Test_DeleteFilesAndDirs()
		{
		_failedList = Object()
		mock = Mock()
		mock.When.MirrorDir_ensureSlash('a').Return('a/')
		mock.When.MirrorDir_getFile('a/list').Return(
			'a_folder/\na_file\nb_folder/b_file\n  \n')
		Assert(mock.Eval(MirrorDir.DeleteFilesAndDirs, 'a', 'list') is: #())
		mock.Verify.MirrorDir_deleteItem('a/', 'a_folder/')
		mock.Verify.MirrorDir_deleteItem('a/', 'a_file')
		mock.Verify.MirrorDir_deleteItem('a/', 'b_folder/b_file')
		mock.Verify.MirrorDir_deleteItem('a/', 'list')

		_failedList = Object()
		mock = Mock()
		mock.When.MirrorDir_ensureSlash('a').Return('a/')
		mock.When.MirrorDir_getFile('a/list').Return(false)
		Assert(mock.Eval(MirrorDir.DeleteFilesAndDirs, 'a', 'list') is: #())
		mock.Verify.Never().MirrorDir_deleteItem([anyArgs:])

		_failedList = Object()
		mock = Mock()
		mock.When.MirrorDir_ensureSlash('a').Return('a/')
		mock.When.MirrorDir_getFile('a/list').Return('')
		Assert(mock.Eval(MirrorDir.DeleteFilesAndDirs, 'a', 'list') is: #())
		mock.Verify.MirrorDir_deleteItem('a/', 'list')
		}

	Test_deleteItem()
		{
		mock = Mock()
		mock.When.MirrorDir_deleteDir('a/a_folder/').Return(true)
		_failedList = Object()
		mock.Eval(MirrorDir.MirrorDir_deleteItem, 'a/', 'a_folder/')
		Assert(_failedList is: #())

		mock = Mock()
		mock.When.MirrorDir_deleteDir('a/b_folder/').Return(false)
		_failedList = Object()
		mock.Eval(MirrorDir.MirrorDir_deleteItem, 'a/', 'b_folder/')
		Assert(_failedList is: #('a/b_folder/'))

		mock = Mock()
		mock.When.MirrorDir_deleteFile('a/a_file').Return(true)
		_failedList = Object()
		mock.Eval(MirrorDir.MirrorDir_deleteItem, 'a/', 'a_file')
		Assert(_failedList is: #())

		mock = Mock()
		mock.When.MirrorDir_deleteFile('a/b_file').Return(false)
		_failedList = Object()
		mock.Eval(MirrorDir.MirrorDir_deleteItem, 'a/', 'b_file')
		Assert(_failedList is: #('a/b_file'))
		}
	}