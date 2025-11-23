package main

import (
	"Avito/internal/controller"
	"Avito/internal/repository"
	"Avito/internal/usecase"
	"fmt"
	"time"
)

func main() {
	time.Sleep(5 * time.Second)
	baza := repository.NewDatabase()
	if baza == nil {
		fmt.Println("error to connect db")
	}

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
