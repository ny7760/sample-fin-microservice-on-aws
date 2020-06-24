-- StockItems
create database if not exists StockItems;
grant all on StockItems.* to 'user'@'%';
commit;

use StockItems;

create table attributes (
    stock_code varchar(4), 
    stock_name varchar(200), 
    industry_code varchar(4),
    update_time varchar(50)
    );

-- sample data
insert into attributes values ('A001', 'Sample IT Company', "5250", "2020-05-17T09:00:00+09:00");
insert into attributes values ('A002', 'Sample IT Company2', "5250", "2020-05-17T09:00:00+09:00");

create table industries (
    industry_code varchar(4), 
    industry_name varchar(50)
    );
-- sample data
insert into industries values ('5250', 'Information & Communication');

create table time_prices (
    stock_code varchar(4), 
    date int,
    time varchar(5),
    price double,
    primary key(stock_code, date, time)
    );
-- sample data
insert into time_prices values ("A002", 20200526, "09:00", 2650);
insert into time_prices values ("A002", 20200526, "09:05", 2660);
insert into time_prices values ("A002", 20200526, "09:10", 2551);
insert into time_prices values ("A002", 20200526, "09:15", 2648);
insert into time_prices values ("A002", 20200526, "09:20", 2645);
insert into time_prices values ("A002", 20200526, "09:25", 2647);
insert into time_prices values ("A002", 20200526, "09:30", 2643);
insert into time_prices values ("A002", 20200526, "09:35", 2640);
insert into time_prices values ("A002", 20200526, "09:40", 2639);
insert into time_prices values ("A002", 20200526, "09:45", 2638);
insert into time_prices values ("A002", 20200526, "09:50", 2638);
insert into time_prices values ("A002", 20200526, "09:55", 2640);
insert into time_prices values ("A002", 20200526, "10:00", 2643);

commit;

-- StockTrade
create database if not exists StockTrade;
grant all on StockTrade.* to 'user'@'%';
commit;

use StockTrade;

create table trades (
    trade_no int,
    order_date int,
    order_time varchar(15),
    stock_code varchar(4),
    trade_type varchar(4),
    order_type varchar(1),  -- 0: 成行  1: 指値
    order_price double,
    order_quantity int,
    fee double,
    tax double,
    estimated_value double,
    trade_status varchar(1),  -- 0: 未約定  1: 約定済み  2: Fail  9: 取り消し
    update_time varchar(50),
    primary key(trade_no)
    );
-- sample data
insert into trades values (10001, 20200528, "09:00:00", "A002", "BUY", "1", 2400, 100, 1000, 100, 241100, "0", "2020-05-28T09:00:00+09:00");
insert into trades values (10002, 20200528, "11:20:00", "A002", "BUY", "0", 2300, 70, 1000, 100, 162100, "0", "2020-05-28T11:20:00+09:00");


create table contract_trades (
    trade_no int,
    contract_date int,
    contract_time varchar(15),
    stock_code varchar(4),
    trade_type varchar(4),
    order_type varchar(1),  -- 0: 成行  1: 指値
    contract_price double,
    contract_quantity int,
    fee double,
    tax double,
    contract_value int,  -- 約定金額
    settlement_amount int,  -- 精算金額
    update_time varchar(50),
    primary key(trade_no)
    );

commit;

-- StockBalance
create database if not exists StockBalance;
grant all on StockBalance.* to 'user'@'%';
commit;

use StockBalance;

create table cash_balances (
    cash_value int,
    update_time varchar(50)
);

-- sample data
insert into cash_balances values (10000000, "2020-05-28T11:20:00+09:00");

create table stock_balances (
    stock_code varchar(4),
    stock_value int,
    stock_quantity int,
    update_time varchar(50),
    primary key(stock_code)
);

/*
create table stock_moves (
    stock_code varchar(4),
    trade_no int,
    order_date int,
    order_time varchar(15),
    trade_type varchar(4),
    price double,
    quantity int,
    fee double,
    tax double,
    update_time varchar(50),
    primary key(stock_code, trade_no, order_date, order_time, trade_type)
);
*/