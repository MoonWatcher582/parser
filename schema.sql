CREATE TABLE restaurant (
	id							INTEGER PRIMARY KEY AUTO_INCREMENT,
	name						TEXT NOT NULL,
	rank						INTEGER NOT NULL,
	restaurant_avg_stars	FLOAT NOT NULL
) DEFAULT CHARSET latin1 COLLATE latin1_general_cs;

CREATE TABLE reviewer (
	id								INTEGER PRIMARY KEY AUTO_INCREMENT,
	user_name					VARCHAR(255) UNIQUE NOT NULL,
	user_reviews				INTEGER NOT NULL,
	user_restaurant_reviews INTEGER NOT NULL,
	user_helpful_votes		INTEGER NOT NULL
) DEFAULT CHARSET latin1 COLLATE latin1_general_cs;

CREATE TABLE review (
	id					INTEGER PRIMARY KEY AUTO_INCREMENT,
	restaurant_id	INTEGER NOT NULL,
	user_id			INTEGER NOT NULL,
	review_stars	INTEGER NOT NULL,
	review_date		TEXT NOT NULL,
	FOREIGN KEY(restaurant_id) REFERENCES restaurant(id) ON DELETE CASCADE,
	FOREIGN KEY(user_id) REFERENCES reviewer(id) ON DELETE CASCADE
) DEFAULT CHARSET latin1 COLLATE latin1_general_cs;
