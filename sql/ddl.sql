-- Create table locations
CREATE TABLE IF NOT EXISTS locations (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    address    TEXT,
    latitude   DECIMAL(10, 7),
    longitude  DECIMAL(10, 7),
    is_open    BOOLEAN      NOT NULL DEFAULT true,
    created_at TIMESTAMP    DEFAULT NOW()
);

-- Create table users
CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255)        NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    phone         VARCHAR(20),
    password_hash VARCHAR(255)        NOT NULL,
    role          VARCHAR(20)         NOT NULL DEFAULT 'user'
                  CHECK (role IN ('user', 'admin', 'master_admin')),
    location_id   INT REFERENCES locations(id),     
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW()
);

-- Create table waste_category
CREATE TABLE IF NOT EXISTS waste_category (
    id                 SERIAL PRIMARY KEY,
    name               VARCHAR(100) NOT NULL,
    parent_category_id INT REFERENCES waste_category(id),
    unit               VARCHAR(10)  NOT NULL DEFAULT 'kg'
                       CHECK (unit IN ('kg', 'liter')),
    is_active          BOOLEAN      NOT NULL DEFAULT true
);

-- Create table waste_price
CREATE TABLE IF NOT EXISTS waste_price (
    id              SERIAL PRIMARY KEY,
    category_id     INT            NOT NULL REFERENCES waste_category(id),
    price_per_unit  NUMERIC(12, 2) NOT NULL,
    effective_from  TIMESTAMP      NOT NULL DEFAULT NOW(),
    effective_until TIMESTAMP,                       
    set_by          INT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW()
);

-- Create table queue 
CREATE TABLE IF NOT EXISTS queue (
    id             SERIAL PRIMARY KEY,
    user_id        INT         NOT NULL REFERENCES users(id),
    location_id    INT         NOT NULL REFERENCES locations(id),
    queue_number   VARCHAR(20) NOT NULL,
    preferred_time TIMESTAMP,
    status         VARCHAR(20) NOT NULL DEFAULT 'waiting'
                   CHECK (status IN ('waiting', 'verified', 'completed', 'expired')),
    verified_by    INT REFERENCES users(id),
    created_at     TIMESTAMP DEFAULT NOW()
);

-- Create table waste_transaction
CREATE TABLE IF NOT EXISTS waste_transaction (
    id           SERIAL PRIMARY KEY,
    queue_id     INT REFERENCES queue(id),
    user_id      INT            NOT NULL REFERENCES users(id),
    location_id  INT            NOT NULL REFERENCES locations(id),
    admin_id     INT            NOT NULL REFERENCES users(id),
    total_rupiah NUMERIC(12, 2) NOT NULL DEFAULT 0,
    status       VARCHAR(20)    NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending', 'completed', 'cancelled')),
    created_at   TIMESTAMP DEFAULT NOW()
);

-- Create table waste_transaction_detail
CREATE TABLE IF NOT EXISTS waste_transaction_detail (
    id                      SERIAL PRIMARY KEY,
    transaction_id          INT            NOT NULL REFERENCES waste_transaction(id),
    category_id             INT            NOT NULL REFERENCES waste_category(id),
    weight                  NUMERIC(10, 2) NOT NULL,
    price_per_unit_snapshot NUMERIC(12, 2) NOT NULL,
    subtotal                NUMERIC(12, 2) NOT NULL
);

-- Create table wallet
CREATE TABLE IF NOT EXISTS wallet (
    id         SERIAL PRIMARY KEY,
    user_id    INT UNIQUE     NOT NULL REFERENCES users(id),
    balance    NUMERIC(14, 2) NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create table wallet_transaction 
CREATE TABLE IF NOT EXISTS wallet_transaction (
    id             SERIAL PRIMARY KEY,
    wallet_id      INT            NOT NULL REFERENCES wallet(id),
    type           VARCHAR(20)    NOT NULL
                   CHECK (type IN ('credit_waste', 'donation_out', 'withdraw')),
    amount         NUMERIC(14, 2) NOT NULL,
    reference_type VARCHAR(50),
    reference_id   INT,
    created_at     TIMESTAMP DEFAULT NOW()
);

-- Create table donation_campaign
CREATE TABLE IF NOT EXISTS donation_campaign (
    id             SERIAL PRIMARY KEY,
    title          VARCHAR(255)   NOT NULL,
    description    TEXT,
    organizer      VARCHAR(255),
    target_amount  NUMERIC(14, 2) NOT NULL,
    current_amount NUMERIC(14, 2) NOT NULL DEFAULT 0,
    status         VARCHAR(20)    NOT NULL DEFAULT 'draft'
                   CHECK (status IN ('draft', 'active', 'completed', 'closed')),
    created_by     INT REFERENCES users(id),
    created_at     TIMESTAMP DEFAULT NOW()
);

-- Create table donation
CREATE TABLE IF NOT EXISTS donation (
    id          SERIAL PRIMARY KEY,
    user_id     INT            NOT NULL REFERENCES users(id),
    campaign_id INT            NOT NULL REFERENCES donation_campaign(id),
    amount      NUMERIC(14, 2) NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Create table campaign_update 
CREATE TABLE IF NOT EXISTS campaign_update (
    id              SERIAL PRIMARY KEY,
    campaign_id     INT  NOT NULL REFERENCES donation_campaign(id),
    content         TEXT NOT NULL,
    report_file_url TEXT,
    published_at    TIMESTAMP DEFAULT NOW()
);
