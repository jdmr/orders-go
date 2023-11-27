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
    ,foreign key (customer_id) references customers(id) on delete cascade
);

create table order_items (
    id varchar(50) primary key
    ,order_id varchar(50) not null
    ,product_id varchar(50) not null
    ,quantity int not null
    ,price money not null
    ,foreign key (order_id) references orders(id) on delete cascade
    ,foreign key (product_id) references products(id) on delete cascade
);

create table transactions (
    id varchar(50) primary key
    ,order_id varchar(50) not null
    ,transaction_date timestamp not null
    ,amount money not null
    ,foreign key (order_id) references orders(id) on delete cascade
);

insert into customers(id,name) values('1','John');
insert into customers(id,name) values('2','Mary');
insert into customers(id,name) values('3','Bob');

insert into products(id,name,price) values('1','Apple',1.00);
insert into products(id,name,price) values('2','Orange',2.00);
insert into products(id,name,price) values('3','Banana',3.00);

insert into orders(id,customer_id,order_date,"status") values('1','1','2019-01-01','PENDING');
insert into orders(id,customer_id,order_date,"status") values('2','1','2019-01-02','PAID');
insert into orders(id,customer_id,order_date,"status") values('3','2','2019-01-03','CANCELLED');
insert into orders(id,customer_id,order_date,"status") values('4','3','2019-01-04','PAID');
insert into orders(id,customer_id,order_date,"status") values('5','3','2019-01-05','PENDING');
insert into orders(id,customer_id,order_date,"status") values('6','3','2019-01-06','PAID');

insert into order_items(id,order_id,product_id,quantity,price) values('1','1','1',1,1.00);
insert into order_items(id,order_id,product_id,quantity,price) values('2','1','2',2,2.00);
insert into order_items(id,order_id,product_id,quantity,price) values('3','2','1',3,1.00);
insert into order_items(id,order_id,product_id,quantity,price) values('4','2','2',4,2.00);
insert into order_items(id,order_id,product_id,quantity,price) values('5','2','3',5,3.00);
insert into order_items(id,order_id,product_id,quantity,price) values('6','3','1',6,1.00);
insert into order_items(id,order_id,product_id,quantity,price) values('7','3','2',7,2.00);
insert into order_items(id,order_id,product_id,quantity,price) values('8','3','3',8,3.00);
insert into order_items(id,order_id,product_id,quantity,price) values('9','4','1',9,1.00);
insert into order_items(id,order_id,product_id,quantity,price) values('10','4','2',10,2.00);
insert into order_items(id,order_id,product_id,quantity,price) values('11','4','3',11,3.00);
insert into order_items(id,order_id,product_id,quantity,price) values('12','5','1',12,1.00);
insert into order_items(id,order_id,product_id,quantity,price) values('13','5','2',13,2.00);
insert into order_items(id,order_id,product_id,quantity,price) values('14','5','3',14,3.00);
insert into order_items(id,order_id,product_id,quantity,price) values('15','6','1',15,1.00);
insert into order_items(id,order_id,product_id,quantity,price) values('16','6','2',16,2.00);
insert into order_items(id,order_id,product_id,quantity,price) values('17','6','3',17,3.00);

insert into transactions(id,order_id,transaction_date,amount) values('2','2','2019-01-02',10.00);
insert into transactions(id,order_id,transaction_date,amount) values('4','4','2019-01-04',62.00);
insert into transactions(id,order_id,transaction_date,amount) values('6','6','2019-01-06',156.00);
