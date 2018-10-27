# feelings

This is a simple AWS SAM-based example application for a blog post.

## Building and deploying

Use make target "deps" to install dependencies, "build" to build, "package" to package and upload application to S3 (NOTE: please update the bucket name in Makefile to point to an S3 bucket owned by you) and "deploy" to deploy using CloudFormation.
