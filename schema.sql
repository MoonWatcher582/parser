CREATE TABLE restaurant (
	id							INTEGER PRIMARY KEY AUTOINCREMENT,
	/* restaurant name is NOT unique */
	restaurant				TEXT NOT NULL,
	rank						INTEGER NOT NULL,
	restaurant_avg_stars NUMBER NOT NULL
);

CREATE TABLE reviewer (
	id								INTEGER PRIMARY KEY AUTOINCREMENT,
	user_name					TEXT UNIQUE NOT NULL,
	user_reviews				INTEGER NOT NULL,
	user_restaurant_reviews INTEGER NOT NULL,
	user_helpful_votes		INTEGER NOT NULL
);

CREATE TABLE review (
	id					INTEGER PRIMARY KEY AUTOINCREMENT,
	restaurant_id	INTEGER NOT NULL,
	user_id			INTEGER NOT NULL,
	review_stars	INTEGER NOT NULL,
	review_date		TEXT NOT NULL,
	FOREIGN KEY(restaurant_id) REFERENCES restaurant(id),
	FOREIGN KEY(user_id) REFERENCES reviewer(id)
);
