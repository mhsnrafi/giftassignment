-- Create Employees Table
CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    interests VARCHAR(255)[],
    gift VARCHAR(255),
    returned BOOLEAN DEFAULT false
    );

-- Create Gifts Table
CREATE TABLE IF NOT EXISTS gifts (
     id SERIAL PRIMARY KEY,
     name VARCHAR(255) NOT NULL,
    categories VARCHAR(255)[],
    istaken BOOLEAN DEFAULT false,
    returned BOOLEAN DEFAULT false
    );

-- Insert Employees data
INSERT INTO employees (name, interests)
VALUES
    ('severin', ARRAY['high end audio', 'photography', 'sitting comfortably']),
    ('oliver', ARRAY['pets', 'scifi', 'music making']),
    ('henry', ARRAY['scifi', 'music', 'drinking', 'motorcross']),
    ('karsten', ARRAY['family', 'reading']),
    ('mehmet', ARRAY['eating', 'sleeping', 'power lifting']),
    ('jan', ARRAY['triathlons', 'techno', 'family']),
    ('alex', ARRAY['drinking', 'football', 'camping']),
    ('sarah', ARRAY['meditation', 'crossfit', 'yoga']),
    ('daniel', ARRAY['handball', 'pets']),
    ('werner', ARRAY['techno', 'music making']),
    ('robert', ARRAY['photography', 'eating']),
    ('daniel', ARRAY['pets']),
    ('jonathan', ARRAY['football', 'drinking', 'cocktails', 'running']),
    ('maxim', ARRAY['scuba', 'motorcross', 'warhammer']),
    ('michel', ARRAY['motorcross', 'family', 'warhammer']),
    ('marc', ARRAY['skydiving', 'reading', 'family']);

-- Insert Gifts data
INSERT INTO gifts (name, categories)
VALUES
    ('netflix card', ARRAY['music', 'film', 'pets', 'scifi']),
    ('gold plated water bottle', ARRAY['power lifting', 'triathlons', 'football', 'crossfit', 'handball', 'running']),
    ('second hand earplugs', ARRAY['techno', 'music']),
    ('6 pack of beer', ARRAY['drinking']),
    ('peterborough united football shirt', ARRAY['football']),
    ('bonsai tree', ARRAY['botany']),
    ('dog treats', ARRAY['pets']),
    ('running shoes', ARRAY['running', 'triathlons']),
    ('spotify voucher', ARRAY['music']),
    ('goggles', ARRAY['scuba', 'skydiving']),
    ('recipe book', ARRAY['eating']),
    ('family vacation to disneyland', ARRAY['family']),
    ('bottle of wine', ARRAY['drinking']),
    ('book of german phrases', ARRAY['german', 'reading']),
    ('cocktail mixer', ARRAY['cocktails', 'drinking']),
    ('dog halloween costume', ARRAY['pets']),
    ('yoga ball', ARRAY['yoga', 'sitting comfortably']),
    ('book on photographing nature', ARRAY['photography']);