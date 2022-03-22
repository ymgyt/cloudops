pub(super) mod s3;

use crate::operator::aws::AwsContext;
use clap::Args;

/// Common options to request aws endpoints
#[derive(Args, Debug)]
#[clap(next_help_heading = "AWS")]
pub struct AwsOptions {
    // Should be enum.
    #[clap(long, env = "AWS_REGION")]
    pub region: String,

    #[clap(long)]
    pub endpoint: Option<String>,

    #[clap(flatten)]
    pub credentials: AwsCredentials,
}

// Should impl Debug for mask ?
/// AWS Credentials.
#[derive(Args, Debug)]
pub struct AwsCredentials {
    #[clap(long, env = "AWS_ACCESS_KEY_ID", hide_env_values = true)]
    pub access_key_id: String,

    #[clap(long, env = "AWS_SECRET_ACCESS_KEY", hide_env_values = true)]
    pub secret_access_key: String,
}

impl AwsOptions {
    pub fn into_context(self) -> AwsContext {
        AwsContext {
            region: self.region,
            secret_access_key: self.credentials.secret_access_key,
            access_key_id: self.credentials.access_key_id,
            endpoint: self.endpoint,
        }
    }
}
