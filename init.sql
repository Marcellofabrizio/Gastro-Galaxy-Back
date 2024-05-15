CREATE TABLE category (
  id SERIAL PRIMARY KEY,
  name TEXT
);

CREATE TABLE recipe (
  id SERIAL PRIMARY KEY,
  name TEXT,
  description TEXT,
  imageUrl TEXT,
  category_id INTEGER,
  CONSTRAINT fk_category
    FOREIGN KEY(category_id)
      REFERENCES category(id)
);

CREATE TABLE ingredient(
  id SERIAL PRIMARY KEY,
  name TEXT,
  amount TEXT,
  imageUrl TEXT,
  isAvailable BOOLEAN
);

CREATE TABLE ingredient_recipe (
  id SERIAL PRIMARY KEY,
  ingredient_id INTEGER,
  recipe_id INTEGER,
  CONSTRAINT fk_ingredient 
    FOREIGN KEY (ingredient_id) 
      REFERENCES ingredient(id),
  CONSTRAINT fk_recipe 
    FOREIGN KEY (recipe_id) 
      REFERENCES recipe(id)
);