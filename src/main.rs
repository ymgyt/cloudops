use std::error::Error;

fn init_logger() {
    use tracing_subscriber::{filter, fmt::time};

    tracing_subscriber::fmt()
        .with_timer(time::UtcTime::rfc_3339())
        .with_ansi(true)
        .with_file(true)
        .with_line_number(true)
        .with_target(false)
        .with_env_filter(filter::EnvFilter::from_env(
            cloudops::envspec::LOG_DIRECTIVES,
        ))
        .init();
}

#[tokio::main]
async fn main() {
    init_logger();

    let app = cloudops::CloudOpsApp::parse();
    if let Err(app_err) = app.exec().await {
        tracing::error!("{}", app_err);

        let mut err: &dyn Error = &app_err;
        while let Some(source) = err.source() {
            tracing::error!("{}", source);
            err = source;
        }

        std::process::exit(app_err.exit_code());
    }
}
