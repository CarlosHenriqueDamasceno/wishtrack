ALTER TABLE contents
    ALTER COLUMN genres SET DATA TYPE VARCHAR(100) [] USING string_to_array(genres, '|');