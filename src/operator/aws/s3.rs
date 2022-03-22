use crate::errors::AppError;
use aws_sdk_s3::model::Object;
use aws_sdk_s3::output::PutObjectOutput;
use aws_sdk_s3::types::ByteStream;
use aws_sdk_s3::{model::Bucket, Client};
use camino::Utf8Path;

pub struct S3Operator {
    client: Client,
}

impl S3Operator {
    pub fn new(client: Client) -> Self {
        Self { client }
    }

    /// Fetch s3 bucket models from s3 service.
    pub async fn list_buckets(&self) -> Result<Vec<Bucket>, AppError> {
        Ok(self
            .client
            .list_buckets()
            .send()
            .await?
            .buckets
            .unwrap_or_default())
    }

    pub async fn list_objects(&self, bucket: impl Into<String>) -> Result<Vec<Object>, AppError> {
        let bucket = bucket.into();

        let mut objects = Vec::new();
        let mut next_token = None;
        loop {
            let list_output = self
                .client
                .list_objects_v2()
                .bucket(bucket.clone())
                .set_continuation_token(next_token.take())
                .send()
                .await?;

            objects.extend(list_output.contents.unwrap_or_default());

            match list_output.continuation_token {
                Some(token) if !token.is_empty() => next_token = Some(token),
                Some(_) | None => break,
            }
        }

        Ok(objects)
    }

    pub async fn put(
        &self,
        src: &Utf8Path,
        bucket: impl Into<String>,
        object_key: impl Into<String>,
    ) -> Result<PutObjectOutput, AppError> {
        let src = tokio::fs::File::open(src.as_std_path())
            .await
            .map_err(|err| AppError::OpenFile {
                err,
                path: src.to_path_buf(),
            })?;

        let src = ByteStream::from_file(src)
            .await
            .map_err(|err| AppError::AwsError(Box::new(err)))?;

        let put_output = self
            .client
            .put_object()
            .bucket(bucket)
            .key(object_key)
            .body(src)
            .send()
            .await?;

        Ok(put_output)
    }
}
