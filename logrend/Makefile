build:
	docker build -t stephenhillier/logrend . && docker push stephenhillier/logrend && kubectl -n earthworks replace --force -f ../kubernetes/04-logrend-deploy.yaml