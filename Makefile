PROJECT=$(GCP_PROJECT)

.PHONY: dev
dev:
	docker-compose exec app realize start --run

deploy:
	gcloud app deploy -q

invoke:
	open https://$(PROJECT).uc.r.appspot.com/slack/invoke
