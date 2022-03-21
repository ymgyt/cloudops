use crate::errors::AppError;
use clap::Parser;

#[derive(Parser, Debug)]
#[clap(version)]
pub struct CloudOpsApp {}

impl CloudOpsApp {
    pub fn parse() -> Self {
        clap::Parser::parse()
    }

    /// Entry point of execution.
    pub fn exec(self) -> Result<(), AppError> {
        Ok(())
    }
}
