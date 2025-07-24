// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_init()
		{
		init = AmazonS3_Multipart.AmazonS3_Multipart_init

		_makeRequest = function (@unused) { return false }
		Assert(init('test_bucket', 'test_file') is: false)

		_makeRequest = function (@unused)
			{ return `<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Bucket>example-bucket</Bucket>
  <Key>example-object</Key>
  <UploadId>VXBsb2FkIElEIGZvciA2aWWpbmcncyBteS1tb3ZpZS5tMnRzIHVwbG9hZA</UploadId>
</InitiateMultipartUploadResult>`
			}
		Assert(init('test_bucket', 'test_file')
			is: 'VXBsb2FkIElEIGZvciA2aWWpbmcncyBteS1tb3ZpZS5tMnRzIHVwbG9hZA')
		}

	Test_uploadParts()
		{
		mock = Mock(AmazonS3_Multipart)
		mock.When.splitFile([anyArgs:]).Return(3)

		_makeRequest = function (@unused) { return false }
		Assert(mock.Eval(AmazonS3_Multipart.AmazonS3_Multipart_uploadParts,
				'test_bucket', 'test_file', 'test_file', 'test_upload_id')
			is: false)
		mock.Verify.cleanupFiles([anyArgs:])

		_makeRequest = function (@args)
			{
			if args.fromFile.Suffix?('_p0')
				return Object(content: '', header: 'ETag: "upload0"')
			if args.fromFile.Suffix?('_p1')
				return Object(content: '', header: 'ETag: "upload1"')
			if args.fromFile.Suffix?('_p2')
				return Object(content: '', header: 'ETag: "upload2"')
			throw "should not get here"
			}
		Assert(mock.Eval(AmazonS3_Multipart.AmazonS3_Multipart_uploadParts,
				'test_bucket', 'test_file', 'test_file', 'test_upload_id')
			is: #(#(partNumber: "1", etag: '"upload0"'),
				#(partNumber: "2", etag: '"upload1"'),
				#(partNumber: "3", etag: '"upload2"')))
		mock.Verify.Times(2).cleanupFiles([anyArgs:])

		_makeRequest = function (@args)
			{
			if args.fromFile.Suffix?('_p0')
				return Object(content: '', header: 'ETag: "upload0"')
			if args.fromFile.Suffix?('_p1')
				return false
			throw "should not get here"
			}
		Assert(mock.Eval(AmazonS3_Multipart.AmazonS3_Multipart_uploadParts,
				'test_bucket', 'test_file', 'test_file', 'test_upload_id')
			is: false)
		mock.Verify.Times(3).cleanupFiles([anyArgs:])
		}

	Test_complete()
		{
		complete = AmazonS3_Multipart.AmazonS3_Multipart_complete
		_makeRequest = function(@args)
			{
			Assert(args.content.Tr('\r\n \t')
				is:
				'<CompleteMultipartUpload>
					<Part>
						<PartNumber>1</PartNumber>
						<ETag>&quot;upload0&quot;</ETag>
					</Part>
					<Part>
						<PartNumber>2</PartNumber>
						<ETag>&quot;upload1&quot;</ETag>
					</Part>
					<Part>
						<PartNumber>3</PartNumber>
						<ETag>&quot;upload2&quot;</ETag>
					</Part>
				</CompleteMultipartUpload>'.Tr('\r\n \t'))
			return false
			}
		etags = #(#(partNumber: "1", etag: '"upload0"'),
			#(partNumber: "2", etag: '"upload1"'),
			#(partNumber: "3", etag: '"upload2"'))
		Assert(complete('test_bucket', 'test_file', 'test_upload_id', etags)
			is: false)
		}

	Test_splitFile()
		{
		s = "helloworld".Repeat(100) // = 1000 bytes
		test = {|cn|
			file = .MakeFile(s)
			n = AmazonS3_Multipart.AmazonS3_Multipart_splitFile(file, cn)
			Assert(n is: (1000 / cn).Ceiling())
			Assert(Dir(file $ "_p*").Size() is: n)
			for i in ..n
				Assert(GetFile(file $ "_p" $ i) is: s[i * cn :: cn])
			}
		test(64)
		test(100)
		}
	}