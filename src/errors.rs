use aws_smithy_http::result::SdkError;
use camino::Utf8PathBuf;
use std::fmt;
use std::fmt::Formatter;
use std::io;

/// Error contracting in CLI execution.
#[derive(Debug)]
pub enum AppError {
    /// Invalid endpoint provided.
    InvalidEndpoint(http::uri::InvalidUri),
    AwsError(Box<dyn std::error::Error + Send + Sync + 'static>),
    OpenFile {
        err: io::Error,
        path: Utf8PathBuf,
    },
}

impl fmt::Display for AppError {
    fn fmt(&self, f: &mut Formatter<'_>) -> fmt::Result {
        use AppError::*;
        match self {
            InvalidEndpoint(err) => write!(f, "invalid endpoint: {}", err),
            AwsError(err) => write!(f, "aws error: {}", err),
            OpenFile { err, path } => write!(f, "open file {}: {}", path, err),
        }
    }
}

impl std::error::Error for AppError {}

impl AppError {
    pub fn exit_code(&self) -> i32 {
        1
    }
}

impl<E> From<aws_smithy_http::result::SdkError<E>> for AppError
where
    E: std::error::Error + Send + Sync + 'static,
{
    fn from(err: SdkError<E>) -> Self {
        AppError::AwsError(Box::new(err))
    }
}
