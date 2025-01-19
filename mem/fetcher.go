package mem

type BytesFetcher interface {
	Fetch() ([]byte, error)
}

type CsvFetcher interface {
	FetchCsv() (*Table, error)
}

type NoSqlFetcher interface {
	FetchNoSql() (*NoSql, error)
}
