use const_format::concatcp;

const PREFIX: &str = "CLOUDOPS";

/// Specify logging directives.
/// for more information see https://docs.rs/tracing-subscriber/latest/tracing_subscriber/filter/struct.EnvFilter.html#directives
pub const LOG_DIRECTIVES: &str = concatcp!(PREFIX, "_", "LOG");
