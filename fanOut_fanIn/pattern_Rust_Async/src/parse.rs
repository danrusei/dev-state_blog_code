use serde::Deserialize;
use serde_json;
use url::Url;
use surf;

pub type Error = Box<(dyn std::error::Error + Send + Sync + 'static)>;

#[derive(Deserialize, Debug)]
pub struct Post {
    #[serde(rename(deserialize = "postId"))]
    post_id: u32,
    id: u32,
    name: String,
    email: String,
    body: String,
}

pub async fn parse(id: u32) -> Result<String, Error> {
    let body: String = get_body(id).await?;
    let posts: Vec<Post> = serde_json::from_str(&body)?;

    let mut longest_post = 0;
    let mut longest_post_id = 0;
    let mut longest_post_email = String::new();

    for post in posts {
        if post.body.len() > longest_post {
            longest_post = post.body.len();
            longest_post_id = post.post_id;
            longest_post_email = post.email;
        }
    }

    let the_result = format!(
        "{} {} {}",
        longest_post_id, longest_post_email, longest_post
    );

    Ok(the_result)
}

async fn get_body(id: u32) -> Result<String, Error> {
    const BASE: &'static str = "https://test-apps-257216.ew.r.appspot.com/";
    let base = Url::parse(BASE).expect("hardcoded URL is known to be valid");
    let site = base.join(&format!("/comments/{}", id))?;

    let body = surf::get(site).recv_string().await?;

    Ok(body)
}
