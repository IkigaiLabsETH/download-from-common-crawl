resource "aws_s3_bucket" "bucket" {
  bucket = "${terraform.workspace}-common-crawl-query-${random_pet.dash.id}"

  tags = {
    environment = terraform.workspace
  }
}

resource "aws_s3_bucket_acl" "bucket_acl" {
  bucket = aws_s3_bucket.bucket.id
  acl    = "private"
}

resource "aws_lambda_permission" "lambda_s3_permission" {
  statement_id  = "AllowS3Invoke"
  action        = "lambda:InvokeFunction"
  function_name = module.s3update_function.function.function_name
  principal     = "s3.amazonaws.com"
  source_arn    = "arn:aws:s3:::${aws_s3_bucket.bucket.id}"
}

resource "aws_s3_bucket_notification" "aws_lambda_trigger" {
  bucket = aws_s3_bucket.bucket.id
  lambda_function {
    lambda_function_arn = module.s3update_function.function.arn
    events              = ["s3:ObjectCreated:*", "s3:ObjectRemoved:*"]
  }
  depends_on = [
    aws_lambda_permission.lambda_s3_permission
  ]
}