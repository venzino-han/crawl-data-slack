init:
# wget chrome driver to data/
	@echo "\
	GROUPWARE_ID=\n\
	GROUPWARE_PW=\n\
	SLACK_BOT_TOKEN=\n\
	" >> .env
	@mkdir -p driver

cmd=
up:
	@COMMAND=${cmd} docker-compose up --build app
