PROJECT=$(GCP_PROJECT)

.PHONY: dev
dev:
	docker-compose exec app realize start --run

deploy:
	gcloud app deploy -q

invoke:
	open https://$(PROJECT).uc.r.appspot.com/slack/invoke

build:
	gcloud builds submit --tag gcr.io/tarsierapps/tarsier-integration

run: build
	gcloud run deploy --image gcr.io/tarsierapps/tarsier-integration --platform managed
#https://tarsier-integration-iskizwszpq-uc.a.run.app
