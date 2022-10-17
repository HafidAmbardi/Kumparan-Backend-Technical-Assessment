# Kumparan Backend Technical Assessment

1. GET http://localhost/article

	Retrieves articles from the Database
	Result of GET will be sorted by latest first
	
			Query Parameters (optional)
	
			1. Query:
			URL: http://localhost/article?query=searchword
			
	
			2. Author:
			URL: http://localhost/article?author=authorname
	
		Response:
		[	
			{
			"id": "id",
			"author_id": "author_id",
			"author_name": "author_name",
			"title": "title",
			"body": "body",
			"created_at": "created_at"
			}
		]

2. POST: http://localhost/article

	Add articles to the Database

		Body:
		{
		  "author_name": "insertAuthorName",
		  "title": "insertArticleTitle",
		  "body": "insertArticleBody"
		}

		Response:
		{
		"id": "id",
		"author_id": "author_id",
		"author_name": "author_name",
		"title": "title",
		"body": "body",
		"created_at": "created_at"
		}
