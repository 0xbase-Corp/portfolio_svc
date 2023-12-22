@echo off
SETLOCAL EnableDelayedExpansion

REM Directory for migrations
IF NOT EXIST migrations (mkdir migrations)

REM Defining table categories and names
SET category[0]=app_specific
SET category[1]=users
SET category[2]=wallets
SET category[3]=debank
SET category[4]=moralis
SET category[5]=btc_com
SET category[6]=coingecko

SET tables_app_specific[0]=pseudonymous_portfolios
SET tables_app_specific[1]=portfolio_annotations

SET tables_users[0]=users
SET tables_users[1]=devices
SET tables_users[2]=user_tags

SET tables_wallets[0]=global_wallets

SET tables_debank[0]=evm_assets_debank_v1
SET tables_debank[1]=tokens_list
SET tables_debank[2]=nft_list

SET tables_moralis[0]=solana_assets_moralis_v1
SET tables_moralis[1]=tokens
SET tables_moralis[2]=nfts

SET tables_btc_com[0]=bitcoin_address_info

SET tables_coingecko[0]=token_prices
SET tables_coingecko[1]=bitcoin_btc_com_v1

REM Creating migration files with dummy SQL for each category and table
FOR /L %%i IN (0,1,6) DO (
    SET cat=!category[%%i]!
    IF NOT EXIST migrations\!cat! (mkdir migrations\!cat!)
    
    FOR /F "tokens=*" %%a IN ('set tables_!cat!') DO (
        FOR /F "tokens=2 delims==" %%b IN ("%%a") DO (
            SET tablename=%%b
            echo Creating migration files for !tablename! in category !cat!

            REM Writing dummy SQL in up script
            echo -- Create !tablename! table > migrations\!cat!\0001_create_!tablename!_table.up.sql
            echo CREATE TABLE !tablename! (); >> migrations\!cat!\0001_create_!tablename!_table.up.sql

            REM Writing dummy SQL in down script
            echo -- Drop !tablename! table > migrations\!cat!\0001_create_!tablename!_table.down.sql
            echo DROP TABLE IF EXISTS !tablename!; >> migrations\!cat!\0001_create_!tablename!_table.down.sql
        )
    )
)

echo Migration files with dummy SQL created in categorized directories.
ENDLOCAL
