drop database if exists starter_development;
drop user starter;

create database starter_development;
create user starter with password 'starter';
grant all privileges on database starter_development to starter;

\connect starter_development
create extension if not exists vector;
grant usage, create on schema public to starter;
