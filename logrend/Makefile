build:
	docker build -t stephenhillier/logrend -f ./Dockerfile ../ && docker push stephenhillier/logrend && kubectl -n earthworks replace --force -f ../kubernetes/04-logrend-deploy.yaml