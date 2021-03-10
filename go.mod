module wiskcache

go 1.16

replace utils => ./utils

replace cache => ./cache

replace exec => ./exec

replace manifest => ./manifest

require (
	cache v0.0.0-00010101000000-000000000000
	exec v0.0.0-00010101000000-000000000000
	manifest v0.0.0-00010101000000-000000000000
	utils v0.0.0-00010101000000-000000000000
)
