package usecase

import "github.com/GoodCodingFriends/gpay/source"

var defaultInteractor *interactor

type interactor struct {
	InteractorParams
}

type InteractorParams struct {
	Source source.Source
}

func Inject(p InteractorParams) {
	defaultInteractor = &interactor{
		InteractorParams: p,
	}
}
