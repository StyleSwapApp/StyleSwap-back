package model

type ArticleResponse struct {
	UserPseudo  string  `json:"user_pseudo"`
	ArticleName string  `json:"article_name"`
	ArticlePrice int    `json:"article_price"`
	ArticleDescription string `json:"article_description"`
	ArticleImage string `json:"article_image"`
}

