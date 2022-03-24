/// Print package version with 'v' prefix
fn main() {
    println!(
        "v{}",
        std::env::var("CARGO_PKG_VERSION")
            .expect("environment variable CARGO_PKG_VERSION expected")
    );
}
