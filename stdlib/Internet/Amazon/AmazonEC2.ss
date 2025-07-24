// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
AmazonAWS
	{
	Host(region = 'us-east-1')
		{
		return 'ec2.' $ region $ '.amazonaws.com'
		}

	ContentType()
		{
		return 'application/x-www-form-urlencoded; charset=UTF-8'
		}

	Service()
		{
		return 'ec2'
		}

	CanonicalQueryString(unused)
		{
		return ''
		}

	PayloadHash(params)
		{
		return Sha256(params).ToHex()
		}

	requiredArgs: #(Version: '2016-11-15')
	StartSpotInstances(instanceSpecs, type)
		{
		instanceSpecs.Type = type
		instanceSpecs.Action = 'RequestSpotInstances'
		return .makePost(instanceSpecs)
		}

	CreateLaunchTemplateVersion(templateName, userData)
		{
		config = Object()
		config["Action"] = "CreateLaunchTemplateVersion"
		config["LaunchTemplateData.UserData"] = Base64.Encode(userData)
		config["LaunchTemplateName"] = templateName
		config["SourceVersion"] = '$Latest'
		response = .makePost(config)
		return XmlParser(response.content).launchtemplateversion.versionnumber.Text()
		}

	CreateFleet(templateName, version, availabilityZones, count, validUntil = false)
		{
		config = Object()
		config["Action"] = "CreateFleet"
		config["Type"] = "maintain"
		if validUntil isnt false
			config["ValidUntil"] = validUntil.GMTime().Format('yyyy-MM-ddTHH:mm:ssZ')
		config["LaunchTemplateConfigs.1.LaunchTemplateSpecification.LaunchTemplateName"] =
			templateName
		config["LaunchTemplateConfigs.1.LaunchTemplateSpecification.Version"] =
			version
		config["TargetCapacitySpecification.DefaultTargetCapacityType"] = 'spot'
		config["TargetCapacitySpecification.TotalTargetCapacity"] = String(count)
		config["SpotOptions.AllocationStrategy"] = 'price-capacity-optimized'
		for m, v in availabilityZones
			{
			namePrefix = "LaunchTemplateConfigs.1.Overrides." $ (m + 1)
			config[namePrefix $ ".AvailabilityZone"] = v.name
			config[namePrefix $ ".SubnetId"] = v.subnet
			}
		response = .makePost(config)
		Assert(Http.ResponseCode(response.header) is: '200',
			msg: 'CreateFleet failed - ' $ response.content)
		return XmlParser(response.content).fleetid.Text()
		}

	ModifyFleet(fleetId, count)
		{
		config = Object()
		config["Action"] = "ModifyFleet"
		config["FleetId"] = fleetId
		config["TargetCapacitySpecification.TotalTargetCapacity"] = String(count)
		response = .makePost(config)
		Assert(Http.ResponseCode(response.header) is: '200',
			msg: 'ModifyFleet failed - ' $ response.content)
		return response.content
		}

	DeleteFleets(fleetId)
		{
		config = Object()
		config["Action"] = "DeleteFleets"
		config["FleetId.1"] = fleetId
		config["TerminateInstances"] = 'true'
		response = .makePost(config)
		Assert(Http.ResponseCode(response.header) is: '200',
			msg: 'DeleteFleets failed - ' $ response.content)
		}

	DescribeFleetInstances(fleetId)
		{
		config = Object()
		config["Action"] = "DescribeFleetInstances"
		config["FleetId.1"] = fleetId
		response = .makePost(config)
		Assert(Http.ResponseCode(response.header) is: '200',
			msg: 'DescribeFleetInstances failed - ' $ response.content)
		return response.content
		}

	DescribeFleets(fleetId = false)
		{
		config = Object()
		config["Action"] = "DescribeFleets"
		if fleetId isnt false
			config["FleetId.1"] = fleetId
		response = .makePost(config)
		Assert(Http.ResponseCode(response.header) is: '200',
			msg: 'DescribeFleets failed - ' $ response.content)
		return response.content
		}

	GetSpotRequestInfo(requestId = '', extraFilters = #())
		{
		ob = Object(Action: 'DescribeSpotInstanceRequests')
		if requestId isnt ''
			ob['SpotInstanceRequestId.1'] = requestId
		.buildExtraFilters(extraFilters, ob)
		.makePost(ob)
		}

	buildExtraFilters(extraFilters, ob)
		{
		i = 1
		for m, v in extraFilters
			{
			ob["Filter." $ i $ ".Name"] = m
			ob["Filter." $ i $ ".Value"] = v
			++i
			}
		}

	GetInstanceInfo(requestId = '', extraFilters = #())
		{
		ob = Object(Action: 'DescribeInstances')
		if requestId isnt ''
			ob['InstanceId.1'] = requestId
		.buildExtraFilters(extraFilters, ob)
		.makePost(ob)
		}

	GetInstanceVolumes(instanceId, prevSnapshotId)
		{
		ob = Object(Action: 'DescribeVolumes',
			'Filter.1.Name':'attachment.instance-id',
			'Filter.1.Value': instanceId,
			'Filter.2.Name': 'snapshot-id',
			'Filter.2.Value': prevSnapshotId)
		.makePost(ob)
		}

	MakeSnapshotFromVolume(volumeId, desc = '')
		{
		ob = Object(Action: 'CreateSnapshot', VolumeId: volumeId, Description: desc)
		.makePost(ob)
		}

	DeleteSnapshot(snapshotId)
		{
		ob = Object(Action: 'DeleteSnapshot', SnapshotId: snapshotId)
		.makePost(ob)
		}

	DescribeInstances(instanceId = '', region = 'us-east-1', options = #())
		{
		ob = Object(Action: 'DescribeInstances')
		if instanceId isnt ''
			ob.InstanceId = instanceId
		ob.Merge(options)
		.makePost(ob, region)
		}

	LoopThroughInstances(region, block)
		{
		next = ""
		do
			{
			response = XmlParser(AmazonEC2.DescribeInstances(:region,
				options: Object(MaxResults: '25', NextToken: next)).content)
			next = response.nexttoken.Text().Trim()
			for item in (response.reservationset.item.instancesset.item.List())
				{
				block(item.instanceid.Text())
				}
			}
		while(next isnt "")
		}

	StartInstance(instanceId)
		{
		ob = Object(Action: 'StartInstances', 'InstanceId.1': instanceId)
		.makePost(ob)
		}

	// 'HttpTokens': 'required'      -- for IMDSv2
	// always returns instance all metadata options
	ModifyInstanceMetadataOptions(instanceId, options = #(), region = 'us-east-1')
		{
		ob = Object(Action: 'ModifyInstanceMetadataOptions', 'InstanceId': instanceId)
		ob.Merge(options)
		response = .makePost(ob, region)
		options = Object()
		for o in XmlParser(response.content).instancemetadataoptions[0].Children()
			options[o.Name()] = o.Text()
		return options
		}

	StopInstance(instanceId)
		{
		ob = Object(Action: 'StopInstances', 'InstanceId.1': instanceId)
		.makePost(ob)
		}

	TerminateInstance(@instances)
		{
		if instances.Empty?()
			return
		ob = Object(Action: 'TerminateInstances')
		for i, id in instances
			ob['InstanceId.' $ (i + 1)] = id
		.makePost(ob)
		}

	CancelSpotRequest(@requestIds)
		{
		if requestIds.Empty?()
			return
		ob = Object(Action: 'CancelSpotInstanceRequests')
		for i, id in requestIds
			ob['SpotInstanceRequestId.' $ (i + 1)] = id
		.makePost(ob)
		}

	CreateTags(tags)
		{
		ob = Object(Action: 'CreateTags').MergeNew(tags)
		return .makePost(ob)
		}

	DescribeAvailabilityZone(filters = #())
		{
		ob = Object(Action: 'DescribeAvailabilityZones').MergeNew(filters)
		.makePost(ob)
		}

	makePost(callSpecs, region = 'us-east-1')
		{
		return .postRequest(AmazonAWS.UrlEncodeValues(callSpecs.MergeNew(.requiredArgs)),
			:region)
		}

	postRequest(message, region)
		{
		extraHeaderInfo = Object(X_Amz_Security_Token: .SecurityToken())
		header = AmazonV4Signing(this, 'POST', region, message,
			:extraHeaderInfo).AuthorizationHeader()
		result = Https('POST', 'https://' $ .Host(region), content: message, :header)
		return result
		}
	}
