default_goal: run

install-test:
	test -f themes/papermod/README.md

run: install-test
	@echo "Running..."
	hugo server
	@echo "Done."

build: install-test
	@echo "Building..."
	rm -rf public
	hugo
	@echo "Done."

deploy: install-test build
	@echo "Deploying..."
	hugo deploy
