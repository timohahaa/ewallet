CREATE TABLE wallets (
    id UUID PRIMARY KEY NOT NULL,
    balance NUMERIC(10, 3) NOT NULL DEFAULT 0 CHECK ( balance >= 0 )
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    made_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    transfered_from UUID REFERENCES wallets (id),
    transfered_to UUID REFERENCES wallets (id),
    amount NUMERIC(10, 3) NOT NULL DEFAULT 0 CHECK ( amount >= 0 )
);
