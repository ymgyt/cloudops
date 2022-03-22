mod aws;

use crate::errors::AppError;
use clap::{Parser, Subcommand};

#[derive(Parser, Debug)]
#[clap(version, propagate_version = true)]
pub struct CloudOpsApp {
    #[clap(subcommand)]
    command: Command,
}

#[derive(Subcommand, Debug)]
enum Command {
    /// S3 operations.
    S3(aws::s3::S3Command),
}

impl CloudOpsApp {
    pub fn parse() -> Self {
        clap::Parser::parse()
    }

    /// Entry point of execution.
    pub async fn exec(self) -> Result<(), AppError> {
        use Command::*;
        match self.command {
            S3(cmd) => cmd.exec().await,
        }
    }
}
