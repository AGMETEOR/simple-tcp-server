.PHONY: run stop

run:
	docker run -p 8080:8080 simple-tcp-server

stop:
	docker ps | grep simple-tcp-server | awk '{ print $$1 }' | xargs -I {} docker stop {}
