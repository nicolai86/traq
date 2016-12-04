variable "bucket_name" {
  default = "nicolai86-traq-data"
}

provider "aws" {}

resource "aws_s3_bucket" "traq_data" {
  tags {
    Application = "traq"
  }

  bucket = "${var.bucket_name}"
  acl    = "private"
}

resource "aws_iam_user" "traq" {
  name = "traq-rw"
}

resource "aws_iam_access_key" "rw" {
  user = "${aws_iam_user.traq.name}"
}

resource "aws_iam_user_policy" "traq_rw" {
  name = "traq-rw-access"
  user = "${aws_iam_user.traq.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Action": "s3:*",
      "Resource": [
        "${aws_s3_bucket.traq_data.arn}/*",
        "${aws_s3_bucket.traq_data.arn}"
      ]
    }
  ]
}
EOF
}

output "aws_access_key_id" {
  value = "${aws_iam_access_key.rw.id}"
}

output "aws_secret_access_key" {
  value = "${aws_iam_access_key.rw.secret}"
}
