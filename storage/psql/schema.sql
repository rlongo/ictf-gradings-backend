CREATE TABLE IF NOT EXISTS belt_tests (
  id SERIAL PRIMARY KEY,
  test_name VARCHAR(32) NOT NULL,
  test_date INT NOT NULL,
  dojang VARCHAR(32) NOT NULL,
  admins TEXT NOT NULL
);


