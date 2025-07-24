// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
// NOTE: It's a good practice to configure a lifecycle rule in the bucket properties
// to abort incomplete multipart uploads
class
	{
	PutFile(bucket, fileFrom, fileTo, makeRequest)
		{
		_makeRequest = makeRequest
		if false is uploadId = .init(bucket, fileTo)
			return false
		if false is etags = .uploadParts(bucket, fileFrom, fileTo, uploadId)
			return .abort(bucket, fileTo, uploadId)
		else
			return .complete(bucket, fileTo, uploadId, etags)
		}

	init(bucket, fileTo)
		{
		// it seems s3 expects "?uploads=" and value does not matter
		content = _makeRequest("POST", [uploads: 'x'], '/' $ bucket $ '/' $ fileTo)
		return String?(content) and content isnt ''
			? XmlParser(content).uploadid.Text()
			: false
		}

	uploadParts(bucket, fileFrom, fileTo, uploadId)
		{
		partCount = .splitFile(fileFrom, 8.Mb()) /*= reasonable part size */
		etags = Object()
		for i in ..partCount
			{
			partNumber = String(i+1)
			response = Object(content: false)
			result = RetryBool(maxretries: 3, min: 10)
				{
				response = Object(content: false)
				try response = _makeRequest('PUT',
					[:partNumber, :uploadId],
					'/' $ bucket $ '/' $ fileTo,
					fromFile: fileFrom $ '_p' $ i,
					fullResponse?:)
				Object?(response) and response.content is ''
				}
			if result is false
				{
				.cleanupFiles(fileFrom, partCount)
				return false
				}
			etag = InetMesg(response.header).Field('ETag')
			etags.Add(Object(:partNumber, :etag))
			}
		.cleanupFiles(fileFrom, partCount)
		return etags
		}

	splitFile(fileFrom, partSize)
		{
		partCount = (FileSize(fileFrom) / partSize).Ceiling()
		File(fileFrom)
			{ |src|
			for i in ..partCount
				File(fileFrom $ '_p' $ i, 'w')
					{ |dst|
					unused = src.CopyTo(dst, partSize) // last part may be smaller
					}
			}
		return partCount
		}

	cleanupFiles(fileFrom, partCount)
		{
		for i in ..partCount
			DeleteFile(fileFrom $ '_p' $ i)
		}

	complete(bucket, fileTo, uploadId, etags)
		{
		content = Razor('<CompleteMultipartUpload>
			@for(part in .etags)
				{
				<Part>
					<PartNumber>@part.partNumber</PartNumber>
					<ETag>@part.etag</ETag>
				</Part>
				}
			</CompleteMultipartUpload>', Object(:etags))
		response = _makeRequest('POST', [:uploadId], '/' $ bucket $ '/' $ fileTo,:content)
		return String?(response) and response isnt ''
		}

	abort(bucket, fileTo, uploadId)
		{
		_makeRequest('DELETE', [:uploadId], '/' $ bucket $ '/' $ fileTo)
		return false
		}
	}