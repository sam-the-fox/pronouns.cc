use pronouns::models::{member::Member, user::User};
use rocket::{response::Redirect, Config, State};
use simple_logger;
use sqlx::{postgres::PgPoolOptions, Pool, Postgres};
use std::{env, path::PathBuf};
#[macro_use]
extern crate rocket;

struct AppState {
    pool: Pool<Postgres>,
    base_url: String,
}

#[get("/<path..>")]
async fn redirect(state: &State<AppState>, path: PathBuf) -> Redirect {
    let sid = path.to_str().unwrap_or("");

    match sid.len() {
        5 => sqlx::query_as!(
            User,
            "SELECT id, sid, username FROM users WHERE sid = $1",
            sid
        )
        .fetch_one(&state.pool)
        .await
        .map_or_else(
            |_| Redirect::temporary(state.base_url.clone()),
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
                Err(_) => return Redirect::temporary(state.base_url.clone()),
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
                |_| Redirect::temporary(state.base_url.clone()),
                |u| {
                    Redirect::temporary(format!(
                        "{}/@{}/{}",
                        state.base_url, u.username, member.name
                    ))
                },
            )
        }
        _ => Redirect::temporary(state.base_url.clone()),
    }
}

#[launch]
async fn rocket() -> _ {
    simple_logger::SimpleLogger::new().env().init().unwrap();

    info!("Loading .env");

    dotenvy::dotenv().unwrap();

    info!("Connecting to database");

    let db_dsn = env::var("DATABASE_URL").expect("DATABASE_URL is not set or invalid!");
    let base_url = env::var("BASE_URL").expect("BASE_URL is not set or invalid!");
    let port = match env::var("PRNS_PORT") {
        Ok(val) => val
            .parse::<u16>()
            .expect("PRNS_PORT is not a valid integer!"),
        Err(_) => {
            error!("PRNS_PORT is not set, falling back to port 8087");
            8087
        }
    };

    let pool = PgPoolOptions::new()
        .max_connections(10)
        .connect(db_dsn.as_str())
        .await
        .expect("Could not connect do database!");

    info!("Initializing prns.cc server");

    let config = Config {
        port,
        ..Config::default()
    };

    rocket::build()
        .configure(config)
        .manage(AppState { pool, base_url })
        .mount("/", routes![redirect])
}
