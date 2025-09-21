-- Add name column to actions table
ALTER TABLE actions ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '';