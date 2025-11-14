-- Sample seed data for Auto LMK platform
-- Run this after migrations: psql -U autolmk -d autolmk < scripts/seed.sql

-- Create sample tenants
INSERT INTO tenants (domain, name, whatsapp_number, pairing_status, status) VALUES
('showroom-jaya.localhost', 'Showroom Jaya Motor', '6281234567890', 'unpaired', 'active'),
('mobilindo.localhost', 'Mobilindo Premium', '6289876543210', 'unpaired', 'active');

-- Create sample users (password: password123)
-- Hash generated with: bcrypt.GenerateFromPassword([]byte("password123"), 12)
INSERT INTO users (tenant_id, email, password_hash, role) VALUES
(1, 'admin@showroom-jaya.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5yvWEZjxo8Wlm', 'owner'),
(2, 'admin@mobilindo.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5yvWEZjxo8Wlm', 'owner');

-- Create sample sales for tenant 1
INSERT INTO sales (tenant_id, phone_number, name, status) VALUES
(1, '628123456001', 'Budi Santoso', 'active'),
(1, '628123456002', 'Ani Wijaya', 'active'),
(2, '628987654001', 'Dedi Kurniawan', 'active');

-- Create sample cars for tenant 1 (Showroom Jaya)
INSERT INTO cars (tenant_id, brand, model, year, price, mileage, transmission, fuel_type, engine_cc, seats, color, description, status, is_featured) VALUES
(1, 'Toyota', 'Avanza', 2023, 250000000, 5000, 'automatic', 'bensin', 1500, 7, 'Silver', 'Toyota Avanza 2023 kondisi sangat baik, service record lengkap di Auto2000. Mobil keluarga irit dan nyaman.', 'available', true),
(1, 'Honda', 'CR-V', 2022, 550000000, 15000, 'automatic', 'bensin', 2000, 7, 'Black', 'Honda CR-V Turbo Prestige 2022, full original, pajak panjang. SUV premium dengan fitur lengkap.', 'available', true),
(1, 'Mitsubishi', 'Xpander', 2023, 280000000, 8000, 'automatic', 'bensin', 1500, 7, 'White', 'Mitsubishi Xpander Ultimate 2023, seperti baru. Cocok untuk keluarga modern.', 'available', false),
(1, 'Toyota', 'Fortuner', 2021, 520000000, 25000, 'automatic', 'diesel', 2400, 7, 'Grey', 'Toyota Fortuner VRZ 2021, terawat istimewa. Diesel irit untuk perjalanan jauh.', 'available', false),
(1, 'Daihatsu', 'Terios', 2020, 210000000, 30000, 'manual', 'bensin', 1500, 7, 'Red', 'Daihatsu Terios X 2020, manual transmission. Harga nego, siap pakai.', 'available', false);

-- Create sample cars for tenant 2 (Mobilindo)
INSERT INTO cars (tenant_id, brand, model, year, price, mileage, transmission, fuel_type, engine_cc, seats, color, description, status, is_featured) VALUES
(2, 'BMW', 'X5', 2022, 1200000000, 10000, 'automatic', 'bensin', 3000, 5, 'Black', 'BMW X5 xDrive40i 2022, luxury SUV dengan performa tinggi. Full spec, sunroof, dan premium audio system.', 'available', true),
(2, 'Mercedes-Benz', 'C-Class', 2023, 950000000, 3000, 'automatic', 'bensin', 2000, 5, 'White', 'Mercedes-Benz C200 AMG 2023, seperti baru. Elegant sedan dengan teknologi canggih.', 'available', true),
(2, 'Audi', 'A4', 2021, 750000000, 18000, 'automatic', 'bensin', 2000, 5, 'Silver', 'Audi A4 TFSI 2021, kondisi mint. German engineering at its best.', 'available', false);

-- Create sample car specs (EAV pattern for specific attributes)
INSERT INTO car_specs (car_id, key, value) VALUES
(1, 'Airbag', 'Dual front airbag'),
(1, 'Audio System', '2DIN Touchscreen'),
(1, 'Velg', 'R15 original Toyota'),
(2, 'Sunroof', 'Panoramic sunroof'),
(2, 'Cruise Control', 'Adaptive cruise control'),
(2, 'Parking Sensor', 'Front & rear with camera'),
(6, 'Head-up Display', 'Yes'),
(6, 'Sound System', 'Harman Kardon premium'),
(7, 'MBUX System', 'Latest generation'),
(7, 'Ambient Lighting', '64 colors');

-- Create sample conversation and messages (for testing)
INSERT INTO conversations (tenant_id, sender_phone, is_sales) VALUES
(1, '6281111111111', false), -- Customer conversation
(1, '628123456001', true);    -- Sales conversation

INSERT INTO messages (conversation_id, sender_phone, message_text, direction) VALUES
(1, '6281111111111', 'Ada mobil Toyota budget 300 juta?', 'inbound'),
(1, 'BOT', 'Tentu! Kami memiliki beberapa pilihan Toyota dalam budget Anda. Ada Toyota Avanza 2023 (Rp 250 juta) dan Toyota Fortuner 2021 (Rp 520 juta). Yang Avanza sangat cocok untuk budget Anda. Mau lihat detailnya?', 'outbound'),
(1, '6281111111111', 'Yang Avanza boleh, kirim fotonya dong', 'inbound'),
(2, '628123456001', '/list', 'inbound'),
(2, 'BOT', 'Berikut daftar inventory Showroom Jaya:\n1. Toyota Avanza 2023 - Rp 250 juta\n2. Honda CR-V 2022 - Rp 550 juta\n3. Mitsubishi Xpander 2023 - Rp 280 juta\n...', 'outbound');

-- Create sample leads
INSERT INTO leads (tenant_id, phone_number, name, interested_car_id, conversation_id, status) VALUES
(1, '6281111111111', 'Pak Andi', 1, 1, 'new'),
(1, '6282222222222', 'Ibu Sari', 2, NULL, 'contacted'),
(2, '6283333333333', 'Mr. John', 6, NULL, 'new');

-- Show created data
SELECT 'Tenants created:' as info;
SELECT id, domain, name FROM tenants;

SELECT 'Cars created:' as info;
SELECT tenant_id, brand, model, year, price FROM cars ORDER BY tenant_id, id;

SELECT 'Sales registered:' as info;
SELECT tenant_id, name, phone_number FROM sales ORDER BY tenant_id;

SELECT 'Leads generated:' as info;
SELECT tenant_id, phone_number, name, status FROM leads ORDER BY tenant_id;
