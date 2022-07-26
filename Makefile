all: up

up:
	sudo docker compose up --force-recreate --remove-orphans -d
	make ps
	sudo docker logs -f tr_app

start:
	sudo docker compose restart
	make ps
	sudo docker logs -f tr_app

stop:
	sudo docker compose stop

down:
	sudo docker compose down --remove-orphans

ps:
	sudo docker compose ps
