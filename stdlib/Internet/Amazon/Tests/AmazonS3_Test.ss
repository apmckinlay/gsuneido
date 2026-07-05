// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ListOperations()
		{
		mock = .initializeMock()
		mock.When.listTruncated([anyArgs:]).CallThrough()
		mock.When.https([anyArgs:]).Return(AmazonTestReponses.emptyListXML)
		Assert(mock.Eval(AmazonS3.ListBucketContents, 'bucket') is: #())
		mock.Verify.https('GET', 'https://s3.us-east-1.amazonaws.com/bucket/', '',
			'', 'signed request', '', '')

		mock.When.https([anyArgs:]).Return(AmazonTestReponses.hasFileXML)
		mock.When.signRequest([anyArgs:]).Return('signed request2')
		Assert(mock.Eval(AmazonS3.ListBucketContents, 'bucket')
			is: Object(
				Object(size: '171 bytes', key: "test.txt", owner: "axoneta",
					last_modified: #20111109.171628.GMTimeToLocal(),
					storage: 'STANDARD')
				))
		mock.Verify.https('GET', 'https://s3.us-east-1.amazonaws.com/bucket/', '',
			'', 'signed request2', '', '')

		Assert(mock.Eval(AmazonS3.ListBucketContents, 'bucket', rawResponse:)
			is: Object(
				Object(size: '171', key: "test.txt", owner: "axoneta",
					last_modified: '2011-11-09T17:16:28.000Z',
					storage: 'STANDARD')
				))

		mock.When.https([anyArgs:]).Return(AmazonTestReponses.invalidXML)
		mock.Eval(AmazonS3.ListBucketContents, 'bucket')
		mock.Verify.addToLog('LIST', 'listTruncated - unmatched tag: etag',
			AmazonTestReponses.invalidXML.content, params: [anyObject:])

		mock.When.ListBucketContents([anyArgs:]).CallThrough()
		mock.When.https([anyArgs:]).Return(AmazonTestReponses.hasFileXML)
		Assert(mock.Eval(AmazonS3.ListBucketFiles, 'bucket')
			is: #("test.txt"))

		// this tests prefix option on ListBucketContents as well.
		mock.When.ListBucketFiles([anyArgs:]).CallThrough()
		mock.When.https([anyArgs:]).Return(AmazonTestReponses.hasFolderFileXML)
		Assert(mock.Eval(AmazonS3.ListBucketFolderFiles, 'bucketWithFolder',
			'testFolder') is: #('testFolder/folderTest.txt'))
		mock.Verify.https('GET',
			"https://s3.us-east-1.amazonaws.com/bucketWithFolder/?prefix=testFolder", '',
			'', 'signed request2', '', '')

		// truncated list
		mock.When.https([anyArgs:]).
			Return(AmazonTestReponses.hasFileXMLTruncated, AmazonTestReponses.hasFileXML)
		mock.When.signRequest([anyArgs:]).Return('signed request3')
		Assert(mock.Eval(AmazonS3.ListBucketFiles, 'bucket', '') is:
			#('test0.txt', 'test.txt'))
		mock.Verify.https('GET',
			"https://s3.us-east-1.amazonaws.com/bucket/", '', '', 'signed request3', '',
			'')
		mock.Verify.https('GET',
			"https://s3.us-east-1.amazonaws.com/bucket/?marker=test0.txt", '', '',
			'signed request3', '', '')

		// bad response code
		xmlReturn = AmazonTestReponses.hasFileXML.Copy()
		xmlReturn.header = 'HTTP/1.1 400'
		mock.When.https([anyArgs:]).Return(xmlReturn)
		Assert(mock.Eval(AmazonS3.ListBucketContents, 'bucket') is: #())
		mock.Verify.addToLog('GET', 'Bad HTTP Status Code (400)',
			'HTTP/1.1 400\r\n\r\n' $ AmazonTestReponses.hasFileXML.content, [anyObject:])

		// bad/empty header
		xmlReturn = #(header: 'bad header')
		mock.When.https([anyArgs:]).Return(xmlReturn)
		Assert(mock.Eval(AmazonS3.ListBucketContents, 'bucket') is: #())
		mock.Verify.addToLog('GET',
			'Bad HTTP Status Code (Invalid HTTP response code in: bad header)',
			'bad header\r\n\r\n', [anyObject:])

		// list versioned files
		mock.When.ListBucketFileVersions([anyArgs:]).CallThrough()
		mock.When.https([anyArgs:]).Return(AmazonTestReponses.hasVersionsOfFileXML)
		Assert(mock.Eval(AmazonS3.ListBucketFileVersions, 'bucket')
			is: Object(Object(versionId: "iMccVr_1xYzUJI40moR0dsGHBlgYDHkU",
				last_modified: #20210111.193341.GMTimeToLocal(), owner: "axoneta",
				size: "48505180 bytes",	key: "axon.betacad.gpg",
				storage: 'STANDARD'),
				Object(versionId: "70YsAq4iFPYCwHI0pScpa4D_TWG07fod",
					last_modified: #20210108.205825.GMTimeToLocal(), owner: "axoneta",
					size: "48522359 bytes", key: "axon.betacad.gpg",
					storage: 'STANDARD')))

		// invalid credentials
		mock = .initializeMock()
		mock.When.https([anyArgs:]).Return('') // ensure we don't make a real request
		mock.When.listTruncated([anyArgs:]).CallThrough()
		mock.When.signRequest([anyArgs:]).Return(false)
		Assert(mock.Eval(AmazonS3.ListBucketContents, 'bucket') is: #())
		mock.Verify.Never().https([anyArgs:])

		// throttle response
		mock = .initializeMock()
		throttleResponse = AmazonTestReponses.hasFileXML.Copy()
		throttleResponse.header = 'HTTP/1.1 503'
		mock.When.https([anyArgs:]).
			Return(throttleResponse, AmazonTestReponses.hasFileXML)
		mock.When.ListBucketContents([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonS3.ListBucketFiles, 'bucket')
			is: #("test.txt"))
		}

	initializeMock()
		{
		mock = Mock(AmazonS3)
		mock.When.addToLog([anyArgs:]).Return('') // ensure we never add to the log file
		mock.When.SecurityToken().Return('token')
		mock.When.signRequest([anyArgs:]).Return('signed request')
		return mock
		}

	Test_buildPutFileRequest()
		{
		.putFile_testNoPolicy()
		.putFile_testWithPolicy()
		.putFile_testBadFile()
		.putFile_testBadCredentials()
		.putFile_requestFailed()
		.putFile_throttle()
		.putFile_addTags()
		}

	putFile_testNoPolicy()
		{
		mock = .initializeMock()
		mock.When.validFile?('test1_got', false).Return(true)
		mock.When.makeRequest([anyArgs:]).CallThrough()
		mock.When.validateAndGetFileTo([anyArgs:]).CallThrough()
		mock.When.https([anyArgs:]).Return(AmazonTestReponses.httpsReturn)
		Assert(mock.Eval(AmazonS3.PutFile, 'bucket', 'test1_got', policy: ''))
		mock.Verify.signRequest('PUT', 'us-east-1', '', '/bucket/test1_got',
			#(X_Amz_Content_Sha256: 'UNSIGNED-PAYLOAD',	X_Amz_Security_Token: 'token'))
		mock.Verify.throttle(
			'PUT', expectedResponse: '200', fullResponse?: false, block: [any:])
		mock.Verify.https('PUT', 'https://s3.us-east-1.amazonaws.com/bucket/test1_got',
			'', 'test1_got', 'signed request', '', '')
		}

	putFile_testWithPolicy()
		{
		mock = .initializeMock()
		mock.When.validFile?('test1_got', false).Return(true)
		mock.When.makeRequest([anyArgs:]).CallThrough()
		mock.When.validateAndGetFileTo([anyArgs:]).CallThrough()
		mock.When.https([anyArgs:]).Return(AmazonTestReponses.httpsReturn)
		Assert(mock.Eval(AmazonS3.PutFile, 'bucket', 'test1_got',
			policy: 'x-amz-acl:public-read'))
		mock.Verify.signRequest('PUT', 'us-east-1', '', '/bucket/test1_got',
			#(X_Amz_Content_Sha256: 'UNSIGNED-PAYLOAD',	X_Amz_Security_Token: 'token',
				X_Amz_Acl: 'public-read'))
		mock.Verify.throttle(
			'PUT', expectedResponse: '200', fullResponse?: false, block: [any:])
		mock.Verify.https('PUT', 'https://s3.us-east-1.amazonaws.com/bucket/test1_got',
			'', 'test1_got', 'signed request', '', '')
		}

	putFile_testBadFile()
		{
		mock = .initializeMock()
		mock.When.validateAndGetFileTo([anyArgs:]).CallThrough()
		mock.When.validFile?('test3', false).Return(false)
		Assert(mock.Eval(AmazonS3.PutFile,
			'bucket', 'test3', 'test3_got', 'x-amz-acl:public-read')
			is: false)
		mock.Verify.Never().makeRequest([anyArgs:])
		}

	putFile_testBadCredentials()
		{
		mock = .initializeMock()
		mock.When.validFile?('test4', false).Return(true)
		mock.When.makeRequest([anyArgs:]).CallThrough()
		mock.When.validateAndGetFileTo([anyArgs:]).CallThrough()
		mock.When.signRequest([anyArgs:]).Return(false)
		Assert(mock.Eval(AmazonS3.PutFile,
			'bucket', 'test4', 'test4_got', 'x-amz-acl:public-read')
			is: false)
		mock.Verify.Never().https([anyArgs:])
		}

	putFile_requestFailed()
		{
		mock = .initializeMock()
		mock.When.validFile?('test4', false).Return(true)
		mock.When.makeRequest([anyArgs:]).CallThrough()
		mock.When.validateAndGetFileTo([anyArgs:]).CallThrough()
		mock.When.https([anyArgs:]).Return(AmazonTestReponses.httpsFailure)
		Assert(mock.Eval(AmazonS3.PutFile,
			'bucket', 'test4', 'test4_got', 'x-amz-acl:public-read')
			is: false)
		}

	putFile_throttle()
		{
		mock = .initializeMock()
		mock.When.validFile?('test1_got', false).Return(true)
		mock.When.makeRequest([anyArgs:]).CallThrough()
		mock.When.validateAndGetFileTo([anyArgs:]).CallThrough()
		httpsReturn1 = AmazonTestReponses.httpsReturn.Copy()
		httpsReturn1.header = AmazonTestReponses.httpsReturn.header.
			Replace('200 OK','503')

		mock.When.https([anyArgs:]).Return(httpsReturn1, AmazonTestReponses.httpsReturn)
		Assert(mock.Eval(AmazonS3.PutFile, 'bucket', 'test1_got', policy: ''))
		mock.Verify.signRequest('PUT', 'us-east-1', '', '/bucket/test1_got',
			#(X_Amz_Content_Sha256: 'UNSIGNED-PAYLOAD',	X_Amz_Security_Token: 'token'))
		mock.Verify.throttle(
			'PUT', expectedResponse: '200', fullResponse?: false, block: [any:])
		mock.Verify.Times(2).https('PUT',
			'https://s3.us-east-1.amazonaws.com/bucket/test1_got', '',
			'test1_got', 'signed request', '', '')
		}

	putFile_addTags()
		{
		mock = .initializeMock()
		mock.When.makeRequest([anyArgs:]).CallThrough()
		mock.When.signRequest([anyArgs:]).Return('signedRequest')
		mock.When.https([anyArgs:]).Return(#(header: "HTTP/1.1 200 OK", content: ''))
		tags = Object(testingTag: true)
		mock.Eval(AmazonS3.PutFileTagging, 'bucket', 'randomFile', tags)
		expectedXml = '<?xml version="1.0" encoding="UTF-8"?>
		<Tagging xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><TagSet><Tag>
				<Key>testingTag</Key>
				<Value>true</Value>
			</Tag></TagSet></Tagging>'
		mock.Verify.https('PUT',
			'https://s3.us-east-1.amazonaws.com/bucket/randomFile?tagging=1',
			'', '', 'signedRequest', expectedXml, '')
		}

	httpsNoSuchBucket: #(content: `<?xml version="1.0" encoding="UTF-8"?>
<Error><Code>NoSuchBucket</Code><Message>The specified bucket does not exist</Message>` $
`<BucketName>notARealBucket</BucketName><RequestId>96954F08AA184283</RequestId>` $
`<HostId>gbmbdjZWtOYFTAqFptERYgDNfm/sAG/IXAakgeTGMgDOpa+fXiUkkpBBqtoyGPcsd60GB1I0zhE=` $
`</HostId></Error>`, header: "HTTP/1.1 404 Not Found
x-amz-request-id: 96954F08AA184283
x-amz-id-2: gbmbdjZWtOYFTAqFptERYgDNfm/sAG/IXAakgeTGMgDOpa+fXiUkkpBBqtoyGPcsd60GB1I0zhE=
Content-Type: application/xml
Transfer-Encoding: chunked
Date: Tue, 21 Aug 2018 15:02:31 GMT
Server: AmazonS3")

	Test_GetFile()
		{
		mock = .initializeMock()
		mock.When.https([anyArgs:]).Return(AmazonTestReponses.httpsReturn)
		mock.When.makeRequest([anyArgs:]).CallThrough()
		Assert(mock.Eval(AmazonS3.GetFile, 'bucket', 'test1.txt'))
		mock.Verify.https('GET', 'https://s3.us-east-1.amazonaws.com/bucket/test1.txt',
			'test1.txt', '', 'signed request', '', '')

		mock.When.https([anyArgs:]).Return(.httpsNoSuchBucket)
		Assert(mock.Eval(AmazonS3.GetFile, 'badBucket', 'test2.txt') is: false)
		mock.Verify.https('GET', 'https://s3.us-east-1.amazonaws.com/badBucket/test2.txt',
			'test2.txt', '', 'signed request', '', '')

		mock.When.https([anyArgs:]).Return(.httpsReturnTagging)
		Assert(mock.Eval(AmazonS3.GetFileTagging, 'bucket', 'File')
			is: #(TestTagging: true))

		mock.When.https([anyArgs:]).Return(AmazonTestReponses.httpsReturn)
		mock.When.makeRequest([anyArgs:]).Return(false)
		Assert(mock.Eval(AmazonS3.GetFile, 'bucket', 'BadFile.txt') is: false)
		mock.Verify.cleanupFileIfFailed(false, 'BadFile.txt')
		}

	httpsReturnTagging: #(content: '<?xml version="1.0" encoding="UTF-8"?>
<Tagging xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><TagSet><Tag>' $
'<Key>TestTagging</Key><Value>true</Value></Tag></TagSet></Tagging>',
	header: "HTTP/1.1 200 OK")

	httpBucketNotEmpty: #(content: '<?xml version="1.0" encoding="UTF-8"?>
<Error><Code>BucketNotEmpty</Code><Message>The bucket you tried to delete is not empty' $
'</Message><BucketName>TestBucketAxon</BucketName><RequestId>45FAD6EA6D8BC614' $
'</RequestId><HostId>xu4WXgQKJdpdd359GSlY95Hp75MDtQFhch+MZkiuG8JxyZSalqZZGsa9/' $
'hAMNpmFt00BLDE1CYc=</HostId></Error>', header: "HTTP/1.1 409 Conflict
x-amz-request-id: 45FAD6EA6D8BC614
x-amz-id-2: xu4WXgQKJdpdd359GSlY95Hp75MDtQFhch+MZkiuG8JxyZSalqZZGsa9/hAMNpmFt00BLDE1CYc=
Content-Type: application/xml
Transfer-Encoding: chunked
Date: Tue, 21 Aug 2018 15:15:39 GMT
Server: AmazonS3")

	httpDeleteSuccess: #(content: "", header: "HTTP/1.1 204 No Content
x-amz-id-2: iYimqXmECnV58xUHtr4kslbVWLrfB50ae9Qwb/NueMINQjQGFL1kC63kB1e/k1758t8CxE0+aEE=
x-amz-request-id: BFD154BA3EE59FB1
Date: Tue, 21 Aug 2018 15:18:01 GMT
Server: AmazonS3")

	Test_DeleteFile()
		{
		mock = .initializeMock()
		mock.When.makeRequest([anyArgs:]).CallThrough()
		mock.When.https([anyArgs:]).Return(AmazonTestReponses.httpsFailure)
		mock.When.GetBucketLocationCached([anyArgs:]).Return('fakeregion')
		Assert(mock.Eval(AmazonS3.DeleteFile, 'NotARealBucket', 'test1.txt') is: false)
		mock.Verify.https('DELETE',
			'https://s3.fakeregion.amazonaws.com/NotARealBucket/test1.txt', '', '',
			'signed request', '', '')

		mock.When.https([anyArgs:]).Return(.httpBucketNotEmpty)
		Assert({ mock.Eval(AmazonS3.DeleteFile, 'NotAnEmptyBucket', '') }
			throws: 'expected the value to not be "" but it was')
		mock.Verify.Never().https('DELETE',
			'https://s3.fakeregion.amazonaws.com/NotAnEmptyBucket/', '', '',
			'signed request', '', '')

		mock.When.https([anyArgs:]).Return(.httpDeleteSuccess)
		Assert(mock.Eval(AmazonS3.DeleteFile, 'bucket', 'fileToDelete'))
		mock.Verify.https('DELETE',
			'https://s3.fakeregion.amazonaws.com/bucket/fileToDelete', '', '',
			'signed request', '', '')
		}

	Test_DeleteFiles()
		{
		mock = .initializeMock()
		mock.When.DeleteFile([anyArgs:]).Return(false)
		mock.When.GetBucketLocationCached([anyArgs:]).Return('fakeregion')

		Assert(mock.Eval(AmazonS3.DeleteFiles, 'testBucket',
			Object('file1.txt', 'file2.log'))
				is: 'Deleting files from Amazon S3 failed:\n\n\tfile1.txt, file2.log')

		mock = .initializeMock()
		mock.When.GetBucketLocationCached([anyArgs:]).Return('fakeregion')
		mock.When.DeleteFile([anyArgs:]).CallThrough()
		mock.When.https([anyArgs:]).Return(.httpDeleteSuccess)
		Assert(mock.Eval(AmazonS3.DeleteFiles, 'bucket',
			Object('delete1.txt', 'delete2.txt')) is: '')
		mock.Verify.https('DELETE',
			'https://s3.fakeregion.amazonaws.com/bucket/delete1.txt', '', '',
			'signed request', '', '')
		mock.Verify.https('DELETE',
			'https://s3.fakeregion.amazonaws.com/bucket/delete2.txt', '', '',
			'signed request', '', '')

		mock.When.https([anyArgs:]).
			Return(.httpDeleteSuccess, AmazonTestReponses.httpsFailure)
		Assert(mock.Eval(AmazonS3.DeleteFiles, 'bucket',
			Object('deleteWorked.txt', 'deleteFailed.txt'))
				is: 'Deleting files from Amazon S3 failed:\n\n\tdeleteFailed.txt')
		mock.Verify.https('DELETE',
			'https://s3.fakeregion.amazonaws.com/bucket/deleteWorked.txt', '', '',
			'signed request', '', '')
		mock.Verify.https('DELETE',
			'https://s3.fakeregion.amazonaws.com/bucket/deleteFailed.txt', '', '',
			'signed request', '', '')
		}

	Test_DeleteFiles2()
		{
		_count = Object(0)
		_requestFailed? = Object(false)
		_makeRequestRes = Object(
`<DeleteResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
<Deleted>
<Key>sample1</Key>
<DeleteMarker>true</DeleteMarker>
<DeleteMarkerVersionId>null</DeleteMarkerVersionId>
</Deleted>
<Deleted>
<Key>sample2</Key>
<DeleteMarker>true</DeleteMarker>
<DeleteMarkerVersionId>null</DeleteMarkerVersionId>
</Deleted>
</DeleteResult>`,
`<DeleteResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
<Deleted>
<Key>sample1</Key>
<DeleteMarker>true</DeleteMarker>
<DeleteMarkerVersionId>null</DeleteMarkerVersionId>
</Deleted>
<Error>
<Key>sample2</Key>
<Code>AccessDenied</Code>
<Message>Access Denied</Message>
</Error>
</DeleteResult>`)

		cl = AmazonS3
			{
			AmazonS3_deleteFileLimit: 2
			AmazonS3_makeRequest(@unused)
				{
				_count[0]++
				if _requestFailed?[0] is true
					return false

				res = _makeRequestRes.PopFirst()
				Assert(res isnt _makeRequestRes)
				return res
				}
			GetBucketLocationCached(bucket)
				{
				return bucket
				}
			}

		fn = cl.DeleteFiles2
		Assert(fn('sample_bucket', #('sample1', 'sample2')) is: '')
		Assert(_count[0] is: 1)
		Assert(_makeRequestRes isSize: 1)

		_count[0] = 0
		_requestFailed?[0] = true
		Assert(fn('sample_bucket', #('sample1', 'sample2')) is:
			'Deleting files from Amazon S3 failed:\n\n\tsample1, sample2')
		Assert(_count[0] is: 1)

		_count[0] = 0
		_requestFailed?[0] = false
		result = fn('sample_bucket', #('sample1', 'sample2'))
		Assert(result is: 'Deleting files from Amazon S3 failed:\n\n\tsample2')
		Assert(_count[0] is: 1)
		Assert(_makeRequestRes isSize: 0)

		_count[0] = 0
		_makeRequestRes.Append(Object(
`<DeleteResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
<Deleted>
<Key>sample1</Key>
<DeleteMarker>true</DeleteMarker>
<DeleteMarkerVersionId>null</DeleteMarkerVersionId>
</Deleted>
<Deleted>
<Key>sample2</Key>
<DeleteMarker>true</DeleteMarker>
<DeleteMarkerVersionId>null</DeleteMarkerVersionId>
</Deleted>
</DeleteResult>`,
`<DeleteResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
<Error>
<Key>sample3</Key>
<Code>AccessDenied</Code>
<Message>Access Denied</Message>
</Error>
<Error>
<Key>sample4</Key>
<Code>AccessDenied</Code>
<Message>Access Denied</Message>
</Error>
</DeleteResult>`,
`<DeleteResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
<Deleted>
<Key>sample5</Key>
<DeleteMarker>true</DeleteMarker>
<DeleteMarkerVersionId>null</DeleteMarkerVersionId>
</Deleted>
</DeleteResult>`))

		Assert(fn('sample_bucket', #('sample1', 'sample2', 'sample3',
			'sample4', 'sample5')) is:
			'Deleting files from Amazon S3 failed:\n\n\tsample3, sample4')
		Assert(_makeRequestRes isSize: 0)
		Assert(_count[0] is: 3)
		}

	Test_FileSize()
		{
		_callLog = Object()
		cl = AmazonS3
			{
			GetBucketLocationCached(bucket)
				{
				return bucket
				}
			ListBucketContents(bucket/*unused*/, file, region/*unused*/ = false)
				{
				result = false
				switch (file)
					{
				case 'empty':
					result = #()
				case 'retry_empty':
					if _callLog.Size() is 1
						result = #()
					else
						result = #(#(key: "retry_empty", size: "123 bytes"))
				case 'not_existing':
					result = #(#(key: 'file', size: "123 bytes"))
				default :
					result = Object(Object(key: file, size: "456 bytes"))
					}
				return result
				}
			AmazonS3_fileSizeSleep(i)
				{
				_callLog.Add(i)
				}
			}
		fn = cl.FileSize

		// empty
		file = 'empty'
		try
			fn('', file)
		catch (err)
			Assert(err is: 'Could not get file list for ' $ file)
		Assert(_callLog isSize: 2)
		Assert(_callLog[0] is: 0)
		Assert(_callLog[1] is: 1)
		_callLog.Delete(all:)

		// retry one time
		file = 'retry_empty'
		try
			result = fn('', file)
		Assert(_callLog isSize: 2)
		Assert(_callLog[0] is: 0)
		Assert(_callLog[1] is: 1)
		Assert(result is: 123)
		_callLog.Delete(all:)

		// not existing
		file = 'not_existing'
		try
			fn('', file)
		catch (err)
			Assert(err is: file $ ' does not exist')
		Assert(_callLog isSize: 1)
		Assert(_callLog[0] is: 0)
		_callLog.Delete(all:)

		// existing file
		file = 'an_existing_file'
		try
			result = fn('', file)
		Assert(result is: 456)
		Assert(_callLog isSize: 1)
		Assert(_callLog[0] is: 0)
		}

	Test_DeleteFolder()
		{
		mock = .initializeMock()
		mock.When.makeRequest([anyArgs:]).CallThrough()
		mock.When.DeleteFile([anyArgs:]).CallThrough()
		mock.When.GetBucketLocationCached([anyArgs:]).Return('fakeregion')
		mock.When.https([anyArgs:]).Return(AmazonTestReponses.httpsFailure)
		Assert({ mock.Eval(AmazonS3.DeleteFolder, 'NotARealBucket', 'test1') }
			throws: 'expected a string ending with')

		Assert(mock.Eval(AmazonS3.DeleteFolder, 'NotARealBucket', 'test1/') is: false)
		mock.Verify.https('DELETE',
			'https://s3.fakeregion.amazonaws.com/NotARealBucket/test1/', '', '',
			'signed request', '', '')

		Assert(mock.Eval(AmazonS3.DeleteFolder, 'NotARealBucket', 't1/t2/') is: false)
		mock.Verify.https('DELETE',
			'https://s3.fakeregion.amazonaws.com/NotARealBucket/t1/t2/', '', '',
			'signed request', '', '')
		}

	Test_CreateFolder()
		{
		mock = .initializeMock()
		mock.When.makeRequest([anyArgs:]).CallThrough()
		mock.When.signRequest([anyArgs:]).CallThrough()
		mock.When.https([anyArgs:]).Return(AmazonTestReponses.httpsReturn)
		Assert(mock.Eval(AmazonS3.CreateFolder, 'bucket', 'test1_got/'))
		mock.Verify.signRequest('EMPTYPUT', 'us-east-1', '', '/bucket/test1_got/',
			#(X_Amz_Content_Sha256: 'UNSIGNED-PAYLOAD',	X_Amz_Security_Token: 'token',
				Content_Length: 0))
		mock.Verify.https('EMPTYPUT',
			'https://s3.us-east-1.amazonaws.com/bucket/test1_got/',
			'', '', 'signed request', '','')
		}

	Test_ObjectMetaData()
		{
		mock = .initializeMock()
		mock.When.makeRequest([anyArgs:]).CallThrough()
		response = AmazonTestReponses.httpsMetaDataReturn.Copy()
		response.content = response.header
		mock.When.https([anyArgs:]).Return(response)
		Assert(mock.Eval(AmazonS3.ObjectMetaData, 'bucket', 'test-file', 'test-version')
			has: 'x-amz-storage-class')
		}

	Test_RestoreArchivedFile()
		{
		mock = .initializeMock()
		mock.When.makeRequest([anyArgs:]).CallThrough()
		mock.When.https([anyArgs:]).Return(AmazonTestReponses.httpRestoreArchivedReturn)
		Assert(mock.Eval(AmazonS3.RestoreArchivedFile, 'bucket', 'test-file',
			'test-version') has: '202 Accepted')
		}

	Test_checkResponseAndLog()
		{
		mock = .initializeMock()
		mock.When.addToLog([anyArgs:])
		resultOb = #(header: 'HTTP/1.1 200 OK')
		Assert(mock.Eval(AmazonS3.AmazonS3_checkResponseAndLog, 'PUT', resultOb))
		Assert(mock.Eval(
			AmazonS3.AmazonS3_checkResponseAndLog, 'PUT', resultOb, #('200', '202'))
			)
		resultOb = #(header: 'HTTP/1.1 202 OK')
		Assert(mock.Eval(
			AmazonS3.AmazonS3_checkResponseAndLog, 'PUT', resultOb, #('200', '202'))
			)
		resultOb = #(header: 'HTTP/1.1 404 OK')
		Assert(mock.Eval(
			AmazonS3.AmazonS3_checkResponseAndLog, 'PUT', resultOb, #('200', '202'))
				is: false)

		resultOb = #(header: '')
		Assert(mock.Eval(
			AmazonS3.AmazonS3_checkResponseAndLog, 'PUT', resultOb, #('200', '202'))
				is: false)
		mock.Verify.Times(1).addToLog('PUT',
			'Bad HTTP Status Code (Invalid HTTP response code in: (empty header))',
			'\r\n\r\n', '')
		}

	amazonUnprocessedFileList: #(
		#(key: "1.txt", size: "68 bytes", last_modified: #20250424.103009),
		#(key: "1093.pdf", size: "5978 bytes", last_modified: #20250424.083222),
		#(key: "1099_Test.pdf", size: "5978 bytes", last_modified: #20250424.083222),
		#(key: "202501/1.txt", size: "68 bytes", last_modified: #20250424.112849),
		#(key: "202501/BizUsers.txt", size: "378 bytes", last_modified: #20250424.115410),
		#(key: "202501/", size: "0 bytes", last_modified: #20250423.112849),
		#(key: "202501/aaaaccccc", size: "487 bytes", last_modified: #20250424.115958),
		#(key: "202501/jp.jpg", size: "65929 bytes", last_modified: #20250424.113032),
		#(key: "202501/openfilename", size: "421 bytes", last_modified: #20250424.115536),
		#(key: "202501/test.png", size: "45533 bytes", last_modified: #20250424.115827),
		#(key: "202502/1.txt", size: "68 bytes", last_modified: #20250424.112020),
		#(key: "202503/1.txt", size: "68 bytes", last_modified: #20250424.102603),
		#(key: "202503/yay.PNG", size: "203639 bytes", last_modified: #20250331.151029),
		#(key: "202504/1.doc", size: "68 bytes", last_modified: #20250424.091559),
		#(key: "202504/1.txt", size: "68 bytes", last_modified: #20250424.091306),
		#(key: "202505/Beach.jpg", size: "13677 bytes", last_modified: #20250516.114429),
		#(key: "202505/jpTest.pdf", size: "25194 bytes", last_modified: #20250502.080505),
		#(key: "202506/Order10.jpg", size: "1161 bytes", last_modified: #20250605.105327),
		#(key: "202507/Beach.jpg", size: "13677 bytes", last_modified: #20250704.115313),
		#(key: "5078.pdf", size: "5978 bytes", last_modified: #20250424.083222),
		#(key: "ACI Test 123.jpg", size: "99488 bytes", last_modified: #20250424.083223),
		#(key: "truck.gif", size: "20116 bytes", last_modified: #20250424.083214))

	expectedProcessedRootList: #(
		#(name: "1.txt", size: 68, date: #20250424.103009),
		#(name: "1093.pdf", size: 5978, date: #20250424.083222),
		#(name: "1099_Test.pdf", size: 5978, date: #20250424.083222),
		#(name: "202501/", size: 0, date: #20250423.112849),
		#(name: "202502/", size: 0, date: #20250424.112020),
		#(name: "202503/", size: 0, date: #20250331.151029),
		#(name: "202504/", size: 0, date: #20250424.091306),
		#(name: "202505/", size: 0, date: #20250502.080505),
		#(name: "202506/", size: 0, date: #20250605.105327),
		#(name: "202507/", size: 0, date: #20250704.115313),
		#(name: "5078.pdf", size: 5978, date: #20250424.083222),
		#(name: "ACI Test 123.jpg", size: 99488, date: #20250424.083223),
		#(name: "truck.gif", size: 20116, date: #20250424.083214))

	Test_mimicExeDir()
		{
		nameOnlyList = .amazonUnprocessedFileList.Map({ it.key })
		Assert(AmazonS3.AmazonS3_mimicExeDir('', nameOnlyList, false)
			equalsSet: #("1.txt", "1093.pdf", "1099_Test.pdf", "202501/", "202502/",
				"202503/", "202504/", "202505/", "202506/", "202507/", "5078.pdf",
				"ACI Test 123.jpg", "truck.gif"))

		Assert(AmazonS3.AmazonS3_mimicExeDir(
			'', .amazonUnprocessedFileList.DeepCopy(), true)
			is: .expectedProcessedRootList)

		nameOnlyList = .amazonUnprocessedFileList.Map({ 'att/' $ it.key })
		Assert(AmazonS3.AmazonS3_mimicExeDir('att/', nameOnlyList, false)
			equalsSet: #("1.txt", "1093.pdf", "1099_Test.pdf", "202501/", "202502/",
				"202503/", "202504/", "202505/", "202506/", "202507/", "5078.pdf",
				"ACI Test 123.jpg", "truck.gif"))

		list = .amazonUnprocessedFileList.DeepCopy().Each({ it.key = 'att/' $ it.key })
		Assert(AmazonS3.AmazonS3_mimicExeDir('att/', list, true)
			is: .expectedProcessedRootList)

		Assert(AmazonS3.AmazonS3_mimicExeDir('att/', #(), true) is: #())
		Assert(AmazonS3.AmazonS3_mimicExeDir('att/', #('att/'), false) is: #())
		Assert(AmazonS3.AmazonS3_mimicExeDir('att/',
			Object(Object(key: 'att/', last_modified: #20250916, size: '0')), true)
			is: #())
		Assert(AmazonS3.AmazonS3_mimicExeDir('att/', #(), false) is: #())
		Assert(AmazonS3.AmazonS3_mimicExeDir('', #(), true) is: #())
		Assert(AmazonS3.AmazonS3_mimicExeDir('', #(), false) is: #())

		// dir one file
		Assert(AmazonS3.AmazonS3_mimicExeDir('test.abc', #(), true) is: #())
		Assert(AmazonS3.AmazonS3_mimicExeDir('test.abc', #(), false) is: #())
		Assert(AmazonS3.AmazonS3_mimicExeDir('att/test.abc', #(), true) is: #())
		Assert(AmazonS3.AmazonS3_mimicExeDir('att/test.abc', #(), false) is: #())

		Assert(AmazonS3.AmazonS3_mimicExeDir('att/test.abc',
			#('att/test.abc', 'att/test.abcdefg'), false)
			is: #('test.abc'))
		Assert(AmazonS3.AmazonS3_mimicExeDir('att/test.abc',
			[[key: 'att/test.abcdefg', size: '2 bytes', last_modified: #20010101],
			[key: 'att/test.abc', size: '1 bytes', last_modified: #20000101]], true)
			is: #(#(name: 'test.abc', size: 1, date: #20000101)))
		}

	Test_formatDirPath()
		{
		fn = AmazonS3.AmazonS3_formatDirPath
		Assert(fn('c:/hello') is: 'c:/hello')
		Assert(fn('/srv/hello') is: '/srv/hello')

		Assert(fn('c:/hello/*.*') is: 'c:/hello/')
		Assert(fn('/srv/hello/*') is: '/srv/hello/')

		Assert({ fn('/srv/hello/') } throws: 'only accept *, *.* or file listing')
		Assert({ fn('c:/hello/') } throws: 'only accept *, *.* or file listing')
		Assert({ fn('c:/hello/*.jpg') } throws: 'wildcards are not handled')
		Assert({ fn('c:/hel*lo/*') } throws: 'wildcards are not handled')
		Assert({ fn('c:/hel?lo/*') } throws: 'wildcards are not handled')
		}

	Test_GetBucketRegion()
		{
		mock = Mock(AmazonS3)
		mock.When.GetBucketRegion([anyArgs:]).CallThrough()

		mock.When.makeRequest([anyArgs:]).Return()
		Assert({mock.GetBucketRegion('bucket0')} throws: 'member not found: 0')
		header = 'HTTP/1.1 200 OK\n' $
			'x-amz-bucket-region: us-east-1'
		mock.When.makeRequest([anyArgs:]).Return([:header])
		Assert(mock.GetBucketRegion('bucket1') is: 'us-east-1')

		header = 'HTTP/1.1 301 Moved Permanently\n' $
			'x-amz-bucket-region: ca-central-1'
		mock.When.makeRequest([anyArgs:]).Return([:header])
		Assert(mock.GetBucketRegion('bucket2') is: 'ca-central-1')

		mock.When.makeRequest([anyArgs:]).Return([headker: header.Lines()[0]])
		Assert({mock.GetBucketRegion('bucket3')}
			throws: 'method not found: boolean.Lower')

		header = 'HTTP/1.1 403 Forbidden\n' $
			'x-amz-bucket-region: us-west-2'
		mock.When.makeRequest([anyArgs:]).Return([:header])
		Assert(mock.GetBucketRegion('bucket3') is: 'us-west-2')
		}

	Test_throttleRetry()
		{
		fn = AmazonS3.AmazonS3_throttleRetry
		// inner retry: retries on 503
		attempts = 0
		result = fn({ attempts++
			if attempts < 3
				throw 'Bad HTTP Status Code (503)'
			'success' })
		Assert(result is: 'success')
		Assert(attempts is: 3)

		// outer retry: retries on curl (52)
		attempts = 0
		result = fn({ attempts++
			if attempts < 3
				throw 'curl: (52) Empty reply from server'
			'success' })
		Assert(result is: 'success')
		Assert(attempts is: 3)

		// both combined: curl once, then 503s, then success
		attempts = 0
		result = fn({
			attempts++
			if attempts is 1
				throw 'curl: (52) Empty reply from server'
			if attempts < 3
				throw 'Bad HTTP Status Code (503)'
			'success'
			})
		Assert(result is: 'success')
		Assert(attempts is: 3)

		// non-matching exception is NOT retried
		attempts = 0
		Assert({ fn({ attempts++ ; throw 'Some other error' }) }
			throws: 'Some other error')
		Assert(attempts is: 1)

		// exhausted retries throws Retry failed
		attempts = 0
		Assert({ fn({ attempts++ ; throw 'Bad HTTP Status Code (503)' }) }
			throws: 'Retry failed')
		Assert(attempts is: 3)
		}
	}
