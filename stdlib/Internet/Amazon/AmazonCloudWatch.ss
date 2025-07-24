// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
AmazonAWS
	{
	region: 'us-east-1'
	Host(region = false)
		{
		return 'monitoring.' $ (region is false ? .region : region) $ '.amazonaws.com'
		}

	ContentType()
		{
		return 'application/x-www-form-urlencoded; charset=utf-8'
		}

	Service()
		{
		return 'monitoring'
		}

	CanonicalQueryString(params)
		{
		return params
		}

	PayloadHash(params)
		{
		params = '' // params is sent as query string, payload/body is empty
		return Sha256(params).ToHex()
		}

	makeRequest(call, params, path, region = false)
		{
		region = region is false ? .region : region
		params = AmazonAWS.UrlEncodeValues(params)
		url = 'https://' $ .Host(region) $ path $ Opt('?', params)
		extraHeaders = Object()
		extraHeaders.X_Amz_Security_Token = .SecurityToken()
		extraHeaders.Accept = 'application/json'

		if false is header = .signRequest(call, region, params, path, extraHeaders)
			return false
		return .https(call, url, header)
		}

	https(call, url, header = #())
		{
		return Https[call.Capitalize()](url, :header)
		}

	signRequest(call, region, params, path, extraHeaders)
		{
		region = region is false ? .region : region
		return AmazonV4Signing(this, call, region, params, path,
			extraHeaders).AuthorizationHeader()
		}

	GetBucketSize(bucket)
		{
		region = AmazonS3.GetBucketLocationCached(bucket)
		params = .generateParams(bucket, 'BucketSizeBytes', 'StandardStorage')
		res = .makeRequest('POST', params, '/', :region)
		return .parseResult(res)
		}

	GetBucketFilesCount(bucket)
		{
		region = AmazonS3.GetBucketLocationCached(bucket)
		params = .generateParams(bucket, 'NumberOfObjects', 'AllStorageTypes')
		res = .makeRequest('POST', params, '/', :region)
		return .parseResult(res)
		}

	generateParams(bucket, metricName, storage, type = 'Sum')
		{
		// 3 days as GetMetricStatistics only updated once a day and there is delay of 24-48 hours
		startTime = Date().Minus(days:3).GMTime().Format('yyyyMMddTHHmmssZ')
		endTime = Date().GMTime().Format('yyyyMMddTHHmmssZ')
		action = Object(
			Action: 'GetMetricStatistics'
			Version: '2010-08-01'
			Namespace: 'AWS/S3'
			MetricName: metricName
			'Dimensions.member.1.Name': 'BucketName'
			'Dimensions.member.1.Value':  bucket
			'Dimensions.member.2.Name': 'StorageType'
			'Dimensions.member.2.Value': storage
			StartTime:  startTime
			EndTime: endTime
			Period: '86400'
			'Statistics.member.1': type,
			'Format': 'json')
		return action
		}

	parseResult(json)
		{
		res = Json.Decode(json, 'skip')
		stats = res.GetDefault('GetMetricStatisticsResponse', #()).GetDefault(
			'GetMetricStatisticsResult', #()).GetDefault('Datapoints', #())
		if stats.Empty?()
			return 0

		return stats.Copy().Sort!(By('Timestamp')).Last().Sum
		}
	}