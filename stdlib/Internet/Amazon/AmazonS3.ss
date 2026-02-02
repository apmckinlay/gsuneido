// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
AmazonAWS
	{
	S3ExpirySeconds: 60

	Host(region = false)
		{
		return 's3.' $ (region is false ? .region : region) $ '.amazonaws.com'
		}

	ContentType()
		{
		return 'multipart/form-data'
		}

	Service()
		{
		return 's3'
		}

	CanonicalQueryString(params)
		{
		return params
		}

	PayloadHash(unused)
		{
		return 'UNSIGNED-PAYLOAD'
		}

	region: 'us-east-1'
	makeRequest(call, params, path, toFile = '', fromFile = '', expectedResponse = '200',
		policy = '', fullResponse? = false, content = '', extraHeaders = #(),
		region = false, limitRate = '')
		{
		region = region is false ? .region : region
		params = AmazonAWS.UrlEncodeValues(params)
		url = 'https://' $ .Host(region) $ path $ Opt('?', params)
		extraHeaders = extraHeaders.Copy()
		extraHeaders.X_Amz_Content_Sha256 = .PayloadHash('')
		extraHeaders.X_Amz_Security_Token = .SecurityToken()
		if policy isnt ''
			extraHeaders.X_Amz_Acl = policy.AfterLast(':')
		if false is header = .signRequest(call, region, params, path, extraHeaders)
			return false

		return .throttle(call, :expectedResponse, :fullResponse?)
			{
			.https(call, url, toFile, fromFile, header, content, limitRate)
			}
		}

	signRequest(call, region, params, path, extraHeaders)
		{
		region = region is false ? .region : region
		return AmazonV4Signing(this, .mapSignedAction(call), region, params, path,
			extraHeaders).AuthorizationHeader()
		}

	mapSignedAction(call)
		{
		// creating a "Folder" on amazon involves putting an empty object
		// need to handle content/fromFile being empty
		// canonical still needs to use "PUT"
		return #("EMPTYPUT": "PUT").GetDefault(call, call)
		}

	https(call, url, toFile = '', fromFile = '', header = #(), content = '',
		limitRate = '')
		{
		return Https(call, url, :content, :toFile, :fromFile, :header,
			timeoutConnect: 60, :limitRate)
		}

	GetBucketLocationCached(bucket)
		{
		return Memoize
			{
			Func(bucket)
				{
				AmazonS3.GetBucketLocation(bucket)
				}
			}(bucket)
		}

	GetBucketLocation(bucket)
		{
		path = '/' $ bucket $ '/'
		if false is res = .makeRequest('GET', [location: '1'], path)
			return false

		resXml = XmlParser(res)
		if #() is region = resXml.Children()
			return .region
		return String(region[0])
		}

	CopyFile(bucketFrom, fileFrom, bucketTo, fileTo, policy = '')
		{
		region = .GetBucketLocationCached(bucketTo)
		fileFrom = Url.EncodePreservePath(fileFrom)
		fileTo = Url.EncodePreservePath(fileTo)
		path = '/' $ bucketTo $ '/' $ fileTo
		extraHeaders = Object()
		extraHeaders.X_Amz_Copy_Source = '/' $ bucketFrom $ '/' $ fileFrom
		return false isnt .makeRequest(
			'EMPTYPUT', [], path, :region, :extraHeaders, :policy)
		}

	GetFile(bucket, fileFrom, fileTo = "", exceptionOnFailure = false, versionId = '',
		region = false)
		{
		// normally failures are automatically logged and false is returned,
		// the exceptionOnFailure option allows the calling application to handle an
		// exception when the request fails
		_exceptionOnFailure = exceptionOnFailure
		if fileTo is ''
			fileTo = Paths.Basename(fileFrom)

		fileFrom = Url.EncodePreservePath(fileFrom)
		path = '/' $ bucket $ '/' $ fileFrom
		retVal = .makeRequest('GET', [:versionId], path, toFile: fileTo,
			:region) isnt false
		.cleanupFileIfFailed(retVal, fileTo)
		return retVal
		}

	PresignUrl(bucket, fileFrom, region = false, method = 'GET', download? = false)
		{
		path = '/' $ bucket $ '/' $ Url.EncodePreservePath(fileFrom)
		region = region is false ? .GetBucketLocationCached(bucket) : region
		return AmazonV4Signing(this, .mapSignedAction(method), region, '', path,
			#()).PresignUrl(.Host(region), :download?)
		}

	PutFileContent(bucket, fileName, content)
		{
		return false isnt .makeRequest('PUT', [], '/' $ bucket $ '/' $ fileName, :content)
		}

	PutFile(bucket, fileFrom, fileTo = '', policy = '', region = false,
		allowEmpty? = false, limitRate = '')
		{
		if false is fileTo = .validateAndGetFileTo(fileFrom, fileTo, allowEmpty?)
			return false

		path = '/' $ bucket $ '/' $ fileTo
		// toFile needs to be '' when sending, it is already set in the url path
		// if in request will change where curl stores the file prior to sending
		return false isnt
			.makeRequest('PUT', [], path, toFile: '', fromFile: fileFrom, :policy,
				:region, :limitRate)
		}

	validateAndGetFileTo(fileFrom, fileTo, allowEmpty? = false)
		{
		if not .validFile?(fileFrom, allowEmpty?)
			{
			.addToLog('PUT', fileFrom $ ' File does not exist or File is empty')
			return false
			}
		if fileTo is ''
			fileTo = Paths.Basename(fileFrom)
		fileTo = Url.EncodePreservePath(fileTo)
		return fileTo
		}

	validFile?(fileFrom, allowEmpty? = false)
		{
		return FileExists?(fileFrom) and (allowEmpty? or FileSize(fileFrom) isnt 0)
		}

	DeleteFile(bucket, filename, versionId = [])
		{
		Assert(filename isnt: '')
		Assert(filename isString: true)
		region = .GetBucketLocationCached(bucket)
		path = '/' $ bucket $ '/' $ Url.EncodePreservePath(filename)
		return .makeRequest('DELETE', versionId, path, expectedResponse: '204',
			:region) isnt false
		}

	ListBucketFolderFiles(bucket, folder = '')
		{
		// only show files at the current level (do not return files in any subfolders)
		// * Currently only handles one folder level.
		files = .ListBucketFiles(bucket, folder)
		return folder isnt '' ? files : files.Filter( { not it.Has?('/') })
		}

	Dir(bucket, dir, details = false)
		{
		dir = .formatDirPath(dir)
		region = AmazonS3.GetBucketLocationCached(bucket)
		if not details
			return .mimicExeDir(
				dir, AmazonS3.ListBucketFiles(bucket, dir, :region), details)

		return .mimicExeDir(
			dir, AmazonS3.ListBucketContents(bucket, dir, :region), details)
		}

	formatDirPath(dir)
		{
		dir = Paths.ToStd(dir)
		Assert(not dir.Suffix?('/'), msg: 'only accept *, *.* or file listing')
		dir = dir.Replace('(?q)/*.*(?-q)$', '/').Replace('(?q)/*(?-q)$', '/')
		Assert(dir hasnt: '*', msg: 'wildcards are not handled')
		Assert(dir hasnt: '?', msg: 'wildcards are not handled')
		return dir
		}

	// Match the result of Dir call in the same format we would get from the Exe
	mimicExeDir(dir, fileList, details)
		{
		if dir isnt '' and not dir.Suffix?('/')
			return .dirOne(dir, fileList, details)

		if not details
			return fileList.Map(
				{
				path = it.RemovePrefix(dir)
				if path.Has?('/')
					path = path.BeforeFirst('/') $ '/'
				path
				}).UniqueValues().Remove('')

		formattedfileList = Object()
		for f in fileList
			{
			f.key = f.key.RemovePrefix(dir)
			if f.key.Has?('/')
				{
				f.key = f.key.BeforeFirst('/') $ '/'
				if false is f1 = formattedfileList.FindOne({ it.name is f.key })
					formattedfileList.Add(.formatDirOb(f, size: 0))
				else
					f1.date = Min(f1.date, f.last_modified)
				}
			else
				formattedfileList.Add(.formatDirOb(f))
			}
		return formattedfileList.RemoveIf({ it.name is "" })
		}

	dirOne(file, fileList, details)
		{
		if not details
			return fileList.Filter({ it is file }).Map!(.dirFileName)

		return fileList.Filter({ it.key is file }).
			Map!({
				it.key = .dirFileName(it.key)
				.formatDirOb(it)
			})
		}

	dirFileName(file)
		{
		return file.Has?('/') ? file.AfterFirst('/') : file
		}

	formatDirOb(f, size = false)
		{
		return Object(
			name: f.key,
			date: f.last_modified,
			size: size is false ? Number(f.size.Tr(' bytes')) : size)
		}

	ListBucketFiles(bucket, prefix = '', region = false)
		{
		return .ListBucketContents(bucket, prefix, :region).
			Map({it.key}).
			RemoveIf({it.Suffix?('/')}) // remove subfolders
		}

	FileExist?(bucket, file, region = false)
		{
		if region is false
			region = .GetBucketLocationCached(bucket)
		list = .ListBucketFiles(bucket, file, :region)
		return list.Has?(file)
		}

	FileSize(bucket, file)
		{
		region = .GetBucketLocationCached(bucket)
		list = .ListBucketContents(bucket, file, :region)
		if false is f = list.FindOne({ it.key is file })
			throw file $ ' does not exist'
		return Number(f.size.Tr(' bytes'))
		}

	CreateBucket(bucket, region = false)
		{
		fromFile = GetAppTempFullFileName('amazons3')
		content = ''
		// If you don't specify a Region, the bucket is created in the
		// US East (N. Virginia) Region (us-east-1) by default.
		if region isnt false and region isnt 'us-east-1'
			content = '<CreateBucketConfiguration ' $
				'xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
				<LocationConstraint>' $ region $ '</LocationConstraint>
				</CreateBucketConfiguration >'
		PutFile(fromFile, content)
		result = .makeRequest('PUT', [], '/' $ bucket, :fromFile)
		DeleteFile(fromFile)
		return result isnt false
		}

	defaultCors: `<?xml version="1.0" encoding="UTF-8"?>
		<CORSConfiguration xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
		   <CORSRule>
			  <AllowedHeader>*</AllowedHeader>
			  <AllowedMethod>PUT</AllowedMethod>
			  <AllowedMethod>GET</AllowedMethod>
			  <AllowedMethod>DELETE</AllowedMethod>
			  <AllowedOrigin>*</AllowedOrigin>
			  <ExposeHeader>Content-Type</ExposeHeader>
			  <ExposeHeader>Content-Length</ExposeHeader>
			  <ExposeHeader>Content-Disposition</ExposeHeader>
		   </CORSRule>
		</CORSConfiguration>`
	PutBucketCors(bucket, region, cors = false)
		{
		content = cors is false ? .defaultCors : cors
		extraHeaders = Object('Content-MD5': Base64.Encode(Md5(content)))
		return '' is .makeRequest(
			'PUT', ['cors': '1'], '/' $ bucket, :content, :region, :extraHeaders)
		}

	GetBucketCors(bucket)
		{
		region = .GetBucketLocationCached(bucket)
		return .makeRequest('GET', ['cors': '1'], '/' $ bucket, :region)
		}

	PutBucketVersioning(bucket, region = false)
		{
		if region is false
			region = .GetBucketLocationCached(bucket)
		content = `<VersioningConfiguration ` $
			`xmlns="http://s3.amazonaws.com/doc/2006-03-01/">` $
			`<Status>Enabled</Status>` $
			`</VersioningConfiguration>`
		return '' is .makeRequest('PUT', ['versioning': '1'], '/' $ bucket,
			:content, :region)
		}

	GetBucketVersioning(bucket)
		{
		region = .GetBucketLocationCached(bucket)
		return .makeRequest('GET', ['versioning': '1'], '/' $ bucket, :region)
		}

	lifeCycleRule: `<Rule>` $
		`<ID>DeleteOldVersions</ID>` $
		`<Filter><Prefix></Prefix></Filter>` $
		`<Status>Enabled</Status>` $
		`<NoncurrentVersionExpiration>` $
		`<NoncurrentDays>14</NoncurrentDays>` $
		`</NoncurrentVersionExpiration>` $
		`</Rule>`
	PutBucketLifecycleConfig(bucket, rules = false, region = false)
		{
		if region is false
			region = .GetBucketLocationCached(bucket)
		rules = rules is false ? .lifeCycleRule : rules
		content = `<LifeCycleConfiguration>` $
			rules $ `</LifeCycleConfiguration>`
		extraHeaders = Object('Content-MD5': Base64.Encode(Md5(content)))
		return '' is .makeRequest('PUT', ['lifecycle': '1'], '/' $ bucket,
			:content, :region, :extraHeaders)
		}

	GetBucketLifecycleConfig(bucket)
		{
		region = .GetBucketLocationCached(bucket)
		return .makeRequest('GET', ['lifecycle': '1'], '/' $ bucket, :region)
		}

	PutBucketPolicy(bucket, region, policy)
		{
		if region is false
			region = .GetBucketLocationCached(bucket)
		return '' is .makeRequest('PUT', ['policy': '1'], '/' $ bucket,
			content: policy, :region, expectedResponse: '204')
		}

	GetBucketPolicy(bucket)
		{
		region = .GetBucketLocationCached(bucket)
		return .makeRequest('GET', ['policy': '1'], '/' $ bucket, :region)
		}

	CreateFolder(bucket, folder)
		{
		Assert(folder endsWith: '/')
		path = '/' $ bucket $ '/' $ Url.EncodePreservePath(folder)
		return false isnt
			.makeRequest('EMPTYPUT', [], path, extraHeaders: #(Content_Length: 0))
		}

	ListBucketFileVersions(bucket, prefix = '', rawResponse = false)
		{
		fileList = Object()
		params = [marker: '', :prefix, versions: '1']
		while true is .listTruncated(bucket, fileList, params, rawResponse, 'version')
			{
			last = fileList.Last()
			params.marker = last.key
			}
		return fileList
		}

	GetFileTagging(bucket, filename, versionId = '')
		{
		filename = Url.EncodePreservePath(filename)
		xml = .makeRequest('GET', [tagging: '1', :versionId],
			'/' $ bucket $ '/'$ filename)
		parsed = XmlParser(xml)
		fileTags = Object()
		for tag in parsed['tagset']['tag'].List() // only using booleans atm
			fileTags[tag.key.Text()] = tag.value.Text() is 'true'
		return fileTags
		}

	GetFileContents(bucket, filename)
		{
		filename = Url.EncodePreservePath(filename)
		result = .makeRequest('GET', [], '/' $ bucket $ '/'$ filename)
		return result
		}

	xmlHeader: '<?xml version="1.0" encoding="UTF-8"?>
		<PLACEHOLDER xmlns="http://s3.amazonaws.com/doc/2006-03-01/">'
	PutFileTagging(bucket, filename, tags, versionId = '')
		{
		filename = Url.EncodePreservePath(filename)
		xml = .xmlHeader.Replace('PLACEHOLDER', 'Tagging') $ '<TagSet>'
		for tag in tags.Members()
			{
			xml $= '<Tag>
				<Key>' $ tag $ '</Key>
				<Value>' $ tags[tag] $ '</Value>
			</Tag>'
			}
		xml $= '</TagSet></Tagging>'
		return .makeRequest('PUT', [tagging: '1', :versionId],
			'/' $ bucket $ '/'$ filename, content: xml)

		}

	ArchivedDBDaysKept: 5
	RestoreArchivedFile(bucket, filename, versionId)
		{
		filename = Url.EncodePreservePath(filename)
		xml = .xmlHeader.Replace('PLACEHOLDER', 'RestoreRequest') $
			'<Days>' $ Display(.ArchivedDBDaysKept) $ '</Days></RestoreRequest>'

		return .makeRequest('POST', [restore: '1', :versionId],
			'/' $ bucket $ '/'$ filename, content: xml, expectedResponse: #('200', '202'))
		}

	ObjectMetaData(bucket, filename, versionId = '')
		{
		filename = Url.EncodePreservePath(filename)
		return .makeRequest('HEAD', [:versionId], '/' $ bucket $ '/'$ filename,
			fullResponse?:).header
		}

	CheckObjectTagging(bucket, key, versionId)
		{
		key = Url.EncodePreservePath(key)
		return .makeRequest('GET', [tagging: '1', :versionId], '/' $ bucket $ '/'$ key)
		}

	DeleteFolder(bucket, folder)
		{
		Assert(folder endsWith: '/')
		.DeleteFile(bucket, folder)
		}

	ListBuckets()
		{
		content = .makeRequest('GET', [], '/')
		response = XmlParser(content)
		buckets = Object()
		for b in response.buckets.bucket.List()
			buckets.Add(b.name.Text())
		return buckets
		}

	ListBucketFolders(bucket, prefix = '')
		{
		return .ListBucketContents(bucket, prefix).
			Map({ it.key }).
			RemoveIf({ not it.Suffix?('/')}) // remove files
		}

	ListBucketContents(bucket, prefix = '', rawResponse = false, region = false)
		{
		// AmazonS3: The maximum number of items that can be returned is 1000
		// if IsTruncated is true then send another request
		// with the marker parameter set to the last item returned by the previous request
		// keep making these requests until the IsTruncated element has a value of false
		fileList = Object()
		params = [marker: '', :prefix]
		while true is .listTruncated(bucket, fileList, params, rawResponse, :region)
			{
			last = fileList.Last()
			params.marker = last.key
			}
		return fileList
		}

	listTruncated(bucket, fileList, params, rawResponse, xmlContent = 'contents',
		region = false)
		{
		if false is content = .makeRequest('GET', params, '/' $ bucket $ '/', :region)
			return false

		try
			response = XmlParser(content)
		catch (err)
			{
			.addToLog("LIST", 'listTruncated - ' $ err, content, params: Locals(0))
			return false
			}
		if response is false or response is ''
			return false

		for file in response[xmlContent].List()
			{
			last_modified = .getLastModified(rawResponse, file.lastmodified.Text())
			size = file.size.Text()
			if not rawResponse
				size $= " bytes"

			fileOb = Object(key: file.key.Text(),
				owner: file.owner.displayname.Text(),
				:size, :last_modified)

			.addXmlTypeIfNotEmpty(file, fileOb, 'storageclass', 'storage')
			.addXmlTypeIfNotEmpty(file, fileOb, 'versionid', 'versionId')

			fileList.Add(fileOb)
			}
		return response.istruncated.Text().SafeEval()
		}

	addXmlTypeIfNotEmpty(file, fileOb, type, member)
		{
		if file[type].Text() isnt ''
			fileOb[member] = file[type].Text()
		}

	getLastModified(rawResponse, dateStr)
		{
		if rawResponse
			return dateStr
		if false isnt date = Date(dateStr.Replace('.000Z', '').Trim())
			return date.GMTimeToLocal()
		return ''
		}

	throttle(action, block, expectedResponse = '200', fullResponse? = false)
		{
		resultOb = false
		try
			.throttleRetry()
				{
				resultOb = (block)()
				params = Locals(3) /*= call levels up to throttle */
				status = .checkResponseAndLog(action, resultOb, expectedResponse, params)
				return status is true
					? fullResponse?
						? resultOb
						: resultOb.GetDefault(#content, '').Trim()
					: false
				}
		catch (err, "Retry failed")
			{
			detail = resultOb is false ? ''
				: resultOb.header $ '\r\n\r\n' $ resultOb.GetDefault(#content, '')
			.addToLog(action, err, detail, params)
			return false
			}
		return false
		}

	DeleteFiles(bucket, filesToBeRemoved)
		{
		failed = Object()
		for file in filesToBeRemoved
			if false is .DeleteFile(bucket, file)
				failed.Add(file)

		return failed.Empty?() ? ""
			: "Deleting files from Amazon S3 failed:\n\n\t" $ failed.Join(', ')
		}

	deleteFileLimit: 1000
	DeleteFiles2(bucket, files)
		{
		failed = Object()
		for (i = 0; i < files.Size(); i = i + .deleteFileLimit)
			.deleteFiles(bucket, files[i :: .deleteFileLimit], failed)

		return failed.Empty?()
			? ''
			: 'Deleting files from Amazon S3 failed:\n\n\t' $ failed.Join(', ')
		}

	deleteFiles(bucket, files, failed)
		{
		content = .xmlHeader.Replace('PLACEHOLDER', 'Delete')
		for file in files
			{
			Assert(file isnt: '')
			Assert(file isString: true)
			content $= `<Object><Key>` $ file $ `</Key></Object>`
			}
		content $= '</Delete>'

		region = .GetBucketLocationCached(bucket)
		extraHeaders = Object('Content-MD5': Base64.Encode(Md5(content)))
		if false is res = .makeRequest('POST', ['delete': '1'], '/' $ bucket, :content,
			:region, :extraHeaders)
			{
			failed.Append(files)
			return
			}

		res = XmlParser(res)
		for child in res.Children()
			if child.Name() is 'error'
				failed.Add(child['key'].Text())
		}

	// extrated for tests
	throttleRetry(block)
		{
		Retry(block, maxRetries: 3, minDelayMs: 100,
			retryException: 'Bad HTTP Status Code (503)')
		}

	cleanupFileIfFailed(retVal, fileTo)
		{
		if retVal is true
			return
		DeleteFile(fileTo)
		}

	checkResponseAndLog(action, resultOb, expectedResponse = '200', params = '')
		{
		try
			code = Http.ResponseCode(resultOb.header)
		catch (e, 'Invalid HTTP response code')
			code = e

		if code is expectedResponse or
			Object?(expectedResponse) and expectedResponse.Has?(code)
			return true
		if action is 'GET' and code is '404'
			return false
		msg = 'Bad HTTP Status Code (' $ code $ ')'
		// 503 is either 'SlowDown' or 'ServiceUnavailable' - in both cases
		// amazon requsts reducing the request rate (switch to Exponential Fallback)
		if code is '503'
			throw msg
		detail = resultOb.header $ '\r\n\r\n' $ resultOb.GetDefault(#content, '')
		.addToLog(action, msg, detail, params)
		return false
		}

	addToLog(action, msg, detail = '', params = '', _secureLogging = false)
		{
		if Object?(params)
			params = params.Copy().Delete(#resultOb, #result)
		if .exceptionInsteadOfLog?()
			throw 'AmazonS3 - ' $ action $ ': ' $ msg
		SuneidoLog('ERRATIC: AmazonS3 - ' $ action $ ': ' $ msg, calls:,
			:params, switch_prefix_limit: 5)

		// can't be sure what "secure" data might be in msg or detail so don't log.
		if secureLogging
			return

		path = GetContributions('LogPaths').GetDefault('amazonS3', 'logs/amazonS3log')
		log = action $ ' > ' $ msg $ Opt('\r\n', detail.Trim()) $ '\r\n'
		.Log(path, log)
		}

	exceptionInsteadOfLog?()
		{
		try
			return _exceptionOnFailure
		return false
		}

	PutMultipartFile(bucket, fileFrom, fileTo = '')
		{
		if false is fileTo = .validateAndGetFileTo(fileFrom, fileTo)
			return false
		return AmazonS3_Multipart.PutFile(bucket, fileFrom, fileTo, .makeRequest)
		}

	notificationConfig: `<NotificationConfiguration ` $
		`xmlns="http://s3.amazonaws.com/doc/2006-03-01/">` $
		`<QueueConfiguration>` $
		`<Id>ID_PLACEHOLDER</Id>`$
		`<Queue>ARN_PLACEHOLDER</Queue>` $
		`<Event>s3:ObjectCreated:*</Event>`$
		`</QueueConfiguration>` $
		`</NotificationConfiguration>`

	PutNotificationConfig(bucket, id, arn, content = false, region = false)
		{
		region = region is false ? .GetBucketLocationCached(bucket) : region
		if content is false
			{
			content = .notificationConfig.Replace('ID_PLACEHOLDER', id)
			content = content.Replace('ARN_PLACEHOLDER', arn)
			}
		return '' is .makeRequest(
			'PUT', ['notification': '1'], '/' $ bucket, :content, :region)
		}

	GetNotificationConfig(bucket)
		{
		region = .GetBucketLocationCached(bucket)
		return .makeRequest('GET', ['notification': '1'], '/' $ bucket, :region)
		}
	}
