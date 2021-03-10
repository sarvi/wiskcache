module wiskcache

go 1.16

replace utils => ./utils
replace cache => ./cache
replace manifest => ./manifest

require (
	cache v0.0.0-00010101000000-000000000000 // indirect
	lukechampine.com/blake3 v1.1.5 // indirect
	manifest v0.0.0-00010101000000-000000000000 // indirect
	utils v0.0.0-00010101000000-000000000000 // indirect
)
