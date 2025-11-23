package main

import (
	"Avito/internal/controller"
	"Avito/internal/repository"
	"Avito/internal/usecase"
)

func main() {
	baza := repository.NewDatabase()

	user := repository.NewUserRepo(baza)
	team := repository.NewTeamRepo(baza)
	pr := repository.NewPullRequestRepo(baza)

	prUse := usecase.NewPullRequestUseCase(pr, user, team)
	UserUse := usecase.NewUserUseCase(user, pr)
	teamUse := usecase.NewTeamUseCase(user, team)

	hendler := controller.NewHTTPHandler(*prUse, *teamUse, *UserUse)
	server := controller.NewHTTPServer(*hendler)

	server.Run()
}
