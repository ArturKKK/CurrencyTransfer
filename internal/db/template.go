package db

const (
	initRequest = `
		CREATE TABLE IF NOT EXISTS currency(
			char_code CHAR(3) PRIMARY KEY,
			vunit_rate NUMERIC(10,4) 
		);
	`
	dropRequest = `
		DROP TABLE IF EXISTS currency;
	`

	cleanRequest = `
		DELETE FROM currency;
	`

	saveRequest = `
		INSERT INTO currency(char_code, vunit_rate) 
		VALUES ($1, $2)
		ON CONFLICT (char_code) DO UPDATE 
		SET vunit_rate = EXCLUDED.vunit_rate;
	`

	selectCurrencyArr = `
		SELECT char_code, vunit_rate FROM currency
	`

	selectByCharCode = `
		SELECT * FROM currency 
			WHERE char_code=$1;
	`

	selectByValue = `
		SELECT * FROM currency 
			WHERE vunit_rate=$1;
	`
)
