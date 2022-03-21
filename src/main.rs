use std::error::Error;

fn init_logger() {
    use tracing_subscriber::{filter, fmt::time};

    tracing_subscriber::fmt()
        .with_timer(time::UtcTime::rfc_3339())
        .with_ansi(true)
        .with_file(true)
        .with_line_number(true)
        .with_env_filter(filter::EnvFilter::from_env(
            cloudops::envspec::LOG_DIRECTIVES,
        ))
        .init();
}

#[tokio::main]
async fn main() {
    init_logger();

    // handle signal
    let app = cloudops::CloudOpsApp::parse();
    if let Err(err) = app.exec() {
        tracing::error!("{}", err);

        let mut source = err.source();
        loop {
            match source {
                Some(err) => {
                    tracing::error!("{}", err);
                    source = Some(err);
                }
                None => break,
            }
        }

        std::process::exit(err.exit_code());
    }
}
