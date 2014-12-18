# xlsx2sql

xlsx2sql is an open source Microsoft excel import data to mysql tools

## install

#### base install

1. install golang
2. download the project to your workspace
3. add xlsx2sql dir to your pc PATH

#### lib install

1. go get github.com/go-sql-driver/mysql
2. go get github.com/tealeg/xlsx
3. go get github.com/widuu/goini

## using

#### excel
|| mysql_tab_name ||  ||  ||
|| mysql_column1_name || mysql_column2_name || mysql_column3_name ||
|| mysql_comment || mysql_comment || mysql_comment ||
|| data || data || data ||
|| data || data || data ||
|| ... || ... || ... ||

#### mysql

```mysql
CREATE TABLE `mysql_tab_name` (
	`mysql_column1_name` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'mysql_comment',
	`mysql_column2_name` VARCHAR(20) NOT NULL DEFAULT '' COMMENT 'mysql_comment',
	`mysql_column3_name` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'mysql_comment',
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='mysql_comment' AUTO_INCREMENT=1;
```