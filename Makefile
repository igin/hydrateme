deploy_api: # Deploy api to google appengine
	pushd api && gcloud app deploy && popd