-- +goose Up
-- +goose StatementBegin
INSERT INTO ledger_accounts (id, name, description, parent_id) VALUES
-- Root Parents
(1, 'Assets', 'Value owned or held (Wealth)', NULL),
(2, 'Liabilities', 'Debts and obligations owed to others', NULL),
(3, 'Income', 'Inflow of value', NULL),
(4, 'Expenses', 'Outflow of value (Cost of living)', NULL),
(5, 'Equity', 'Net worth and opening balances', NULL),

-- Assets (100-199)
(101, 'Investments', 'Stocks, bonds, and brokerage holdings', 1),
(102, 'Cash', 'Physical currency and petty cash', 1),

-- Expense Top-Level Categories (4000-4999)
(4100, 'Obligations', 'Fixed commitments and taxes', 4),
(4200, 'Home', 'Housing, maintenance, and household logistics', 4),
(4300, 'Lifestyle', 'Discretionary spending and personal shopping', 4),
(4400, 'Food', 'Groceries and dining out', 4),
(4500, 'Transport', 'Commuting and travel costs', 4),
(4600, 'Health & Growth', 'Medical, wellness, and education', 4),

-- Obligations Sub-accounts
(4101, 'Insurance', 'Risk management premiums', 4100),
(4102, 'Tax', 'Income and personal taxes', 4100),
(4103, 'Parents', 'Financial support for family', 4100),

-- Home Sub-accounts
(4201, 'Maintenance', 'Home repairs and upkeep', 4200),
(4202, 'Cleaner', 'Domestic help and cleaning services', 4200),
(4203, 'Dog', 'Pet care, food, and vet bills', 4200),
(4204, 'Phone', 'Communication and data plans', 4200),

-- Lifestyle & Shopping Hierarchy
(4301, 'Shopping', 'General retail purchases', 4300),
(4302, 'Entertainment', 'Movies, music, and leisure', 4300),
(4303, 'Gifts', 'Presents and donations', 4300),
(4304, 'Holiday', 'Vacation and travel spending', 4300),
(4310, 'J', 'Personal shopping for J', 4301),
(4311, 'K', 'Personal shopping for K', 4301),
(4312, 'Handbags', 'Specific luxury accessory tracking', 4301),

-- Food Sub-accounts
(4401, 'Groceries', 'Supermarket and food supplies', 4400),
(4402, 'EatingOut', 'Restaurants, cafes, and takeout', 4400),

-- Transport Sub-accounts
(4501, 'Train', 'Public transport and rail', 4500),
(4502, 'Taxi', 'Ride-sharing and private hire', 4500),

-- Health & Growth Sub-accounts
(4601, 'Doctor', 'Healthcare and medical visits', 4600),
(4602, 'Education', 'Courses, books, and self-improvement', 4600);

-- +goose StatementEnd
