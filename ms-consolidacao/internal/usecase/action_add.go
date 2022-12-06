package usecase

import (
	"context"

	"github.com/lassulfi/imersao11-consolidacao/internal/domain/entity"
	"github.com/lassulfi/imersao11-consolidacao/internal/domain/repository"
	"github.com/lassulfi/imersao11-consolidacao/pkg/uow"
)

type ActionAddInput struct {
	MatchID  string `json:"match_id"`
	TeamID   string `json:"team_id"`
	PlayerID string `json:"player_id"`
	Minute   int    `json:"minutes"`
	Action   string `json:"action"`
}

type AddActionUseCase struct {
	Uow         uow.UowInterface
	ActionTable entity.ActionTableInterface
}

func NewAddActionUseCase(uow uow.UowInterface, actionTable *entity.ActionTable) *AddActionUseCase {
	return &AddActionUseCase{
		Uow:         uow,
		ActionTable: actionTable,
	}
}

func (a *AddActionUseCase) Execute(ctx context.Context, input ActionAddInput) error {
	return a.Uow.Do(ctx, func(uow *uow.Uow) error {
		matchRepository := a.getMatchRepository(ctx)
		myTeamRepository := a.getMyTeamRepository(ctx)
		playerRepository := a.getPlayerRepository(ctx)

		match, err := matchRepository.FindByID(ctx, input.MatchID)
		if err != nil {
			return err
		}

		score, err := a.ActionTable.GetScore(input.Action)
		if err != nil {
			return err
		}
		theAction := entity.NewGameAction(input.PlayerID, input.Minute, input.Action, score)
		match.Actions = append(match.Actions, *theAction)
		err = matchRepository.SaveActions(ctx, match, float64(score))
		if err != nil {
			return err
		}

		player, err := playerRepository.FindByID(ctx, input.PlayerID)
		if err != nil {
			return err
		}
		player.Price += float64(score)
		err = playerRepository.Update(ctx, player)

		myTeam, err := myTeamRepository.FindByID(ctx, input.TeamID)
		if err != nil {
			return err
		}
		err = myTeamRepository.AddScore(ctx, myTeam, float64(score))

		return nil
	})
}

func (a *AddActionUseCase) getMatchRepository(ctx context.Context) repository.MatchRepositoryInterface {
	matchRepository, err := a.Uow.GetRepository(ctx, "MatchRepository")
	if err != nil {
		panic(err)
	}
	return matchRepository.(repository.MatchRepositoryInterface)
}

func (a *AddActionUseCase) getMyTeamRepository(ctx context.Context) repository.MyTeamRepositoryInterface {
	myTeamRepository, err := a.Uow.GetRepository(ctx, "MyTeamRepository")
	if err != nil {
		panic(err)
	}
	return myTeamRepository.(repository.MyTeamRepositoryInterface)
}

func (a *AddActionUseCase) getPlayerRepository(ctx context.Context) repository.PlayerRepositoryInterface {
	playerRepository, err := a.Uow.GetRepository(ctx, "PlayerRepository")
	if err != nil {
		panic(err)
	}
	return playerRepository.(repository.PlayerRepositoryInterface)
}
