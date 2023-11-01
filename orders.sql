drop table if exists transactions;
drop table if exists order_items;
drop table if exists orders;
drop table if exists products;
drop table if exists customers;

create table customers (
    id varchar(50) primary key
    ,name varchar(50) not null
);

create table products (
    id varchar(50) primary key
    ,name varchar(50) not null
    ,price money not null
);

create table orders (
    id varchar(50) primary key
    ,customer_id varchar(50) not null
    ,order_date timestamp not null
    ,"status" varchar(50) not null
    ,foreign key (customer_id) references customers(id)
);

create table order_items (
    id varchar(50) primary key
    ,order_id varchar(50) not null
    ,product_id varchar(50) not null
    ,quantity int not null
    ,price money not null
    ,foreign key (order_id) references orders(id)
    ,foreign key (product_id) references products(id)
);

create table transactions (
    id varchar(50) primary key
    ,order_id varchar(50) not null
    ,transaction_date timestamp not null
    ,amount money not null
    ,foreign key (order_id) references orders(id)
);
