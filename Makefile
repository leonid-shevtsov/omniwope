build:
	go build -o dist/omniwope-dev .

release:
	goreleaser release --clean

staging: build
	env OMNIWOPE_MASTODON_START_DATE=$(shell date +'%Y-%m-%d') \
			OMNIWOPE_TG_START_DATE=$(shell date +'%Y-%m-%d') \
		dist/omniwope-dev --config omniwope-staging.yml
