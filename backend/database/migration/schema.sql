-- Table 1 (User info)
CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	full_name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	hashed_password VARCHAR(255) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Table 2 (Store balance)
CREATE TABLE IF NOT EXISTS wallets (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	balance BIGINT NOT NULL DEFAULT 0,
	currency VARCHAR(10) NOT NULL DEFAULT 'VND',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Table 3 (Store transfer info between 2 wallets)
CREATE TABLE IF NOT EXISTS transfers (
	id BIGSERIAL PRIMARY KEY,
	from_wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    to_wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    amount BIGINT NOT NULL CHECK (amount > 0), -- The transferred money must be > 0
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Table 4 (Record all positive/negative changes of each wallet)
-- Used to audit. The sum of all entries of 1 wallet must equal the balance of this wallet
CREATE TABLE IF NOT EXISTS entries (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    amount BIGINT NOT NULL, -- Negative (withdraw/transfer) or positive (recharge)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);
CREATE INDEX IF NOT EXISTS idx_transfers_from_wallet_id ON transfers(from_wallet_id);
CREATE INDEX IF NOT EXISTS idx_transfers_to_wallet_id ON transfers(to_wallet_id);
CREATE INDEX IF NOT EXISTS idx_entries_wallet_id ON entries(wallet_id);