module j/server

go 1.17

replace j/japi => ./../japi

replace j/schema => ./../schema

require j/japi v0.0.0-00010101000000-000000000000

require (
	github.com/lucasjones/reggen v0.0.0-20200904144131-37ba4fa293bb // indirect
	j/schema v0.0.0-00010101000000-000000000000 // indirect
)
