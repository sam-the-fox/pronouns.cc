#[derive(Debug)]
pub struct User {
    pub id: String,
    pub sid: String,
    pub username: String,
}

#[derive(Debug)]
pub struct Member {
    pub id: String,
    pub user_id: String,
    pub sid: String,
    pub name: String,
}
