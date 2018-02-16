# AWS related
variable region {
  type = "string"
}

variable availability_zones {
  type = "list"
}

variable image_id {
  type = "string"
}

variable subnets {
  type = "list"
}

variable certificate_arn {
  type = "string"
}

# Environment variables
variable consumer_key {
  type = "string"
}

variable consumer_secret {
  type = "string"
}

variable callback_url {
  type = "string"
}

variable port {
  type = "string"
}

variable homepage {
  type = "string"
}

variable host {
  type = "string"
}

variable protocol {
  type = "string"
}

provider "aws" {
  version = "~> 1.9"
  region  = "${var.region}"
}

resource "aws_default_vpc" "default" {
  tags {
    Name = "Default VPC"
  }
}

resource "aws_security_group" "lb" {
  name   = "lb"
  vpc_id = "${aws_default_vpc.default.id}"

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "ecs" {
  name   = "ecs"
  vpc_id = "${aws_default_vpc.default.id}"

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = ["${aws_security_group.lb.id}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_ecs_cluster" "backend" {
  name = "backend"
}

resource "aws_launch_configuration" "backend" {
  name_prefix          = "backend"
  iam_instance_profile = "${aws_iam_instance_profile.backend.name}"
  instance_type        = "t2.micro"
  image_id             = "${var.image_id}"
  security_groups      = ["${aws_security_group.ecs.id}"]

  user_data = <<EOF
#!/bin/bash
echo ECS_CLUSTER=${aws_ecs_cluster.backend.name} >> /etc/ecs/ecs.config
EOF
}

resource "aws_autoscaling_group" "backend" {
  availability_zones   = "${var.availability_zones}"
  name                 = "backend"
  min_size             = 1
  max_size             = 1
  launch_configuration = "${aws_launch_configuration.backend.name}"
}

resource "aws_ecr_repository" "backend" {
  name = "backend"
}

resource "aws_ecs_task_definition" "backend" {
  family       = "backend"
  network_mode = "bridge"

  container_definitions = <<EOF
[{
"name": "backend",
"image": "${aws_ecr_repository.backend.repository_url}:latest",
"privileged": false,
"disableNetworking": false,
"cpu": 0,
"memory": 990,
"memoryReservation": 384,
"essential": true,
"portMappings": [{"containerPort": 8080, "hostPort": 0}],
"logConfiguration": {
  "logDriver": "awslogs",
  "options": {
    "awslogs-group": "${aws_cloudwatch_log_group.backend.name}",
    "awslogs-region": "${var.region}",
    "awslogs-stream-prefix": "backend"
  }
},
"environment": [
  {
    "name": "CONSUMER_KEY",
    "value": "${var.consumer_key}"
  },
  {
    "name": "CONSUMER_SECRET",
    "value": "${var.consumer_secret}"
  },
  {
    "name": "CALLBACK_URL",
    "value": "${var.callback_url}"
  },
  {
    "name": "PORT",
    "value": "${var.port}"
  },
  {
    "name": "HOMEPAGE",
    "value": "${var.homepage}"
  },
  {
    "name": "HOST",
    "value": "${var.host}"
  },
  {
    "name": "PROTOCOL",
    "value": "${var.protocol}"
  }
]
}]
EOF
}

resource "aws_ecs_service" "backend" {
  name            = "backend"
  cluster         = "${aws_ecs_cluster.backend.id}"
  task_definition = "${aws_ecs_task_definition.backend.arn}"
  desired_count   = 1

  load_balancer {
    container_name   = "backend"
    container_port   = "${var.port}"
    target_group_arn = "${aws_lb_target_group.backend.arn}"
  }

  placement_strategy {
    type  = "spread"
    field = "host"
  }

  depends_on = ["aws_lb.backend"]
}

resource "aws_lb" "backend" {
  name            = "backend"
  security_groups = ["${aws_security_group.lb.id}"]
  subnets         = "${var.subnets}"

  enable_deletion_protection = false
}

resource "aws_lb_target_group" "backend" {
  name     = "backend"
  port     = 80
  protocol = "HTTP"
  vpc_id   = "${aws_default_vpc.default.id}"
}

resource "aws_lb_listener" "backend" {
  load_balancer_arn = "${aws_lb.backend.arn}"
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "${var.certificate_arn}"

  default_action {
    target_group_arn = "${aws_lb_target_group.backend.arn}"
    type             = "forward"
  }
}

resource "aws_cloudwatch_log_group" "backend" {
  name = "backend"
}

resource "aws_iam_role" "backend" {
  name = "backend"

  assume_role_policy = <<EOF
{
"Version": "2012-10-17",
"Statement": [
  {
    "Effect": "Allow",
    "Principal": {
      "Service": "ec2.amazonaws.com"
    },
    "Action": "sts:AssumeRole"
  }
]
}
EOF
}

resource "aws_iam_role_policy" "backend" {
  name = "backend"
  role = "${aws_iam_role.backend.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecs:CreateCluster",
        "ecs:DeregisterContainerInstance",
        "ecs:DiscoverPollEndpoint",
        "ecs:Poll",
        "ecs:RegisterContainerInstance",
        "ecs:StartTelemetrySession",
        "ecs:UpdateContainerInstancesState",
        "ecs:Submit*",
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "backend" {
  name = "backend"
  role = "${aws_iam_role.backend.name}"
}
