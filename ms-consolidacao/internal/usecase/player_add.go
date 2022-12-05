package usecase

import (
	"context"

	"github.com/lassulfi/imersao11-consolidacao/internal/domain/entity"
	"github.com/lassulfi/imersao11-consolidacao/internal/domain/repository"
	"github.com/lassulfi/imersao11-consolidacao/pkg/uow"
)

type AddPlayerInput struct {
	ID           string
	Name         string
	InitialPrice float64
}

type AddPlayerUseCase struct {
	Uow uow.UowInterface
}

func (a *AddPlayerUseCase) Execute(ctx context.Context, input AddPlayerInput) error {
	playerRepository := a.getPlayerRepository(ctx)
	player := entity.NewPlayer(input.ID, input.Name, input.InitialPrice)
	err := playerRepository.Create(ctx, player)
	if err != nil {
		return err
	}
	a.Uow.CommitOrRollback()
	return nil
}

func (a *AddPlayerUseCase) getPlayerRepository(ctx context.Context) repository.PlayerRepositoryInterface {
	playerRepository, err := a.Uow.GetRepository(ctx, "PlayerRepository")
	if err != nil {
		panic(err)
	}
	return playerRepository.(repository.PlayerRepositoryInterface)
}
