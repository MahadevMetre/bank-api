package ec2_scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/robfig/cron/v3"
)

type Ec2Scheduler struct {
	cron       *cron.Cron
	client     *ec2.Client
	instanceID []string
}

type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

func NewScheduler(creds AWSCredentials) (*Ec2Scheduler, error) {
	if creds.AccessKeyID == "" || creds.SecretAccessKey == "" || creds.Region == "" {
		return nil, fmt.Errorf("missing required AWS credentials")
	}

	credProvider := credentials.NewStaticCredentialsProvider(
		creds.AccessKeyID,
		creds.SecretAccessKey,
		"",
	)

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credProvider),
		config.WithRegion(creds.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatal(err)
	}

	return &Ec2Scheduler{
		cron: cron.New(
			cron.WithLocation(loc),
			cron.WithSeconds(),
		),
		client: client,
	}, nil
}

func (s *Ec2Scheduler) Start() error {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return fmt.Errorf("failed to load IST timezone: %w", err)
	}

	instances, err := s.GetInstancesByTagValue(context.Background(), "Name", "paydoh-redmine-server")
	if err != nil {
		return fmt.Errorf("failed to get instances: %w", err)
	}

	s.instanceID = instances

	// Morning job - Monday to Friday at 9:30 AM IST (4:30 UTC)
	_, err = s.cron.AddFunc("0 30 4 * * 1-5", func() {
		fmt.Printf("Running morning task for 'Monday to Friday at 9:30 AM IST': Current time is %s (IST)\n",
			time.Now().In(loc).Format("2006-01-02 15:04:05 MST"))
		s.StartInstances(context.TODO(), s.instanceID)
	})
	if err != nil {
		return fmt.Errorf("failed to add morning job: %w", err)
	}

	// Evening job - Monday to Friday at 7:00 PM IST (13:30 UTC)
	_, err = s.cron.AddFunc("0 0 13 * * 1-5", func() {
		fmt.Printf("Running evening task for 'Monday to Friday at 7:00 PM IST': Current time is %s (IST)\n",
			time.Now().In(loc).Format("2006-01-02 15:04:05 MST"))
		s.StopInstances(context.Background(), s.instanceID)
	})
	if err != nil {
		return fmt.Errorf("failed to add evening job: %w", err)
	}

	s.cron.Start()
	return nil
}

func (s *Ec2Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	fmt.Println("Scheduler stopped")
}

func (s *Ec2Scheduler) StartInstances(ctx context.Context, instanceIDs []string) error {
	if len(instanceIDs) == 0 {
		return fmt.Errorf("no instance IDs provided")
	}

	input := &ec2.StartInstancesInput{
		InstanceIds: instanceIDs,
	}

	_, err := s.client.StartInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to start instances: %v", err)
	}

	return nil
}

func (s *Ec2Scheduler) StopInstances(ctx context.Context, instanceIDs []string) error {
	if len(instanceIDs) == 0 {
		return fmt.Errorf("no instance IDs provided")
	}

	input := &ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	}

	_, err := s.client.StopInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to stop instances: %v", err)
	}

	return nil
}

func (s *Ec2Scheduler) GetInstancesByTagValue(ctx context.Context, tagKey, tagValue string) ([]string, error) {
	if tagKey == "" || tagValue == "" {
		return nil, fmt.Errorf("tag key and value must not be empty")
	}

	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String(fmt.Sprintf("tag:%s", tagKey)),
				Values: []string{tagValue},
			},
		},
	}

	result, err := s.client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %v", err)
	}

	var instanceIDs []string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if instance.InstanceId != nil {
				instanceIDs = append(instanceIDs, *instance.InstanceId)
			}
		}
	}

	return instanceIDs, nil
}
