use std::env;
use std::net::SocketAddr;
use std::sync::Arc;

use axum::{
    extract::{Path, State},
    http::StatusCode,
    routing::get,
    Json, Router,
};
use sqlx::{postgres::PgPoolOptions, Pool, Postgres};
use tracing::{error, info};
use tracing_subscriber;
use user::ExportUser;

use crate::user::DBUser;

pub mod user;

struct AppState {
    pool: Pool<Postgres>,
}

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    info!("Loading .env");

    dotenvy::dotenv().unwrap();

    info!("Connecting to database");

    let db_dsn = env::var("DATABASE_URL").expect("DATABASE_URL is not set or invalid!");

    let pool = PgPoolOptions::new()
        .max_connections(10)
        .connect(db_dsn.as_str())
        .await
        .expect("Could not connect do database!");

    info!("Initializing exporter server");

    let shared_state = Arc::new(AppState { pool });

    let app = Router::new()
        .route("/start/:id", get(handler))
        .with_state(shared_state);

    let port = match env::var("EXPORTER_PORT") {
        Ok(val) => val
            .parse::<u16>()
            .expect("EXPORTER_PORT is not a valid integer!"),
        Err(_) => {
            error!("EXPORTER_PORT is not set, falling back to port 8081");
            8081
        }
    };

    info!(port, "Starting server");

    let addr = SocketAddr::from(([127, 0, 0, 1], port));
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn handler(
    State(state): State<Arc<AppState>>,
    Path(user_id): Path<String>,
) -> Result<(StatusCode, Json<ExportUser>), (StatusCode, String)> {
    info!(user_id = user_id, "Received thing :)");

    let user = match sqlx::query_as!(
        DBUser,
        "SELECT id, username, display_name FROM users WHERE id = $1",
        user_id
    )
    .fetch_one(&state.pool)
    .await
    {
        Ok(val) => val,
        Err(why) => return Err((StatusCode::INTERNAL_SERVER_ERROR, format!("{}", why))),
    };

    info!(id = user.id, username = user.username, "Got user!");

    Ok((StatusCode::ACCEPTED, Json(user.to_export())))
}
