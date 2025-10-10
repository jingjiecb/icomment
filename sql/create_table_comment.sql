CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	article_url TEXT NOT NULL,
	parent_id INTEGER,
	nickname TEXT NOT NULL,
	email TEXT NOT NULL,
	content TEXT NOT NULL,
	status TEXT DEFAULT 'pending',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_status ON comments(status);
CREATE INDEX IF NOT EXISTS idx_article_url ON comments(article_url);
CREATE INDEX IF NOT EXISTS idx_email ON comments(email);