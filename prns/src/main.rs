use pronouns::models::{member::Member, user::User};
use rocket::{response::Redirect, State};
use sqlx::{postgres::PgPoolOptions, Pool, Postgres};
use std::env;
use tracing::info;
use tracing_subscriber;
#[macro_use]
extern crate rocket;

struct AppState {
    pool: Pool<Postgres>,
    base_url: String,
}

#[get("/<sid>")]
async fn redirect(state: &State<AppState>, sid: &str) -> Redirect {
    match sid.len() {
        5 => sqlx::query_as!(
            User,
            "SELECT id, sid, username FROM users WHERE sid = $1",
            sid
        )
        .fetch_one(&state.pool)
        .await
        .map_or_else(
            |_| Redirect::temporary(format!("{}", state.base_url)),
            |u| Redirect::temporary(format!("{}/@{}", state.base_url, u.username)),
        ),
        6 => {
            let member = match sqlx::query_as!(
                Member,
                "SELECT id, user_id, sid, name FROM members WHERE sid = $1",
                sid
            )
            .fetch_one(&state.pool)
            .await
            {
                Err(_) => return Redirect::temporary(format!("{}", state.base_url)),
                Ok(val) => val,
            };

            sqlx::query_as!(
                User,
                "SELECT id, sid, username FROM users WHERE id = $1",
                member.user_id
            )
            .fetch_one(&state.pool)
            .await
            .map_or_else(
                |_| Redirect::temporary(format!("{}", state.base_url)),
                |u| {
                    Redirect::temporary(format!(
                        "{}/@{}/{}",
                        state.base_url, u.username, member.name
                    ))
                },
            )
        }
        _ => Redirect::temporary(format!("{}", state.base_url)),
    }
}

#[launch]
async fn rocket() -> _ {
    tracing_subscriber::fmt::init();

    info!("Loading .env");

    dotenvy::dotenv().unwrap();

    info!("Connecting to database");

    let db_dsn = env::var("DATABASE_URL").expect("DATABASE_URL is not set or invalid!");
    let base_url = env::var("BASE_URL").expect("BASE_URL is not set or invalid!");

    let pool = PgPoolOptions::new()
        .max_connections(10)
        .connect(db_dsn.as_str())
        .await
        .expect("Could not connect do database!");

    info!("Initializing prns.cc server");

    rocket::build()
        .manage(AppState { pool, base_url })
        .mount("/", routes![redirect])
}
