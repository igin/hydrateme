build:
	go build

run: # Run the app locally with the appengine dev server
	dev_appserver.py .

deploy: # Deploy api to google appengine
	gcloud app deploy

test:
	DATASTORE_EMULATOR_HOST=localhost:8081 go test 

test-setup:
	gcloud beta emulators datastore start  \
		--data-dir .emulator \
		--no-store-on-disk \
		--host-port localhost:8081