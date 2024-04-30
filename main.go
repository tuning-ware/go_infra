package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		sgArgs := &ec2.SecurityGroupArgs{
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(8080),
					ToPort:     pulumi.Int(8080),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22),
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		}

		sg, err := ec2.NewSecurityGroup(ctx, "jenkins-sg", sgArgs)
		if err != nil {
			return err
		}

		kp, err := ec2.NewKeyPair(ctx, "local-ssh", &ec2.KeyPairArgs{
			PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7n+yFsr+ZuLEL/wL4azx2oiPNUvBYU+aUI2Q21OUH93I1chTDsnndtXPVWPMRG3xn6KMKRxJNCOkwhLS+guV/aBm/NYmSDpT2lFiWMqif1MH3P5TlIKeYFqjX0jMSZfW5BrCzYnyqRawv0tBFVoH1tRTBCqu3v8U8sd++A7ypIp/sbP/zXi5xchHGQd7yIwBcOdVLng7fUflouanvDUI3JnXYkjC0uYwatPm2u4hZPppcvtDFhZLFCTuf9qKF5g8Y0Occ438DOe6UyMidf6CP7hLuas76oTOPrCrm4cHZwW80jgFZd/etypC2phXlLW+7/zWQyPktqvT4WGiHghuWrBV8ARY/N5/eAr8P7eJswIF8m1i6Lr2lsXspYQwB2ZenobZIiKKnfoc8HFqZ0VkZpZjo15dUP3kHldfYRotqzuVrxa0Buu+qdLQ9+xPNN1wpsZxQZO12LYHNShowJeycL1muV5PbNTnnhAf6oJ84Oibn6JHWii1rV05mF0wuUY8= ecche@aero16elvin"),
		})
		if err != nil {
			return err
		}

		jenkinsServer, err := ec2.NewInstance(ctx, "jenkins-server", &ec2.InstanceArgs{
			InstanceType:        pulumi.String("t2.micro"),
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
			Ami:                 pulumi.String("ami-07caf09b362be10b8"),
			KeyName:             kp.KeyName,
		})

		fmt.Println(jenkinsServer.PublicIp)
		fmt.Println(jenkinsServer.PublicDns)

		ctx.Export("publicIp", jenkinsServer.PublicIp)
		ctx.Export("publicHostName", jenkinsServer.PublicDns)

		return nil
	})
}
