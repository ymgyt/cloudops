use crate::errors::AppError;
use crate::operator::aws::s3::S3Operator;
use clap::{Args, Subcommand};
use tabled::style::Style;

#[derive(Args, Debug)]
pub struct BucketCommand {
    #[clap(subcommand)]
    pub command: BucketSubcommand,
}

#[derive(Subcommand, Debug)]
pub enum BucketSubcommand {
    /// List s3 buckets.
    #[clap(visible_alias = "ls")]
    List {},
}

impl BucketCommand {
    pub async fn exec(self, operator: S3Operator) -> Result<(), AppError> {
        use BucketSubcommand::*;
        match self.command {
            List { .. } => self.exec_list_buckets(operator).await,
        }
    }

    async fn exec_list_buckets(&self, operator: S3Operator) -> Result<(), AppError> {
        let buckets = operator.list_buckets().await?;

        let mut table = tabled::builder::Builder::default().set_header(["bucket", "creation_date"]);

        for bucket in buckets {
            table = table.add_row([
                bucket.name.unwrap_or_default(),
                bucket
                    .creation_date
                    .and_then(|date| date.fmt(aws_smithy_types::date_time::Format::DateTime).ok())
                    .unwrap_or_default(),
            ]);
        }

        println!("{}", table.build().with(Style::blank()));

        Ok(())
    }
}
