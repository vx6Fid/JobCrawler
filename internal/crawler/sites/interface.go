package sites

import (
	"context"

	"github.com/gocolly/colly/v2"
	"github.com/vx6fid/job-crawler/pkg"
)

var knownTechnologies = []string{
	"kafka", "kubernetes", "docker", "terraform", "ansible", "jenkins", "github actions",
	"ci/cd", "prometheus", "grafana", "datadog", "aws", "gcp", "azure", "python", "go", "java",
	"javascript", "react", "vue", "node.js", "typescript", "postgresql", "mysql", "mongodb",
	"redis", "elasticsearch", "linux", "bash", "shell", "git", "nginx", "apache", "sql", "nosql",
	"rest", "grpc", "graphql", "selenium", "puppet", "chef", "circleci", "travisci", "bitbucket",
	"openshift", "helm", "istio", "argocd", "flux", "zabbix", "new relic", "splunk", "pagerduty",
	"opsgenie", "cloudflare", "s3", "ec2", "iam", "cloudformation", "load balancer", "ecs", "eks",
	"fargate", "cloudwatch", "vpc", "lambda", "serverless", "tdd", "bdd", "junit", "pytest",
	"rspec", "mocha", "chai", "kafka streams", "kafka connect", "kinesis", "rabbitmq", "activemq",
	"celery", "airflow", "snowflake", "bigquery", "redshift", "zipkin", "jaeger",
}

type SiteParser interface {
	Matches(url string) bool
	Parse(ctx context.Context, e *colly.HTMLElement) ([]pkg.JobPosting, error)
	ParseJobDescription(e *colly.HTMLElement) (pkg.JobPosting, error)
}
