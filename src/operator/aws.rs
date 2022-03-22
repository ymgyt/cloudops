use crate::errors::AppError;
use aws_smithy_http::endpoint::Endpoint;
use aws_types::region::Region;
use std::str::FromStr;

pub(crate) mod s3;

pub struct AwsClientBuilder {}

impl AwsClientBuilder {
    pub fn new() -> Self {
        Self {}
    }

    pub async fn s3_client(&self, context: AwsContext) -> Result<aws_sdk_s3::Client, AppError> {
        let mut builder = aws_sdk_s3::Config::builder()
            .region(context.region())
            .credentials_provider(context.credential_provider());

        if let Some(endpoint) = context.endpoint()? {
            builder = builder.endpoint_resolver(endpoint);
        }

        Ok(aws_sdk_s3::Client::from_conf(builder.build()))
    }
}

impl Default for AwsClientBuilder {
    fn default() -> Self {
        AwsClientBuilder::new()
    }
}

pub struct AwsContext {
    pub region: String,
    pub access_key_id: String,
    pub secret_access_key: String,
    pub endpoint: Option<String>,
}

impl AwsContext {
    fn endpoint(&self) -> Result<Option<Endpoint>, AppError> {
        match self.endpoint {
            Some(ref endpoint) => Ok(Some(Endpoint::immutable(
                http::uri::Uri::from_str(endpoint).map_err(AppError::InvalidEndpoint)?,
            ))),
            None => Ok(None),
        }
    }

    fn credential_provider(&self) -> aws_types::Credentials {
        aws_types::Credentials::new(
            self.access_key_id.clone(),
            self.secret_access_key.clone(),
            None,
            None,
            "cloudops",
        )
    }

    fn region(&self) -> Region {
        Region::new(self.region.clone())
    }
}
