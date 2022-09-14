module github.com/qunv/visaurus/adapter

go 1.19

replace github.com/qunv/visaurus/core => ./../core

require (
	github.com/qunv/visaurus/core v0.0.0-00010101000000-000000000000
	gorm.io/gorm v1.23.8
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
)
