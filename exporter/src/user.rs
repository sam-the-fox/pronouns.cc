use serde::Serialize;

pub struct DBUser {
    pub id: String,
    pub username: String,
    pub display_name: Option<String>,
}

impl DBUser {
    pub fn to_export(self) -> ExportUser {
        ExportUser {
            id: self.id,
            username: self.username,
            display_name: self.display_name,
        }
    }
}

#[derive(Serialize, Debug)]
pub struct ExportUser {
    pub id: String,
    pub username: String,
    pub display_name: Option<String>,
}
