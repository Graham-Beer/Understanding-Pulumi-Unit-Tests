package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type infrastructure struct {
	group  *ec2.SecurityGroup
	server *ec2.Instance
}

func createInfrastructure(ctx *pulumi.Context) (*infrastructure, error) {
	group, err := ec2.NewSecurityGroup(ctx, "web-secgrp", &ec2.SecurityGroupArgs{
		Ingress: ec2.SecurityGroupIngressArray{
			// Uncomment to fail a test:
			// ec2.SecurityGroupIngressArgs{
			// 	Protocol:   pulumi.String("tcp"),
			// 	FromPort:   pulumi.Int(22),
			// 	ToPort:     pulumi.Int(22),
			// 	CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			// },
			ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("tcp"),
				FromPort:   pulumi.Int(80),
				ToPort:     pulumi.Int(80),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	mostRecent := true
	ami, err := ec2.LookupAmi(ctx, &ec2.LookupAmiArgs{
		Filters: []ec2.GetAmiFilter{
			{
				Name:   "name",
				Values: []string{"ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"},
			},
		},
		Owners:     []string{"137112412989"},
		MostRecent: &mostRecent,
	})
	if err != nil {
		return nil, err
	}

	server, err := ec2.NewInstance(ctx, "web-server-www", &ec2.InstanceArgs{
		InstanceType:        pulumi.String("t2-micro"),
		VpcSecurityGroupIds: pulumi.StringArray{group.ID()}, // reference the group object above
		Ami:                 pulumi.String(ami.Id),
		// Comment out to fail a test:
		Tags: pulumi.StringMap{"Name": pulumi.String("webserver")},
	})
	if err != nil {
		return nil, err
	}

	return &infrastructure{
		group:  group,
		server: server,
	}, nil
}
