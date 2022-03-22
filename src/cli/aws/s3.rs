mod bucket;

use crate::cli::aws::AwsOptions;
use crate::errors::AppError;
use crate::operator::aws::s3::S3Operator;
use crate::operator::aws::AwsClientBuilder;
use camino::{Utf8Path, Utf8PathBuf};
use clap::{Args, Subcommand};

#[derive(Args, Debug)]
#[clap(subcommand_required = true)]
pub struct S3Command {
    #[clap(flatten)]
    pub aws: AwsOptions,

    #[clap(subcommand)]
    pub command: S3Subcommand,
}

#[derive(Subcommand, Debug)]
pub enum S3Subcommand {
    /// S3 buckets operations.
    #[clap(visible_alias = "bkt")]
    Bucket(bucket::BucketCommand),

    /// Put file onto s3 bucket.
    #[clap(next_help_heading = "PUT_OPTIONS")]
    Put {
        /// Src file path to put.
        #[clap(long, value_name = "FILE_PATH")]
        src: Utf8PathBuf,

        /// Dest bucket.
        #[clap(long)]
        bucket: String,

        /// Object key.
        #[clap(long, alias = "object-key")]
        key: String,
    },

    /// List bucket objects.
    #[clap(visible_alias = "ls", next_help_heading = "LIST_OPTIONS")]
    List {
        /// Target bucket.
        #[clap(long)]
        bucket: String,
    },
}

impl S3Command {
    pub async fn exec(self) -> Result<(), AppError> {
        let s3_client = AwsClientBuilder::new()
            .s3_client(self.aws.into_context())
            .await?;
        let operator = S3Operator::new(s3_client);

        use S3Subcommand::*;
        match self.command {
            Bucket(bucket) => bucket.exec(operator).await,
            Put { src, bucket, key } => {
                S3Command::exec_put(operator, src.as_path(), bucket, key).await
            }
            List { bucket } => S3Command::exec_list(operator, bucket).await,
        }
    }

    pub async fn exec_put(
        operator: S3Operator,
        src: &Utf8Path,
        bucket: impl Into<String>,
        object_key: impl Into<String>,
    ) -> Result<(), AppError> {
        let bucket = bucket.into();
        let object_key = object_key.into();

        let output = operator
            .put(src, bucket.clone(), object_key.clone())
            .await?;

        tracing::debug!("{:?}", output);

        println!(
            "put s3://{}/{} {}",
            bucket,
            object_key,
            output.e_tag.unwrap_or_default().replace('"', "")
        );

        Ok(())
    }

    pub async fn exec_list(
        operator: S3Operator,
        bucket: impl Into<String>,
    ) -> Result<(), AppError> {
        let objects = operator.list_objects(bucket).await?;

        let mut table = tabled::builder::Builder::default().set_header([
            "object",
            "last_modified",
            "e_tag",
            "size",
            "storage_class",
        ]);

        for object in objects {
            table = table.add_row([
                object.key.unwrap_or_default(),
                object
                    .last_modified
                    .and_then(|date| date.fmt(aws_smithy_types::date_time::Format::DateTime).ok())
                    .unwrap_or_default(),
                object.e_tag.unwrap_or_default().replace('"', ""),
                object.size.to_string(),
                object
                    .storage_class
                    .map(|sc| sc.as_str().to_owned())
                    .unwrap_or_default(),
            ]);
        }

        println!("{}", table.build().with(tabled::style::Style::blank()));

        Ok(())
    }
}
