package usecase

import (
	"context"
	"time"

	"github.com/lassulfi/imersao11-consolidacao/internal/domain/entity"
	"github.com/lassulfi/imersao11-consolidacao/internal/domain/repository"
	"github.com/lassulfi/imersao11-consolidacao/pkg/uow"
)

type MatchInput struct {
	ID      string    `json:"id"`
	Date    time.Time `json:"match_date"`
	TeamAID string    `json:"team_a_id"`
	TeamBID string    `json:"team_b_id"`
}

type AddMatchUseCase struct {
	Uow uow.UowInterface
}

func NewAddMatchUseCase(uow uow.UowInterface) *AddMatchUseCase {
	return &AddMatchUseCase{
		Uow: uow,
	}
}

func (a *AddMatchUseCase) Execute(ctx context.Context, input MatchInput) error {
	return a.Uow.Do(ctx, func(uow *uow.Uow) error {
		matchRepository := a.getMatchRepository(ctx)
		teamRepository := a.getTeamRepository(ctx)

		teamA, err := teamRepository.FindByID(ctx, input.TeamAID)
		if err != nil {
			return err
		}
		teamB, err := teamRepository.FindByID(ctx, input.TeamBID)
		if err != nil {
			return err
		}

		match := entity.NewMatch(input.ID, teamA, teamB, input.Date)
		err = matchRepository.Create(ctx, match)
		if err != nil {
			return err
		}

		return nil
	})
}

func (a *AddMatchUseCase) getMatchRepository(ctx context.Context) repository.MatchRepositoryInterface {
	matchRepository, err := a.Uow.GetRepository(ctx, "MatchRepository")
	if err != nil {
		panic(err)
	}
	return matchRepository.(repository.MatchRepositoryInterface)
}

func (a *AddMatchUseCase) getTeamRepository(ctx context.Context) repository.TeamRepositoryInterface {
	teamRepository, err := a.Uow.GetRepository(ctx, "TeamRepository")
	if err != nil {
		panic(err)
	}
	return teamRepository.(repository.TeamRepositoryInterface)
}
