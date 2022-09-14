package ports

import "github.com/qunv/visaurus/core/entities"

type IWordPort interface {
	Insert(word *entities.Word) error
}
