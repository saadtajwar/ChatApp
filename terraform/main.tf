data "aws_ecs_cluster" "chatapp-cluster" {
  cluster_name = "chatapp-cluster"
}


resource "aws_ecs_task_definition" "chatapp-task" {
  family             = "chatapp-task"
  task_role_arn      = aws_iam_role.chatapp.arn
  execution_role_arn = aws_iam_role.chatapp.arn
  network_mode       = "bridge"

  cpu    = 128
  memory = 256

  lifecycle {
    create_before_destroy = true
  }

  container_definitions = jsonencode([{

    cpu    = 128
    memory = 256

    name  = "golang-chat-app"
    image = format("%s:latest", aws_ecr_repository.chatapp.repository_url)

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group" = aws_cloudwatch_log_group.chatapp_logs.name
        "awslogs-region" : data.aws_region.current.name
        "awslogs-stream-prefix" = "chatapp"
      }
    }

    essential  = true

    environment = []
    secrets     = []
  }])

}

resource "aws_cloudwatch_log_group" "chatapp_logs" {
  name              = "/ecs/chatapp_logs"
  retention_in_days = 3
}

resource "aws_iam_role" "chatapp" {
  name               = "chatapp-task-${var.workspace}"
  assume_role_policy = data.aws_iam_policy_document.chatapp_assume_role.json
}

data "aws_iam_policy_document" "chatapp_assume_role" {
  statement {
    effect = "Allow"

    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy" "cloudwatch" {
  name   = "chatapp-cloudwatch"
  role   = aws_iam_role.chatapp.name
  policy = data.aws_iam_policy_document.cloudwatch.json
}

data "aws_iam_policy_document" "cloudwatch" {
  statement {
    effect = "Allow"

    actions = [
      "ecr:GetAuthorizationToken",
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_role_policy" "ecr" {
  name   = "chatapp-ecr"
  role   = aws_iam_role.chatapp.name
  policy = data.aws_iam_policy_document.ecr.json
}

data "aws_iam_policy_document" "ecr" {
  statement {
    effect = "Allow"

    actions = [
      "ecr:GetAuthorizationToken",
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage",
      "ecr:BatchCheckLayerAvailability",
    ]

    resources = [
        aws_ecr_repository.chatapp.arn
    ]
  }
}

resource "aws_iam_role_policy" "dynamo" {
  name   = "chatapp-dynamo"
  role   = aws_iam_role.chatapp.name
  policy = data.aws_iam_policy_document.dynamo.json
}

data "aws_iam_policy_document" "dynamo" {

  statement {
    sid    = "dynamo"
    effect = "Allow"

    actions = [
      "dynamodb:Query",
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:DeleteItem"
    ]

    resources = [
      aws_dynamodb_table.chatapp_users.arn
    ]
  }
}

resource "aws_dynamodb_table" "chatapp_users" {
  name         = "Users"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "Username"
    type = "S"
  }
}

resource "aws_ecr_repository" "chatapp" {
  name = "chatapp"
}

data "aws_region" "current" {}

data "aws_caller_identity" "current" {}

data "aws_subnets" "app_subnets" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.main.id]
  }
}

