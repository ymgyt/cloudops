use std::fmt;
use std::fmt::Formatter;

/// Error contracting in CLI execution.
#[derive(Debug)]
pub enum AppError {}

impl fmt::Display for AppError {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        todo!()
    }
}

impl std::error::Error for AppError {}

impl AppError {
    pub fn exit_code(&self) -> i32 {
        todo!()
    }
}
