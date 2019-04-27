USE punocracy;

/* Load data from the file of parsed homophones */

LOAD DATA INFILE '/var/lib/mysql-files/parsedHomophones.csv'
	INTO TABLE Words_T
	FIELDS TERMINATED BY ','
	LINES TERMINATED BY '\n';
