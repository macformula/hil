package writer

const (
	_decimal = 10
	
	// format for 24-hour clock with minutes, seconds, and 4 digits
	// of precision after decimal (period and colon delimiter)
	_messageTimeFormat = "15:04:05.0000"

	// format for 24-hour clock with minutes and seconds (period delimiter)
	_filenameTimeFormat = "15-04-05"

	// format for year, month and day with two digits each (period delimiter)
	_filenameDateFormat = "2006-01-02"
)