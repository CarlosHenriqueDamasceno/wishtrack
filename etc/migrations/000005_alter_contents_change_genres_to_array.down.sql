ALTER TABLE contents
    ALTER COLUMN genres SET DATA TYPE text USING array_to_string(genres, '|'),
    ALTER COLUMN genres SET DEFAULT '',
    ALTER COLUMN genres DROP NOT NULL;